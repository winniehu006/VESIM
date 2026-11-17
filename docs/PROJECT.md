# VESIM 项目说明

## 1. 项目定位

VESIM 是一个**自托管的浏览器控制台**，用于管理 Qualcomm 4G / LTE / 5G 调制解调器，提供短信、代理、eSIM、WiFi 通话 / IMS 语音、定时任务以及 MCP（AI 助手）接入能力。

本分支是 [HiDeck](https://github.com/yibaiba/hideck) 的改名分支，基于上游 `v2.1.35`。网页控制台、后端与文档均来自 HiDeck，继续遵循 **PolyForm Noncommercial 1.0.0** 许可；本分支与上游作者无隶属或背书关系。

> **重要边界**：本分支**只发布 Windows 便携版**（`vesim_*_windows_amd64.zip`）。驱动真实模组硬件所需的能力（QMI / MBIM / AT、VoWiFi、USB 音频、按网卡绑定出站、udev 热插拔）依赖 Linux 与 USB 直通，本分支不构建、不发布、不测试这些路径；需要模组部署请使用上游 HiDeck。

## 2. 发布与版本

| 项目 | 说明 |
| --- | --- |
| 唯一发布产物 | Windows x86_64 便携 zip，解压即用，无安装程序、无系统服务 |
| 发布页 | <https://github.com/winniehu006/VESIM/releases/latest> |
| 校验方式 | 每次发布同时上传 `.sha256`，并写入 Release 说明 |
| 版本规则 | 顺次递增（v2.1.35 → v2.1.36 → v2.1.37 → v2.1.38 …），不复用、不覆盖上游版本号 |
| 旧版本处理 | 发布页只保留最新版本；git tag 保留以维持版本递增记录 |
| 更新记录 | [CHANGELOG.md](../CHANGELOG.md) |

## 3. 主要功能

- **设备管理**：USB 自动发现；QMI、MBIM、AT、PC/SC；实时无线与 SIM 卡信息。
- **代理**：SOCKS5 / HTTP 代理，通过 `SO_BINDTODEVICE` 绑定到模组网卡（Linux 能力）。
- **短信**：收件箱、发送、联系人；支持从 iOS、Google、Samsung、小米、华为、OPPO、vivo 导入 vcf / csv。
- **语音 / 电话**：WiFi 通话、蜂窝软件 IMS、原生 VoLTE；保持、呼叫等待、仅监听接听。
- **VoWiFi**：SWu / IMS、通话录音；国家代理池、节点固定或直连。
- **原生 VoLTE**：模组 IMS（`phone_mode=volte`），不走 ePDG；仅支持中国特有的 MBN 映射。
- **eSIM**：下载、启用、禁用、重命名、删除；支持激活码或二维码 / PDF。
- **MBN 与 SIM 策略**：运营商配置库、自动选择规则、手动试运行。
- **自动化与通知**：定时任务；Telegram、邮件、Bark、飞书、企业微信、微信、QQ 通知。
- **MCP / AI 助手**：125 个工具，覆盖设备、网络、短信、联系人、通话、eSIM、MBN、代理、命令、余额、任务与诊断。

> 在 Windows 上，控制台启动后设备列表为空是**预期行为**：QMI / MBIM / AT、VoWiFi、USB 音频、按网卡绑定与 udev 热插拔都需要 Linux 主机 USB 直通。Windows 便携包适合体验 Web 控制台、界面、数据管理与 MCP 接入。

## 4. 技术组成

| 层次 | 技术 |
| --- | --- |
| 前端 | Vue 3 + TypeScript，位于 `web/` |
| 后端 | Go 1.26.4+，`cmd/vesim`，编译标签 `with_utls nomsgpack`，`CGO_ENABLED=0` |
| 数据存储 | SQLite（`data/vesim.db`），GORM |
| 实时语音 | WebRTC（ICE UDP 默认 `:7580`，端口占用时顺延） |
| 远程接入 | Streamable HTTP MCP，与 Web 共用 HTTP / HTTPS 端口，路径 `/mcp`，默认关闭 |
| 构建产物 | `dist/vesim_<tag>_windows_amd64.zip`，由 `build-windows.ps1` 生成 |
| 发布工作流 | `.github/workflows/windows-release.yml`（tag 触发），`.github/workflows/web-ci.yml`（前端路径触发） |

## 5. 仓库关键结构

| 路径 | 用途 |
| --- | --- |
| `cmd/vesim/` | 程序入口 |
| `internal/api/` | HTTP API 与 MCP 服务 |
| `internal/phone/`、`internal/volte/`、`internal/modem/` | 语音、IMS、模组相关逻辑 |
| `internal/db/` | 数据库模型与初始化 |
| `internal/web/` | 前端构建结果的嵌入位置（`dist/` 复制至此） |
| `web/`、`web-reader/` | 前端源码 |
| `config/config.example.yaml` | 配置模板，启动时复制为 `config/config.yaml` |
| `packaging/windows/start-vesim.bat` | Windows 启动脚本 |
| `build-windows.ps1` | Windows 便携包构建脚本 |
| `docs/` | 本文档与上游协议 / 专题文档 |
| `Dockerfile*`、`docker-compose*.yml`、`deploy*.sh`、`packaging/openwrt/`、`packaging/fnos/`、`Makefile` | 上游 Linux / Docker / OpenWrt 遗留文件，本分支不再构建、发布或测试；`deploy.sh` 与 `deploy-binary.sh` 启动即提示 Windows-only 并退出 |

## 6. 文档导航

| 文档 | 内容 |
| --- | --- |
| [操作手册](USER_GUIDE.md) | 下载、启动、登录、各功能页面使用、日常操作流程 |
| [附属手册](REFERENCE.md) | 配置项、命令行与环境变量、运行时文件、MCP 接入、备份、故障排查 |
| [维护手册](MAINTENANCE.md) | 构建、CI 发布、版本与 Release 管理、数据维护、平台限制与已知问题 |
| [MCP 指南](mcp.md) | MCP 鉴权、工具契约、客户端配置与能力边界 |
| [HTTPS / WebRTC](https-webrtc.md) | 麦克风权限、HTTPS 与 WebRTC 部署要求 |
| [WebDAV 备份](webdav-archive.md) | 短信、通话与录音的归档与安全本地清理 |

## 7. 许可与免责

- 本项目仅供个人学习、研究与技术测试，不得用于商业用途，也不作为生产电话系统。
- 与 Quectel、Qualcomm 及各运营商无关；使用时应遵守当地法律与运营商条款。
- 主许可：[PolyForm Noncommercial 1.0.0](../LICENSE)；`third_party/vowifi-go` 为 AGPL-3.0。分发二进制或镜像前请阅读 [THIRD_PARTY_NOTICES.md](../THIRD_PARTY_NOTICES.md)。
