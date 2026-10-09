# 北理工羽毛球抢票

一个用于北理工体育馆羽毛球场地预约的命令行小工具。程序会根据你输入的场馆、日期、场地和时间段持续检查目标时段，并在可下单时尝试提交订单；检测到未支付订单后，可选发送钉钉机器人通知。

## 使用提醒

- 本项目仅建议本人账号自用。
- 不要泄露 `token`、`openid`、钉钉 `Webhook` 配置，这些信息可以代表你的账号或通知机器人发起请求。
- 不要把填写了真实配置的 `config/config.go` 上传到公开仓库。
- 程序不能保证一定抢到场地，最终结果以体育馆预约系统实际订单为准。

## 来源与授权说明

本项目基于 [YaphetLee2002/BitBMTGrab](https://github.com/YaphetLee2002/BitBMTGrab) 进行个人学习和自用修改。原仓库当前未声明开源许可证，因此本仓库不对原始代码作额外授权声明，也不建议将其用于公开分发、商业用途或超出个人自用范围的场景。

如果你需要长期公开维护、分发或二次开发，请自行确认原作者授权，或联系原作者补充明确的开源许可证。

## 准备工作

建议优先使用以下环境：

- Windows
- WSL2 Debian

项目只在 Windows 和 WSL2 Debian 上做过实际测试。Go 官方提供 Windows、macOS、Linux 的安装包，所以 Linux/macOS 理论上也可以运行；主要差别通常只是编译输出文件名是否带 `.exe`。

你需要准备：

- Go 环境：[Go 官方下载页](https://go.dev/dl/)
- 抓包工具，例如 Reqable
- 微信，登录你自己的账号
- 钉钉自定义机器人，可选，用于抢到后通知，防止抢到后错过 15 分钟支付有效期
- Windows 弹窗通知，可选，适合电脑前使用

## 抓取 token 与 openid

下面以 Reqable 为例。

1. 打开 Reqable，开始抓包。
2. 使用微信打开“北理工体育馆”公众号的预约功能。
3. 在预约页面随便进行一些操作，例如查看场馆、日期、时段、提交预约的订单信息。
4. 在 Reqable 中筛选：

```text
http://gym.dazuiwl.cn
```

5. 找到请求里的 `token` 和 `openid`：
   - `token`：通常大多数接口的请求头里都能看到，可以找一找。
   - `openid`：通常在点击“提交订单-提交并支付订单”后相关接口（例如 `/api/order/submit`）的请求体里可以找到，不需要真支付。

目前公众号预约接口是 HTTP 请求，通常可以直接查看内容，不需要安装 CA 证书，只有抓 HTTPS 请求时，才涉及安装 CA 证书和信任证书。

## 安装 Go 环境

安装 Go 可参考 Go 官方下载页，安装完成后在终端验证：

```powershell
go version
```

`go.mod` 中写的是：

```text
go 1.23
```

使用 Go 1.23 或更新版本均可尝试，本地测试使用 Go 1.26.3 版本也能运行。

## 获取项目源码

方式一：使用 Git 克隆。

```powershell
git clone <仓库地址>
cd BitBMTGrab
```

方式二：下载 ZIP，解压后进入项目目录。

```powershell
cd <解压后的项目目录>
```

## 配置账号信息

打开 [config/config.go](config/config.go)，填写：

```go
const ApiToken = "替换为你的token"
const OpenId = "替换为你的openid"
```

示例：

```go
const ApiToken = "xxxxxxxx"
const OpenId = "xxxxxxxx"
```

这些配置不要泄露。

如果你经常使用同一个校区，也可以提前配置默认场馆 ID：

```go
const DefaultVenueID = 25
```

留 `0` 表示启动时继续询问场馆 ID：

```go
const DefaultVenueID = 0
```

如果你想少输入年份，也可以提前配置默认年份：

```go
const DefaultYear = 2026
```

配置后，运行时日期只需要输入 `MMDD`。留 `0` 表示运行时输入完整 `YYYYMMDD`。

## 配置钉钉通知（可选）

钉钉自定义机器人可以通过 Webhook 向群里发送消息。参考：[钉钉开发者百科：创建自定义机器人](https://open.dingtalk.com/document/dingstart/custom-bot-creation-and-installation)。

配置步骤：

1. 在钉钉群里添加“自定义机器人”。
2. 安全设置选择“自定义关键词”。
3. 关键词建议设置为：

```text
抢票通知
```

4. 复制机器人 Webhook URL，填写到 [config/config.go](config/config.go)：

```go
const DingTalkWebhook = "https://oapi.dingtalk.com/robot/send?access_token=xxx"
```

如果不需要钉钉通知，保持空字符串即可：

```go
const DingTalkWebhook = ""
```

测试钉钉通知：

```powershell
go run . --testdingtalk
```

如果已经编译过，也可以直接用编译后的文件测试：

```powershell
.\badminton.exe --testdingtalk
```

如果配置正确，钉钉群会收到一条测试消息。默认通知标题包含“抢票通知”，用于匹配钉钉机器人关键词策略。

## 配置 Windows 弹窗通知（可选）

如果你在 Windows 或 WSL2 Debian 中运行程序，可以启用 Windows 弹窗通知。它不是右下角 toast 通知，而是普通弹窗，抢到后需要你手动关闭。

打开 [config/config.go](config/config.go)，把：

```go
const EnableWindowsNotification = false
```

改成：

```go
const EnableWindowsNotification = true
```

测试 Windows 弹窗：

```powershell
go run . --testwin
```

如果已经编译过，也可以直接用编译后的文件测试：

```powershell
.\badminton.exe --testwin
```

WSL2 Debian 下也会尝试调用 Windows 的 `powershell.exe` 来显示弹窗。如果提示找不到 `powershell.exe`，请检查 WSL 是否启用了 Windows interop。Linux/macOS 不支持这种 Windows 弹窗通知，但仍可以使用钉钉通知。

## 运行程序

直接运行：

```powershell
go run .
```

编译后运行：

Windows：

```powershell
go build -o badminton.exe .
.\badminton.exe
.\badminton.exe --testdingtalk
.\badminton.exe --testwin
```

WSL2 Debian / Linux / macOS：

```bash
go build -o badminton .
./badminton
./badminton --testdingtalk
./badminton --testwin
```

输出文件名可以自己改，Windows 通常使用 `.exe` 后缀，Linux/macOS/WSL 通常不需要后缀。编译后的文件也支持测试参数；如果过几天发现通知异常，可以先用编译后的文件带参数重新测试，不必重新 `go run`。

## 输入预约参数

程序启动后会依次要求输入：

1. 场馆 ID
2. 预约日期
3. 场地名称或编号，一个或多个
4. 预约开始小时

场馆 ID：

```text
25 - 良乡
51 - 中关村
54 - 中关村早晚场
```

示例输入：

```text
请输入场馆ID（25-良乡，51-中关村，54-中关村早晚场）: 25
请输入预订日期（格式：YYYYMMDD）: 20260511
请输入场地名称或编号（多个用逗号分隔，最多5个）: 4
请输入预订时间（小时，例如9表示09:00开始的时段）: 17
```

其中：

- 日期推荐输入 `YYYYMMDD` 格式；如果在配置中填写了 `DefaultYear`，则只需要输入 `MMDD`。
- 场地可以输入程序打印的完整名称，也可以只输入名称末尾的数字编号，例如 `4` 自动匹配 `场4` 或 `主馆4`。编号不是表格列序号；同一编号对应多个场地时会报错，需要使用完整名称。
- 多个场地可以用英文逗号 `,` 或中文逗号 `，` 分隔，例如 `1,2,3` 或 `场1，场2，场3`；也可以用范围写法，例如 `1-5`、`场1-场5`，或混合输入 `1-3,5`。数字范围可以跨名称前缀，例如 `8-10` 自动匹配 `主馆8,副馆9,副馆10`。
- 一次最多输入 5 个场地。重复场地（包括 `1,主馆1` 这种匹配到同一场地的输入）、空场地名、场地名前后带空格都会报错退出。日志和通知使用实际场地名称。
- 时间只输入开始小时，例如 `17` 表示 `17:00` 开始的时段。

## 运行状态说明

程序会先打印目标日期的场地预订情况，然后进入检查和下单流程。

如果输入了多个场地，程序不会并发请求，而是按输入顺序错峰检查。总检查周期约 5 秒，相邻两个场地的间隔约为 `5秒 / 场地数量`，例如输入 5 个场地时，大约每 1 秒检查下一个场地，每个场地和自己的下一次检查间隔约 5 秒。抢到任意一个场地后，程序会发送通知并退出。

常见输出含义：

```text
[场1] 检查订单失败: 订单检查失败: 场地该时间段预约中，约5秒后继续检查...
```

目标时段当前不可下单，但属于可等待状态。程序会每 5 秒继续检查。

检查请求超时、连接失败、响应读取中断，或收到 HTTP 408、500、502、503、504 时，也会按原有错峰节奏继续检查，不会因一次临时故障直接退出。

```text
[场1] 检查通过，准备进入下单流程...
```

`CheckSportSchedule` 检查通过，程序准备尝试提交订单。

```text
[场1] 提交订单失败: xxx，下轮回到检查阶段...
```

提交接口失败，程序不会继续硬提交，而是在该场地下轮检查时重新回到检查阶段。

```text
[场1] 订单提交成功，开始确认订单...
```

提交接口返回成功，程序开始查询未支付订单列表。

```text
[场1] 已检测到订单 123，正在等待支付...
```

程序检测到未支付订单，说明已经抢到待支付订单，此时如果配置了钉钉 Webhook，会发送钉钉通知，如果启用了 Windows 弹窗通知，还会弹出一个 Windows 消息框提醒你前往支付。

建议每次使用时先观察一轮检查输出，因为接口返回的业务错误仍然只对“预约中”状态持续等待；如果接口返回的是其他状态，例如：

```text
场地该时间段上课中
场地该时间段包场中
```

程序会认为这不是可恢复的临时状态，并直接退出。

程序在 `07:00:00` 到 `23:30:00` 之间会尝试立刻下单；如果不在这个时间范围内，会等待到下一个早上 7 点再继续。

## 常见问题

### `go` 命令找不到

说明 Go 没安装好，或者环境变量没有配置好。重新安装 Go 后，打开一个新的终端执行：

```powershell
go version
```

### `token` 或 `openid` 失效

重新抓包并更新 [config/config.go](config/config.go) 中的 `ApiToken` 和 `OpenId`。

### 场地名称不匹配

完整场地名称按精确匹配处理，纯数字输入按场地名称末尾的编号匹配。请参考程序打印的场地名称；如果编号对应多个场地，需要输入完整名称以消除歧义。

### 钉钉通知不触发

检查：

- `DingTalkWebhook` 是否为空。
- Webhook URL 是否完整。
- 钉钉机器人是否启用了关键词安全策略。
- 钉钉机器人配置的关键词是否能匹配通知标题里的“抢票通知”。

可以先运行：

```powershell
go run . --testdingtalk
.\badminton.exe --testdingtalk
```

如果仍然失败，根据控制台打印的钉钉错误信息继续排查，也可以把错误信息发给 AI 辅助分析。

### Windows 弹窗通知不触发

检查：

- `EnableWindowsNotification` 是否设置为 `true`。
- 当前环境是否是 Windows 或 WSL2。
- WSL2 中是否可以直接执行 `powershell.exe`。

可以先运行：

```powershell
go run . --testwin
.\badminton.exe --testwin
```

如果仍然失败，根据控制台打印的 Windows 通知错误继续排查。

### 程序一直没有抢到

可能是目标时段一直没有释放、网络较慢、接口限制、账号状态异常或场地不可预约，程序只能自动执行检查和提交，不能绕过预约系统规则。
