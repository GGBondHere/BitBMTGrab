package app

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

func isFormCancelled(err error) bool {
	return errors.Is(err, huh.ErrUserAborted)
}

func runForm(fields ...huh.Field) error {
	return huh.NewForm(huh.NewGroup(fields...)).WithShowHelp(false).Run()
}

func inputForm(title string, value *string, validate func(string) error) error {
	return runForm(huh.NewInput().Title(title).Value(value).Validate(validate))
}

func requiredValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("此项不能为空")
	}
	return nil
}

func positiveNumber(value string) error {
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return fmt.Errorf("请输入大于 0 的整数")
	}
	return nil
}

func optionalNumber(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return 0, fmt.Errorf("请输入非负整数，留空表示不设置")
	}
	return number, nil
}

func validateOptionalNumber(value string) error {
	_, err := optionalNumber(value)
	return err
}

func validateYear(value string) error {
	year, err := optionalNumber(value)
	if err != nil {
		return err
	}
	if year > 9999 {
		return fmt.Errorf("年份不能大于 9999")
	}
	return nil
}

func validateWebhook(value string) error {
	settings := Settings{DingTalkWebhook: strings.TrimSpace(value)}
	return settings.Validate()
}

func secretInput(title string, value *string, validate func(string) error) *huh.Input {
	return huh.NewInput().Title(title).EchoMode(huh.EchoModePassword).Value(value).Validate(validate)
}

func initializeSettings(path string, settings *Settings) error {
	next := *settings
	var venue, year string
	form := huh.NewForm(
		huh.NewGroup(
			secretInput("Token", &next.ApiToken, requiredValue),
			secretInput("OpenID", &next.OpenId, requiredValue),
		).Title("账号信息"),
		huh.NewGroup(
			huh.NewInput().Title("默认场馆 ID（可选）").Value(&venue).Validate(validateOptionalNumber),
			huh.NewInput().Title("默认年份（可选）").Value(&year).Validate(validateYear),
		).Title("默认参数"),
		huh.NewGroup(
			secretInput("钉钉 Webhook（可选）", &next.DingTalkWebhook, validateWebhook),
			huh.NewConfirm().Title("Windows 弹窗通知").Affirmative("启用").Negative("禁用").Value(&next.EnableWindowsNotification),
		).Title("通知"),
	).WithShowHelp(false)
	err := form.Run()
	if err != nil {
		return err
	}
	next.ApiToken = strings.TrimSpace(next.ApiToken)
	next.OpenId = strings.TrimSpace(next.OpenId)
	next.DingTalkWebhook = strings.TrimSpace(next.DingTalkWebhook)
	next.DefaultVenueID, _ = optionalNumber(venue)
	next.DefaultYear, _ = optionalNumber(year)
	err = saveSettings(path, next)
	if err != nil {
		return err
	}
	*settings = next
	fmt.Println("配置已保存")
	return nil
}

func runMenu(version, path string, settings *Settings, initialize bool) error {
	fmt.Printf("BITBadminton %s\n配置位置: %s\n", version, path)
	if initialize {
		err := initializeSettings(path, settings)
		if err != nil {
			return err
		}
	}
	for {
		var action string
		err := runForm(huh.NewSelect[string]().Title("北理工羽毛球抢票").Options(
			huh.NewOption("开始抢票", "book"),
			huh.NewOption("修改配置", "settings"),
			huh.NewOption("测试通知", "notify"),
			huh.NewOption("退出", "exit"),
		).Value(&action))
		if err != nil {
			return err
		}
		switch action {
		case "book":
			err = settings.ValidateAccount()
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = runBooking(*settings)
			if isFormCancelled(err) {
				continue
			}
			return err
		case "settings":
			err = editSettings(path, settings)
		case "notify":
			err = testNoticeMenu(*settings)
		case "exit":
			return nil
		}
		if err != nil && !isFormCancelled(err) {
			fmt.Println(err)
		}
	}
}

func configured(value string) string {
	if value == "" {
		return "未设置"
	}
	return "已设置"
}

