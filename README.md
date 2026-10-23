<h1 align="center">BITBadminton</h1>

<p align="center">
  <a href="https://github.com/GGBondHere/BITBadminton/releases/latest"><img alt="最新 Release 版本" src="https://img.shields.io/github/v/release/GGBondHere/BITBadminton?sort=semver&label=version&style=flat-square"></a>
  <a href="#下载运行"><img alt="支持 Windows、Linux、macOS" src="https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-0078D4?style=flat-square"></a>
</p>

<p align="center">北理工体育馆羽毛球预约工具，支持多场地错峰检查、钉钉和 Windows 弹窗通知。</p>

<p align="center">
  <a href="#准备信息">准备信息</a> ·
  <a href="#下载运行">下载运行</a> ·
  <a href="#首次配置">首次配置</a> ·
  <a href="#常见问题">常见问题</a> ·
  <a href="CHANGELOG.md">变更日志</a>
</p>

## 准备信息

开始前需要准备本人预约账号的 **Token 和 OpenID**。通知可选，使用钉钉时还需准备群机器人 Webhook。

### 获取 Token 和 OpenID

以 Reqable 为例：

1. 打开 Reqable 开始抓包，用微信进入“北理工体育馆”公众号的预约页面。
2. 查看场馆、日期和时段，筛选 `http://gym.dazuiwl.cn` 的请求。
3. 从请求头中获取 `token`。
4. 触发提交订单请求（如 `/api/order/submit`），从请求体中获取 `openid`，不需要实际支付；若产生未支付订单，请自行处理。

预约接口使用 HTTP，不需要安装 CA 证书；抓取 HTTPS 请求时才涉及证书和信任设置。**抓包记录包含账号凭据，请勿公开。**

### 准备钉钉机器人（可选）

在钉钉群创建自定义机器人，安全策略选择“自定义关键词”，填写 `抢票通知`，复制完整 Webhook，留待初始化时填写。例如：

```text
https://oapi.dingtalk.com/robot/send?access_token=xxx
```

