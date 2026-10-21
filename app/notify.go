package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const dingTalkNoticeKeyword = "抢票通知"

type dingTalkMarkdownRequest struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

type dingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func sendDingTalkMarkdown(webhook, title, text string) (err error) {
	webhook = strings.TrimSpace(webhook)
	if webhook == "" {
		return nil
	}
	// HTTP 错误可能携带整个 URL，不能把 access_token 打进日志。
	defer func() {
		if err != nil {
			err = redactWebhookError(err, webhook)
		}
	}()

	var payload dingTalkMarkdownRequest
	payload.MsgType = "markdown"
	payload.Markdown.Title = title
	payload.Markdown.Text = text

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("构造钉钉通知失败: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Post(webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("发送钉钉通知失败: %w", err)
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("读取钉钉通知响应失败: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("钉钉通知 HTTP 状态异常: %s, 响应: %s", res.Status, string(responseBody))
	}

	var response dingTalkResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("解析钉钉通知响应失败: %w, 响应: %s", err, string(responseBody))
	}
	if response.ErrCode != 0 {
		return fmt.Errorf("钉钉通知返回错误: errcode=%d, errmsg=%s", response.ErrCode, response.ErrMsg)
	}

	return nil
}

func redactWebhookError(err error, webhook string) error {
	message := strings.ReplaceAll(err.Error(), webhook, "[Webhook]")
	parsed, parseErr := url.Parse(webhook)
	if parseErr == nil {
		for _, values := range parsed.Query() {
			for _, value := range values {
				if value != "" {
					message = strings.ReplaceAll(message, value, "[已隐藏]")
				}
			}
		}
	}
	return errors.New(message)
}

func sendOrderNotice(settings Settings, venueID int, venueName string, dateStr string, hour int, orderID int, amount int) error {
	title := dingTalkNoticeKeyword + " - 抢票成功，待支付"
	text := fmt.Sprintf(
		"### 抢票成功，待支付\n\n- 场馆ID: %d\n- 场地: %s\n- 日期: %s\n- 时间: %02d:00\n- 订单ID: %d\n- 金额: %d元\n\n请尽快前往支付。",
		venueID,
		venueName,
		dateStr,
		hour,
		orderID,
		amount,
	)

	return sendDingTalkMarkdown(settings.DingTalkWebhook, title, text)
}

func sendDingTalkTestNotice(settings Settings) error {
	if strings.TrimSpace(settings.DingTalkWebhook) == "" {
		return errors.New("钉钉通知未启用，请在“修改配置”中填写机器人 Webhook")
	}

	title := dingTalkNoticeKeyword + " - 测试消息"
	text := "### 测试消息\n\n这是一条钉钉机器人测试消息，用于验证 Webhook 是否可用。"

	return sendDingTalkMarkdown(settings.DingTalkWebhook, title, text)
}

func showWindowsOrderNotice(settings Settings, venueID int, venueName string, dateStr string, hour int, orderID int, amount int) error {
	if !settings.EnableWindowsNotification {
		return nil
	}

	title := dingTalkNoticeKeyword + " - 抢票成功，待支付"
	message := fmt.Sprintf(
		"抢票成功，待支付\n\n场馆ID: %d\n场地: %s\n日期: %s\n时间: %02d:00\n订单ID: %d\n金额: %d元\n\n请尽快前往支付。",
		venueID,
		venueName,
		dateStr,
		hour,
		orderID,
		amount,
	)

	return showWindowsMessageBox(title, message, false)
}

func showWindowsTestNotice(settings Settings) error {
	if !settings.EnableWindowsNotification {
		return errors.New("Windows 通知未启用，请在“修改配置”中启用 Windows 弹窗")
	}

	return showWindowsMessageBox(
		dingTalkNoticeKeyword+" - 测试消息",
		"这是一条 Windows 弹窗测试消息，用于验证弹窗通知是否可用。",
		true,
	)
}

func showWindowsMessageBox(title, message string, wait bool) error {
	script := fmt.Sprintf(
		"Add-Type -AssemblyName System.Windows.Forms\n[System.Windows.Forms.MessageBox]::Show(%s, %s, [System.Windows.Forms.MessageBoxButtons]::OK, [System.Windows.Forms.MessageBoxIcon]::Information) | Out-Null",
		powerShellString(message),
		powerShellString(title),
	)

	cmd := exec.Command("powershell.exe", "-NoProfile", "-STA", "-Command", script)
	if !wait {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("启动 Windows 弹窗失败: %w", err)
		}
		return nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		outputText := strings.TrimSpace(string(output))
		if outputText == "" {
			return fmt.Errorf("显示 Windows 弹窗失败: %w", err)
		}
		return fmt.Errorf("显示 Windows 弹窗失败: %w, 输出: %s", err, outputText)
	}

	return nil
}

func powerShellString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
