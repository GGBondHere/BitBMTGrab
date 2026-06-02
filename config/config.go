package config

const ApiHost = "http://gym.dazuiwl.cn"
const ApiToken = "替换为你的token"
const OpenId = "替换为你的openid"

// DefaultVenueID 留 0 表示启动时询问场馆ID。
// 填写 25、51 或 54 后，程序会直接使用该场馆ID，不再询问。
const DefaultVenueID = 0

// DefaultYear 留 0 表示日期输入完整 YYYYMMDD。
// 填写年份，例如 2026 后，日期只需要输入 MMDD。
const DefaultYear = 0

// DingTalkWebhook 留空表示禁用钉钉通知。
// 启用时填写完整机器人 Webhook，例如：https://oapi.dingtalk.com/robot/send?access_token=xxx
// 钉钉机器人安全设置请选择“自定义关键词”，请确保给机器人配置的关键词包含在通知标题中；
// 默认通知标题包含“抢票通知”。
const DingTalkWebhook = ""

// EnableWindowsNotification 控制是否启用 Windows 弹窗通知。
// 仅适用于 Windows 或可调用 powershell.exe 的 WSL 环境；关闭时不影响程序运行。
const EnableWindowsNotification = false