创建步骤见 [钉钉机器人文档](https://open.dingtalk.com/document/dingstart/custom-bot-creation-and-installation)。

## 下载运行

从 [Releases](https://github.com/GGBondHere/BITBadminton/releases) 或下表下载，无需安装 Go。Intel/AMD 电脑选 `amd64`，ARM 电脑选 `arm64`；WSL 使用 Linux 版本。

| 系统 | amd64 | arm64 |
| --- | --- | --- |
| Windows | [下载](https://github.com/GGBondHere/BITBadminton/releases/latest/download/badminton_windows_amd64.exe) | [下载](https://github.com/GGBondHere/BITBadminton/releases/latest/download/badminton_windows_arm64.exe) |
| Linux / WSL | [下载](https://github.com/GGBondHere/BITBadminton/releases/latest/download/badminton_linux_amd64) | [下载](https://github.com/GGBondHere/BITBadminton/releases/latest/download/badminton_linux_arm64) |
| macOS | [Intel Mac](https://github.com/GGBondHere/BITBadminton/releases/latest/download/badminton_darwin_amd64) | [Apple 芯片 Mac](https://github.com/GGBondHere/BITBadminton/releases/latest/download/badminton_darwin_arm64) |

### Windows

下载后在文件所在目录打开 PowerShell，运行程序，例如：

```powershell
.\badminton_windows_amd64.exe
```

### Linux / WSL

一行安装并启动，自动识别架构：

```bash
curl -fsSL https://raw.githubusercontent.com/GGBondHere/BITBadminton/main/install.sh | sh
```

安装后运行 `bmt` 启动。也可从上表下载二进制，赋予执行权限后直接运行。

### macOS

下载后在终端中赋予执行权限并运行。以下以 Apple 芯片版本为例：

```bash
chmod +x badminton_darwin_arm64 && ./badminton_darwin_arm64
```

## 首次配置

首次启动依次填写三组表单，完成后保存：

| 表单 | 内容 |
| --- | --- |
| 账号信息 | Token、OpenID，均为必填 |
| 默认参数 | 场馆 ID、年份，留空或填 `0` 表示不预设 |
| 通知 | 钉钉 Webhook（可选）、Windows 弹窗开关（默认关闭） |

方向键选择、Enter 确认，Ctrl+C 取消当前表单；中途取消不保存。以后在主菜单的 **修改配置** 中编辑单项。

## 通知与测试

- **钉钉**：填写 Webhook 即启用，留空或选择“禁用钉钉通知（清空 Webhook）”关闭。仅支持关键词策略，不支持加签或 @。
- **Windows 弹窗**：普通消息框，需手动关闭，不需要管理员权限。支持 Windows 和可调用 `powershell.exe` 的 WSL，普通 Linux/macOS 不支持。

启用通知后，在主菜单选择 **测试通知 → 钉钉通知 / Windows 弹窗**。菜单只显示可用渠道；没有可测试项时，先在“修改配置”中启用。测试弹窗关闭后返回菜单。

确认目标未支付订单后才发送通知，**通知失败不影响订单，请及时支付**。正式弹窗启动后不等待关闭，钉钉通知等待接口响应，超时约 10 秒。

## 开始抢票

选择“开始抢票”，输入场馆、日期、场地和开始小时。已设置的默认场馆会直接使用，程序会显示场地预订表供参考。

```text
场馆 ID：25
预约日期：20261025（默认年份为 2026 时可输入 1025）
场地名称或编号：1-3,5
预约开始小时：17
```

| 参数 | 输入规则 |
| --- | --- |
| 场馆 | `25` 良乡、`51` 中关村、`54` 中关村早晚场 |
| 日期 | `YYYYMMDD` 或 `YYYY-MM-DD`；有默认年份时可输入 `MMDD`，不支持 `MM-DD` 或混用横杠 |
| 场地 | 完整名称或末尾编号，如 `场4`、`4`；多个用英文 `,` 或中文 `，` 分隔，如 `1,2,3` |
| 范围 | `1-5`、`场1-场5`、`1-3,5`；数字范围可跨名称前缀，如 `8-10` 匹配 `主馆8,副馆9,副馆10` |
| 时间 | 开始小时，如 `17` 表示 `17:00-18:00` |

一次最多 5 个备选场地，均为同一场馆、同一天、同一小时。不支持空项、空格、重复场地、倒序范围或前缀不一致的名称范围；编号有歧义时填写完整名称，编号不是表格列序号。

场地按输入顺序错峰检查，相邻间隔约 `5秒 / 场地数量`，一轮约 5 秒，请求耗时另计。抢到任意一个后停止其他场地，只提交单场地订单。

## 运行状态

- **预约中**：持续检查。上课中、包场中、超过每日预约次数等其他业务错误会退出，建议开始后先观察一轮。
- **网络异常**：各场地独立计数，正常响应后清零；任一场连续失败 5 次，退出当前任务。
- **检查通过**：立即尝试提交；提交失败则下轮重新检查。
- **提交成功**：每 5 秒确认本次预约的未支付订单，不重复提交。
- **已确认目标订单**：抢票成功但尚未支付，发送通知后退出。抢票过程中可用 Ctrl+C 停止。
- 可下单时段为系统本地时间 `07:00:00-23:30:00`，其他时间等待下一个早上 7 点；预约日期不会自动改变。请将运行环境设为北京时间。

## 常见问题

<details>
<summary><strong>Linux 服务器无法直接下载？</strong></summary>

可为终端配置代理，或使用自行信任的镜像源；也可在其他设备下载对应的 Linux 二进制，上传到服务器后赋予执行权限并运行。

</details>

<details>
<summary><strong>如何更新，配置会被覆盖吗？</strong></summary>

先停止程序，下载同平台、同架构的新版本，覆盖原可执行文件。Linux / WSL 可重新执行一行安装命令更新并启动，`bmt` 保持不变。

配置独立保存在用户目录，更新程序后使用同一系统用户启动即可沿用。

</details>

<details>
<summary><strong>配置保存在哪里？</strong></summary>

启动时会打印实际位置，默认如下：

| 环境 | 默认位置 |
| --- | --- |
| Windows | `%APPDATA%\BitBMTGrab\config.json` |
| Linux / WSL | `~/.config/BitBMTGrab/config.json`，遵循 `XDG_CONFIG_HOME` |
| macOS | `~/Library/Application Support/BitBMTGrab/config.json` |

配置按系统用户保存，不随工作目录改变。Windows、WSL、root 和普通用户默认各用一份。Linux/macOS 配置目录和文件权限分别为 `0700`、`0600`。

Linux 一行安装将程序保存为配置目录中的 `badminton`，更新时保留 `config.json`。短命令位于 `/usr/local/bin/bmt`（root）或 `~/.local/bin/bmt`（普通用户）；若不在 `PATH` 中，按安装提示添加。

</details>

<details>
<summary><strong>以前在 config.go 中填写过信息，如何迁移？</strong></summary>

首次启动时，将旧配置中的账号、默认参数和通知设置填入菜单即可。旧版配置不会自动导入。

</details>

<details>
<summary><strong>配置或通知不生效？</strong></summary>

先核对启动时打印的配置位置，修改后重新启动正在抢票的进程，运行期间不会热加载。Token / OpenID 失效时重新抓包，通过“修改配置”更新。

在主菜单选择“测试通知”：钉钉检查完整 URL 和关键词 `抢票通知`；Windows 检查开关和 `powershell.exe` 是否可用。“弹窗已启动”仅表示 PowerShell 启动成功，不保证弹窗已显示在最前面。

</details>

<details>
<summary><strong>为什么一直没有抢到？</strong></summary>

先在公众号确认目标时段是否可预约，再查看终端是否提示账号失效、网络异常或预约次数限制。

</details>

## 使用提醒

- 仅建议本人账号自用，程序不能保证抢到场地，最终以预约系统订单为准。
- 请遵守学校及预约平台规定；如禁止脚本或自动化预约，请勿使用本工具。违规使用导致的账号限制、预约取消等后果由使用者自行承担。
- Token、OpenID 和 Webhook 保存在本机明文配置中，界面遮蔽不等于加密；不要公开配置文件或账号凭据。

## 开发

安装 [Go 1.23+](https://go.dev/dl/)，克隆仓库后进入项目目录，可用 `go run .` 启动，或自行编译：

```bash
go build -o badminton.exe .  # Windows
go build -o badminton .      # Linux / WSL / macOS
```

源码编译同样使用 TUI 菜单和用户配置，本地构建版本号显示为 `dev`。

## 目录结构

```text
BITBadminton/
├── install.sh                  # Linux / WSL 安装与 bmt 短命令
├── main.go                     # 程序入口
├── app/
│   ├── app.go                  # 启动流程与命令行参数
│   ├── tui.go                  # 初始化、配置编辑与交互菜单
│   ├── badminton.go            # 场地选择、错峰检查与下单流程
│   ├── notify.go               # 钉钉和 Windows 弹窗通知
│   ├── api.go                  # 预约接口请求与网络错误处理
│   ├── config.go               # 用户配置加载、校验与保存
│   ├── model.go                # 接口数据结构
│   └── string.go               # 接口数据的文本输出
├── .github/workflows/
│   └── release.yml             # 按版本标签构建和发布程序
├── go.mod                      # Go 版本与依赖声明
├── go.sum                      # 依赖校验记录
├── CHANGELOG.md                # 版本变更日志
├── LICENSE                     # MIT 许可证
├── THIRD_PARTY_NOTICES.md      # 第三方版权与许可
└── README.md
```

## 许可证与致谢

许可信息见 [LICENSE](LICENSE) 和 [第三方声明](THIRD_PARTY_NOTICES.md)。

感谢 [YaphetLee2002/BitBMTGrab](https://github.com/YaphetLee2002/BitBMTGrab) 提供的初始实现。
