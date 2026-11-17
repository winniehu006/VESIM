# VESIM 维护手册

本手册面向维护者，覆盖 Windows 便携包构建、CI 发布、版本管理、数据维护、平台限制与发布检查清单。使用者请先看[操作手册](USER_GUIDE.md)，配置速查见[附属手册](REFERENCE.md)。

## 1. 维护范围

- 本仓库 `origin` 为 `winniehu006/VESIM`，`upstream` 为 `yibaiba/hideck`（仅用于对照上游）。
- **唯一发布产物**：Windows x86_64 便携 zip（`vesim_*_windows_amd64.zip`）+ 同名 `.sha256`。
- 上游的 Linux / Docker / OpenWrt / fnOS 路径只保留源码作参考，**不构建、不发布、不测试**：
  - `Dockerfile*`、`docker-compose*.yml`、`Makefile`；
  - `deploy.sh`、`deploy-binary.sh`：启动即打印「仅发布 Windows 包」并退出；
  - `packaging/openwrt/`、`packaging/fnos/`：README 已标注工作流已移除。
- 这些路径的自动构建工作流已从 `.github/workflows/` 删除，仅保留 `windows-release.yml` 与 `web-ci.yml`。

## 2. 构建 Windows 便携包

### 2.1 前置条件

| 依赖 | 版本 | 说明 |
| --- | --- | --- |
| Go | 1.26.4+ | 本地构建可用 `-GoBin <path>` 指定 `go\bin` 目录 |
| Node.js | 20+ | 仅在需要重新编译前端时必需 |
| PowerShell | 5.1 或 7 | 运行 `build-windows.ps1` |

构建为 `CGO_ENABLED=0`、`GOOS=windows`、`GOARCH=amd64`，因此不需要额外 C 工具链。前端 `npm run build` 实为 `typecheck && vite build`（前端约 397 个测试）。后端实际命令：

```bash
go build -trimpath -buildvcs=false -tags "with_utls nomsgpack" \
  -ldflags "-s -w -X 'github.com/winniehu006/vesim/internal/global.Version=<版本>' -X '.../global.BuildTime=<时间>'" \
  -o dist/vesim.exe ./cmd/vesim
```

### 2.2 构建命令

```powershell
# 完整构建（编译前端 + 后端）
.\build-windows.ps1

# 使用指定 Go、跳过前端编译（要求 internal/web/dist 已存在）
.\build-windows.ps1 -GoBin C:\go\bin -SkipFrontend

# 指定版本
.\build-windows.ps1 -Version v2.1.38
```

| 参数 | 行为 |
| --- | --- |
| `-Version` | 覆盖版本号；不传时取 `git describe --tags --abbrev=0`，若无 tag 则 `v0.0.0-dev`；自动补 `v` 前缀 |
| `-GoBin` | 临时把该目录加入 `PATH` |
| `-SkipFrontend` | 不再要求安装 npm，也不重新构建前端；但会校验 `internal/web/dist` 存在 |

### 2.3 产物与内容

- 输出：`dist\vesim_<tag>_windows_amd64.zip`。
- 包内文件：`vesim.exe`、`start-vesim.bat`、`config\config.example.yaml`、`LICENSE`、`THIRD_PARTY_NOTICES.md`。
- 版本串与构建时间通过 ldflags 注入 `internal/global.Version` / `BuildTime`；可在二进制中检索 `v<版本>` 验证。
- 本地构建与 CI 构建的 zip 大小通常不同：未压缩的 `vesim.exe` 大小相同，差异来自不同 PowerShell / .NET 的 zip 压缩实现与构建时间写入；**以 Release 的 SHA256 为准**。

### 2.4 本地检查

```bash
go vet ./cmd/... ./internal/... ./pkg/...
go test -run ZZZNoSuchTest -count=1 ./cmd/... ./internal/... ./pkg/...   # 仅验证可编译
npm test --prefix web
npm run typecheck --prefix web
```

