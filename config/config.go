package config

const ApiHost = "http://gym.dazuiwl.cn"
const ApiToken = "替换为你的token"
const OpenId = "替换为你的openid"

// DingTalkWebhook 留空表示禁用钉钉通知。
// 启用时填写完整机器人 Webhook，例如：https://oapi.dingtalk.com/robot/send?access_token=xxx
// 钉钉机器人安全设置请选择“自定义关键词”，请确保给机器人配置的关键词包含在通知标题中；
// 默认通知标题包含“抢票通知”。
const DingTalkWebhook = ""

// EnableWindowsNotification 控制是否启用 Windows 弹窗通知。
// 仅适用于 Windows 或可调用 powershell.exe 的 WSL 环境；关闭时不影响程序运行。
const EnableWindowsNotification = false
