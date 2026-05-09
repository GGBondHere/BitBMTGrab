package main

import (
	"badminton/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func dingTalkWebhook() string {
	return strings.TrimSpace(config.DingTalkWebhook)
}

func sendDingTalkMarkdown(title, text string) error {
	webhook := dingTalkWebhook()
	if webhook == "" {
		return nil
	}

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

func sendOrderNotice(venueID int, venueName string, dateStr string, hour int, orderID int, amount int) error {
	title := dingTalkNoticeKeyword + " - 已检测到订单"
	text := fmt.Sprintf(
		"### 已检测到订单\n\n- 场馆ID: %d\n- 场地: %s\n- 日期: %s\n- 时间: %02d:00\n- 订单ID: %d\n- 金额: %d元\n\n请尽快前往支付。",
		venueID,
		venueName,
		dateStr,
		hour,
		orderID,
		amount,
	)

	return sendDingTalkMarkdown(title, text)
}

func sendDingTalkTestNotice() error {
	if dingTalkWebhook() == "" {
		return errors.New("钉钉通知未启用：请先在 config.DingTalkWebhook 中填写完整机器人 Webhook")
	}

	title := dingTalkNoticeKeyword + " - 测试消息"
	text := "### 测试消息\n\n这是一条钉钉机器人测试消息，用于验证 Webhook 是否可用。"

	return sendDingTalkMarkdown(title, text)
}