实测冒烟：解压 zip → 运行 `vesim.exe`（或 `start-vesim.bat`）→ 浏览器访问 `http://localhost:<port>` 应返回 200 且 `<title>VESIM</title>`；删除 `data\` 后直接运行应自动重建目录并能访问；占用 WebRTC UDP 端口时应只打印一次「已改用」告警并继续启动。

## 3. CI 发布流程

### 3.1 `windows-release.yml`

- **触发**：推送 `v*` tag，或手动 `workflow_dispatch`（输入 `version`）。
- **权限**：`contents: write`。
- **Runner**：`windows-latest`。
- **步骤**：
  1. Checkout；
  2. 设置 Node.js 20（缓存 `web/package-lock.json`）；
  3. 设置 Go 1.26；
  4. `Resolve version`（tag 事件取 `GITHUB_REF_NAME`，手动事件取输入）；
  5. `Prepare frontend`：`npm ci --prefix web` → `npm run build --prefix web` → 复制到 `internal/web/dist`；
  6. `Vet and compile tests`：`go vet` 与 `go test -run ZZZNoSuchTest`（仅编译，不执行上游在 Windows 下会因文件锁失败的用例）；
  7. `Build portable package`：`.\build-windows.ps1 -Version <版本> -SkipFrontend`；
  8. `Write checksum`：生成 `.sha256`；
  9. `Compose release notes`：从校验文件读取 SHA256，生成 Windows 专用发布说明；
  10. `Create the release when missing`：Release 不存在则按同一说明创建；
  11. `Upload assets`：上传 zip 与 `.sha256`，`--clobber`；
  12. `Publish release notes with the checksum`：用同一说明更新标题与正文，保证重复构建同一 tag 结果一致。

### 3.2 `web-ci.yml`

- **触发**：仅当 `web/**` 或该文件自身变化的 PR / `main` push，以及手动触发。
- **步骤**：`npm ci` → `lint` → `typecheck` → `build`。
- 因此发布流程相关的提交不会触发前端 CI，这是路径过滤生效，不是漏跑。

## 4. 版本与 Release 管理

1. 确认 `main` 上的改动已就绪，`CHANGELOG.md` 已加入新版本条目。
2. 创建 annotated tag 并推送：

```bash
git tag -a vX.Y.Z -m "VESIM vX.Y.Z"
git push origin main
git push origin refs/tags/vX.Y.Z
```

3. 等待 `windows-release` run 成功，读取资产名与 `digest`。
4. 下载 zip 与 `.sha256`，实测 SHA256 应与 Release 说明、GitHub `digest` 三方一致；做一次浏览器冒烟测试。
5. 发布页只保留最新版本：删除旧 Release（保留 git tag 以维持版本递增记录）。
6. 版本号**顺次递增、不复用、不覆盖上游版本号**：v2.1.35 → v2.1.36 → v2.1.37 → v2.1.38 …
7. 若需重建同一版本：重推 tag 会触发 CI，`--clobber` 覆盖资产，发布说明自动写新的 SHA256。

> 网络不稳定时 `git push` 可能超时，重试即可；凭据由 Git Credential Manager 临时提供，用完后应吊销使用过的 PAT。

仓库元信息现状（GitHub API 实测）：description = `Windows portable web console for Qualcomm 4G/5G modem management — renamed fork of HiDeck`；homepage = `https://github.com/winniehu006/VESIM/releases/latest`；topics = `esim, golang, modem, quectel, self-hosted, sms, vowifi, vue, windows`。README / 中文 README 顶部提供四份文档导航，下载统一使用 `releases/latest`。

## 5. 数据与升级维护

- **状态目录**：`config\`、`data\`、`logs\` 都要保留；升级时替换 `vesim.exe`、`start-vesim.bat` 等发布文件即可。
- **数据库**：`data\vesim.db`（及 `-wal`、`-shm`）。2.1.37 起初始化前自动创建父目录；驱动由 `VESIM_SQLITE_DRIVER` 选择，默认 `modernc`，连接池 `MaxOpenConns(1)`；PRAGMA 为 `busy_timeout=5000`、`journal_mode=WAL`、`synchronous=NORMAL`、`foreign_keys=ON`。`AutoMigrate` 覆盖约 28 个模型，并按顺序执行短信身份、联系人 ID、国家规则主键、ICCID 重键、卡策略规范化等自定义迁移。修改表结构需新增迁移逻辑并验证旧库可升级。
- **国家表缓存**：`data\mcc-mnc-table.json`，启动时多源下载；离线时按未知国家直连。
- **录音 / MBN / 自动接听音频 / TLS 证书**：分别在 `data\recordings\voice\`、`data\mbn\`、`data\phone\auto-answer.wav`、`data\tls\`；MCP 令牌只存摘要与元数据于 `mcp_tokens` 表。
- **日志**：`logs\app-YYYY-MM-DD.log` 按天轮转、保留 30 天，同时输出到控制台、文件与 SSE 实时日志；`logs\app.log` 为入口（Windows 无创建符号链接权限时跳过软链）。
- **备份**：使用「设置 → WebDAV 归档」增量归档短信 / 通话 / 录音，或手工复制 `config\` 与 `data\`。归档台账保存在本地数据库中，清理需先预览再确认。

## 6. 平台限制与已知问题

| 项 | 说明 |
| --- | --- |
| Windows 无设备 | 预期行为：QMI/MBIM/AT、VoWiFi、USB 音频、按网卡绑定、udev 热插拔均需 Linux + USB 直通 |
| 模组直拨 / ALSA / 原生 AMR、MP3 编解码 | 代码在非 Linux 上明确返回 unavailable；电话页显示禁用原因 |
| PC/SC `decodeMultiString` 测试 | 只在 darwin/linux 后端存在，测试已移入平台限定文件，保证 Windows `go vet` 通过 |
| 上游测试临时文件锁失败 | `TestVoiceCallStoreRejectsBatchContainingActiveRecord`、`TestAnswerUsesReplacementMediaEndpoint`、`TestRecoveredAnswerStartsRecording` 在 Windows 因清理时文件被占用失败；原始代码同样复现，与本分支改动无关；CI 只做测试编译 |
| `gofmt -l` 报警 | 工作区 CRLF 检出导致，未修改文件同样被标记，不是格式问题 |
| 双栈地址误报告警 | 已修复：按端口比较，`:7580` → `[::]:7580` 不再误报 |
| 端口占用 | WebRTC UDP 自动顺延 10 个端口后退临时端口；HTTP 端口占用需改 `server.port` |
| 启动脚本换行符 | CI 检出为 LF、本地为 CRLF；已实测 LF 版在 cmd 下端口解析与标签跳转正常 |

## 7. 发布检查清单

发布前：

- [ ] 代码改动已完成，新增测试或平台限定测试在目标平台可编译；
- [ ] `go vet ./cmd/... ./internal/... ./pkg/...` 通过；
- [ ] `CHANGELOG.md` 已加入新版本条目；
- [ ] 版本号相对上一个 tag 顺次递增；
- [ ] README / 中文 README 下载链接仍为 `releases/latest`，文档导航有效。

发布后：

- [ ] CI run 成功，资产为 2 个且命名含正确版本；
- [ ] 下载后 SHA256 与 `.sha256`、Release 说明、GitHub `digest` 一致；
- [ ] 浏览器冒烟测试 200 / 标题正确，`data\` 能自动创建；
- [ ] 删除旧 Release，发布页只剩最新版本；
- [ ] 吊销使用过的 PAT。