func editSettings(path string, settings *Settings) error {
	for {
		var field string
		err := runForm(huh.NewSelect[string]().Title("修改配置").Options(
			huh.NewOption("Token（"+configured(settings.ApiToken)+"）", "token"),
			huh.NewOption("OpenID（"+configured(settings.OpenId)+"）", "openid"),
			huh.NewOption(fmt.Sprintf("默认场馆 ID：%d", settings.DefaultVenueID), "venue"),
			huh.NewOption(fmt.Sprintf("默认年份：%d", settings.DefaultYear), "year"),
			huh.NewOption("钉钉 Webhook（"+configured(settings.DingTalkWebhook)+"）", "webhook"),
			huh.NewOption("禁用钉钉通知（清空 Webhook）", "clear-webhook"),
			huh.NewOption(fmt.Sprintf("Windows 弹窗：%t", settings.EnableWindowsNotification), "windows"),
			huh.NewOption("返回", "back"),
		).Value(&field))
		if err != nil || field == "back" {
			return err
		}

		next := *settings
		switch field {
		case "token":
			err = runForm(secretInput("Token", &next.ApiToken, requiredValue))
			next.ApiToken = strings.TrimSpace(next.ApiToken)
		case "openid":
			err = runForm(secretInput("OpenID", &next.OpenId, requiredValue))
			next.OpenId = strings.TrimSpace(next.OpenId)
		case "venue", "year":
			number, title, validate := next.DefaultVenueID, "默认场馆 ID（0 或留空表示不设置）", validateOptionalNumber
			if field == "year" {
				number, title, validate = next.DefaultYear, "默认年份（0 或留空表示不设置）", validateYear
			}
			value := ""
			if number != 0 {
				value = strconv.Itoa(number)
			}
			err = inputForm(title, &value, validate)
			if field == "year" {
				next.DefaultYear, _ = optionalNumber(value)
			} else {
				next.DefaultVenueID, _ = optionalNumber(value)
			}
		case "webhook":
			err = runForm(secretInput("钉钉 Webhook", &next.DingTalkWebhook, func(value string) error {
				err := requiredValue(value)
				if err != nil {
					return err
				}
				return validateWebhook(value)
			}))
			next.DingTalkWebhook = strings.TrimSpace(next.DingTalkWebhook)
		case "clear-webhook":
			var clear bool
			err = runForm(huh.NewConfirm().Title("清空钉钉 Webhook 并禁用通知？").Affirmative("清空").Negative("保留").Value(&clear))
			if err == nil && !clear {
				continue
			}
			next.DingTalkWebhook = ""
		case "windows":
			err = runForm(huh.NewConfirm().Title("Windows 弹窗通知").Affirmative("启用").Negative("禁用").Value(&next.EnableWindowsNotification))
		}
		if err != nil {
			return err
		}
		err = saveSettings(path, next)
		if err != nil {
			return err
		}
		*settings = next
		fmt.Println("配置已保存")
	}
}

func testNoticeMenu(settings Settings) error {
	var options []huh.Option[string]
	if strings.TrimSpace(settings.DingTalkWebhook) != "" {
		options = append(options, huh.NewOption("钉钉通知", "dingtalk"))
	}
	if settings.EnableWindowsNotification {
		_, err := exec.LookPath("powershell.exe")
		if err == nil {
			options = append(options, huh.NewOption("Windows 弹窗", "windows"))
		}
	}
	if len(options) == 0 {
		fmt.Println("暂无可测试的通知，请先在“修改配置”中启用通知；Windows 弹窗还需当前环境可调用 powershell.exe。")
		return nil
	}
	options = append(options, huh.NewOption("返回", "back"))

	var target string
	err := runForm(huh.NewSelect[string]().Title("测试通知").Options(options...).Value(&target))
	if err != nil || target == "back" {
		return err
	}
	if target == "dingtalk" {
		err = sendDingTalkTestNotice(settings)
	} else {
		err = showWindowsTestNotice(settings)
	}
	if err != nil {
		return fmt.Errorf("通知测试失败: %w", err)
	}
	fmt.Println("通知测试成功")
	return nil
}
