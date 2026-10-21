package app

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const ApiHost = "http://gym.dazuiwl.cn"

type Settings struct {
	ApiToken                  string `json:"api_token"`
	OpenId                    string `json:"openid"`
	DefaultVenueID            int    `json:"default_venue_id"`
	DefaultYear               int    `json:"default_year"`
	DingTalkWebhook           string `json:"dingtalk_webhook"`
	EnableWindowsNotification bool   `json:"enable_windows_notification"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取用户配置目录失败: %w", err)
	}
	return filepath.Join(dir, "BitBMTGrab", "config.json"), nil
}

func loadSettings(path string) (Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, fmt.Errorf("读取配置 %s 失败: %w", path, err)
	}

	var settings *Settings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return Settings{}, fmt.Errorf("解析配置 %s 失败，请修正文件后重试: %w", path, err)
	}
	if settings == nil {
		return Settings{}, fmt.Errorf("配置 %s 必须是 JSON 对象", path)
	}
	settings.normalize()
	err = settings.Validate()
	if err != nil {
		return Settings{}, fmt.Errorf("配置 %s 无效: %w", path, err)
	}
	return *settings, nil
}

func (s *Settings) normalize() {
	s.ApiToken = strings.TrimSpace(s.ApiToken)
	s.OpenId = strings.TrimSpace(s.OpenId)
	s.DingTalkWebhook = strings.TrimSpace(s.DingTalkWebhook)
}

func (s Settings) Validate() error {
	if s.DefaultVenueID < 0 {
		return fmt.Errorf("默认场馆 ID 不能小于 0")
	}
	if s.DefaultYear < 0 || s.DefaultYear > 9999 {
		return fmt.Errorf("默认年份必须是 1 到 9999，或用 0 表示不设置")
	}
	if s.DingTalkWebhook != "" {
		webhook, err := url.Parse(s.DingTalkWebhook)
		if err != nil || webhook.Scheme != "https" || webhook.Hostname() == "" || webhook.User != nil {
			return fmt.Errorf("钉钉 Webhook 必须是完整的 HTTPS URL")
		}
	}
	return nil
}

func (s Settings) ValidateAccount() error {
	if strings.TrimSpace(s.ApiToken) == "" || strings.TrimSpace(s.OpenId) == "" {
		return fmt.Errorf("账号信息不完整，请在“修改配置”中填写 token 和 openid")
	}
	return nil
}

func saveSettings(path string, settings Settings) error {
	settings.normalize()
	err := settings.Validate()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("编码配置失败: %w", err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	if runtime.GOOS != "windows" {
		err = os.Chmod(dir, 0700)
		if err != nil {
			return fmt.Errorf("设置配置目录权限失败: %w", err)
		}
	}

	// 先完整写入同目录临时文件，失败时不截断原配置。
	file, err := os.CreateTemp(dir, ".config-*.json")
	if err != nil {
		return fmt.Errorf("创建配置临时文件失败: %w", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	_, err = file.Write(data)
	if err == nil {
		err = file.Sync()
	}
	if err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("关闭配置临时文件失败: %w", err)
	}
	err = os.Rename(file.Name(), path)
	if err != nil {
		return fmt.Errorf("替换配置失败，原配置未修改: %w", err)
	}
	return nil
}
