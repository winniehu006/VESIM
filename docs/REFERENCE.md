# VESIM 附属手册（配置与参考）

本手册给出配置项、命令行参数、环境变量、运行时文件、MCP / AI 助手接入、WebDAV 备份与故障排查的速查信息。配置模板为 [`config/config.example.yaml`](../config/config.example.yaml)。

## 1. 配置文件

程序默认读取 `config/config.yaml`；Windows 启动脚本首次运行时会从 `config.example.yaml` 复制。配置支持热加载。

### 1.1 `server`

| 字段 | 默认 | 说明 |
| --- | --- | --- |
| `port` | `7575` | HTTP 管理端口；不写冒号时自动补全为 `:7575` |
| `https_enabled` | 模板 `false`（代码缺省兼容旧版为 `true`） | 内置 HTTPS 开关；反向代理场景请保持 `false` |
| `https_port` | `7576` | 内置 HTTPS 端口 |
| `webrtc_udp_address` | `:7580` | WebRTC 语音 UDP；占用时顺延最多 10 个端口，再退到临时端口，不终止启动 |
| `webrtc_public_host` | 空 | NAT 后的 ICE 公网域名；启动及每次会话重新解析 |
| `tls_data_dir` | `data/tls` | 本地 CA / 证书目录，权限 0700 |
| `tls_cert_file` / `tls_key_file` | 空 | 同时设置时覆盖本地 CA |
| `ice_servers` | `[]` | STUN / TURN 列表，如 `["stun:stun.example.com:3478"]` |
| `debug` | `false` | 调试日志级别 |

### 1.2 `web`

| 字段 | 默认 | 说明 |
| --- | --- | --- |
| `username` / `password` | `admin` / `admin` | 首次登录必须修改；保存后密码以 bcrypt 哈希写入配置 |
| `PROXY_WEB_PASSWORD` | 环境变量 | 若设置，优先于配置文件，且控制台不会覆盖 |

### 1.3 `outbound_limits`（按 ICCID 独立计数，重启不清零）

| 字段 | 默认 | 说明 |
| --- | --- | --- |
| `calls.interval_seconds` / `per_hour` / `per_day` | `10` / `20` / `100` | 外呼间隔与配额 |
| `sms.interval_seconds` / `per_hour` / `per_day` | `5` / `30` / `100` | 短信间隔与配额 |
| `disabled` | 未设置 | 必须显式 `true` 才关闭；`0` 或省略继承默认，负数报错 |
| `cards` | `{}` | 按带引号的 ICCID 单独覆盖；详见 [`outbound-limits.md`](outbound-limits.md) |

旧字段 `server.sms_rate_limit_disabled` 已废弃，启动时仅打印迁移提示。

### 1.4 其他常用段

| 段 | 关键字段 | 说明 |
| --- | --- | --- |
| `modem_voice` | `auto_enable_adb: true` | 模组直拨时准备 QDC507GLEFM21；设 `false` 禁止自动授权 / 改 USB（Linux） |
| `archive` | `enabled:false`、`interval_hours:24`（1–720）、`retention_days:30`（1–3650）、`auto_cleanup:false`、`include_sms/calls/recordings:true` | WebDAV 归档；密码留空保留原值，详见 [`webdav-archive.md`](webdav-archive.md) |
| `system` | `openwrt_dynamic_interfaces:false`、`auto_answer.audio_file: data/phone/auto-answer.wav`、`auto_answer.devices:{}` | 自动接听每设备可设 `enabled`、`delay_seconds` |
| `devices` | `[]`；示例字段 `id/name/device_backend/esim_transport/pcsc_reader_name/pcsc_usb_path/sim_pin_env` | `esim_transport` 枚举 `at|qmi|mbim|pcsc`；`sim_pin_env` 只存环境变量**名** |
| `proxy` | `instances: []` | SOCKS5 / HTTP 代理实例，字段含 `mode/listen_addr/listen_port/auth_enabled` |
| `vowifi` | `enabled:false`（另有 `device_id`、`mode`，默认 `vowifi`） | 全局 VoWiFi 开关（Linux） |
| `webhook` / `telegram` / `qq` / `weixin` / `wecom_bot` / `wecom` | 默认 `enabled:false` | 通知渠道；`telegram.recording_mode` 为 `voice` 或 `audio`；`webhook.timeout_ms:5000`、`retry_max:3` |
| `notification_device_filters` | 注释示例 | 按渠道（telegram/feishu/qq/weixin/wecom_bot/webhook/bark/email/pushplus/wecom）过滤设备 |
| `mcp` | `enabled:false` | 远程 MCP 开关；与 Web 共用端口，令牌在网页管理 |

> 代码还内置 `feishu`、`bark`、`email`、`pushplus` 默认配置，但未写入模板；以实际页面表单为准。

## 2. 命令行参数与环境变量

### 2.1 命令行

```
vesim.exe -c config\config.yaml
```

| 参数 | 默认 | 说明 |
| --- | --- | --- |
| `-c` | `config/config.yaml` | 指定配置文件路径 |
| `-backend-only` | `false` | 不挂载内嵌 Web UI，仅提供 API |

### 2.2 配置覆盖（`PROXY_` 前缀）

viper 使用 `PROXY` 作为环境变量前缀，配置键中的 `.` 替换为 `_`，例如：

| 变量 | 覆盖的配置 |
| --- | --- |
| `PROXY_WEB_PASSWORD` | `web.password`（优先级最高） |
| `PROXY_SERVER_WEBRTC_UDP_ADDRESS` | `server.webrtc_udp_address` |
| `PROXY_SERVER_WEBRTC_PUBLIC_HOST` | `server.webrtc_public_host` |
| `PROXY_DEVICES_0_APN` | `devices[0].apn` |

### 2.3 `VESIM_` 变量

| 变量 | 用途 |
| --- | --- |
| `VESIM_SIM_PIN_READER1` | SIM PIN（配置中只写变量名，不写明文） |
| `VESIM_SQLITE_DRIVER` | SQLite 驱动，默认 `modernc`；可选 `purego/go/glebarez`，cgo 驱动禁用 |
| `VESIM_QMI_SERIAL_REQUESTS=1` | EC20 + macOS/OrbStack 的 QMI 串行兼容，仅精确值 `1` 启用 |
| `VESIM_PROBE_TOKEN` | `modem-voice-probe` 工具鉴权 |
| `VESIM_FORCE_GO_RUN_LOG` | 覆盖 `go run` 判定以调整日志格式 |
| `VESIM_E911_TRACE_RAW` | E911 原始报文跟踪 |
| `VITE_API_PROXY_TARGET` | 前端开发代理目标，如 `http://127.0.0.1:7575`（`web/.env.local`） |

部署脚本还会读写 `VESIM_DIR`、`VESIM_ARCH`、`VESIM_VERSION`、`VESIM_DOMAIN` 等变量，但这些脚本在本分支只会提示「仅发布 Windows 版」后退出。

## 3. 运行时文件与目录

| 路径 | 作用 |
| --- | --- |
| `config/config.yaml` | 主配置；模板来自 `config.example.yaml` |
| `config/notification-state.json` | 通知渠道绑定状态（Telegram 自动绑定、微信 iLink 登录态） |
| `config/carrier_overrides.yaml` | 运营商覆盖，可选；加载失败回退内置配置 |
| `data/vesim.db`（含 `-wal`、`-shm`） | SQLite 主库，含短信、通话、MCP 令牌摘要、外呼计数、归档台账；父目录自动创建 |
| `data/mcc-mnc-table.json` | MCC/MNC 国家表缓存，启动时从多个源下载 |
| `data/recordings/voice/` | 通话录音（MP3 / PCAP） |
| `data/mbn/` | MBN 运营商配置库 |
| `data/phone/auto-answer.wav` | 自动接听音频 |
| `data/tls/` | 内置 HTTPS 时生成本地 CA 与服务器证书（0600） |
| `logs/app-YYYY-MM-DD.log` | 按天轮转的日志，保留 30 天；另有 `logs/app.log` 入口 |

> **升级 / 迁移**：保留整个 `data/` 与 `config/` 目录即可保留数据与设置；替换为新版本的 `vesim.exe`、`start-vesim.bat` 与其他发布文件。

## 4. MCP / AI 助手接入

1. 登录 Web 控制台并接受免责声明。
2. **设置 → MCP 接入**：开启开关，创建命名令牌。令牌明文**只显示一次**，请立即保存。
3. 客户端配置（`Streamable HTTP`，鉴权头 `Authorization: Bearer <MCP_TOKEN>`）：

```json
{
  "mcpServers": {
    "vesim": {
      "type": "http",
      "url": "http://localhost:7575/mcp",
      "headers": { "Authorization": "Bearer <MCP_TOKEN>" }
    }
  }
}
```

| 客户端 | 配置文件位置 |
| --- | --- |
| Codex | `~/.codex/config.toml`（`mcp_servers.vesim` + `http_headers`） |
| Claude Code | `~/.claude.json` |
| Cursor | `~/.cursor/mcp.json` |
| Gemini CLI | `~/.gemini/settings.json`（`httpUrl`） |
| Hermes Agent | `~/.hermes/config.yaml` |
| DeepSeek dsh | `$DSH_HOME/cordis.patch.yml`（`streamable-http`） |
| Pi Agent | `~/.pi/agent/mcp.json`（`pi mcp list` 验证） |
| Grok Build | `~/.grok/config.toml`（`grok mcp doctor vesim` 验证） |

- 端点为 VESIM 地址 + `/mcp`，支持 HTTP / HTTPS；远程访问需使用 Agent 能到达的地址，如 `http://<主机>:7575/mcp`。
- 令牌拥有全部已暴露的业务权限，**不能**用于网页登录或 REST 管理；网页会话也不能访问 MCP。可在设置中撤销 / 删除。
- 共 125 个工具，均提供 input / output schema；**写操作不会自动重试**。
- 浏览器音频通话：电话页「连接 MCP」，授权麦克风，保持页面打开；Agent 使用 `audio_connection_id`。麦克风要求见 [`https-webrtc.md`](https-webrtc.md)。
- 详细鉴权、工具契约、反代排障见 [`mcp.md`](mcp.md) 与 [`mcp-validation.md`](mcp-validation.md)。

## 5. WebDAV 备份

入口：**系统设置 → WebDAV 归档**，增量写入短信、通话记录与录音；远端失败只记录失败，不会误标为已归档。

- 填写地址、账号、密码与远端目录后先「测试连接」：会在远端写删临时文件，需要建目录 / 写 / 删除权限。地址必须是最终地址，PUT 不跟随重定向。
- 远端结构：`base_path/sms/YYYY/MM/sms-<id>.json`、`call/.../call-<id>.json`、`recording/.../<文件名>.mp3`。
- 内容哈希与路径记录在本地台账：未变化跳过，变化覆盖。
- 清理分「预览」和「确认」两步，令牌与数据范围绑定。仅删除同时满足：超过保留期、哈希与最后一次成功归档匹配、录音也已归档、该类型已启用。自动清理只在归档完整成功后运行。
- WebDAV 不是端到端加密；公网请务必使用受信任 HTTPS 并限制账号目录。

## 6. 故障排查速查

| 现象 | 处理 |
| --- | --- |
| 浏览器打不开 / 打开错误页面 | `start-vesim.bat` 会读取 `server.port`、延迟约 3 秒并打开 `http://localhost:<port>`（localhost 可触达 IPv6 监听）；默认 `http://localhost:7575` |
| 日志出现 `WebRTC UDP 端口 … 不可用，已改用 …` | 端口真实被占用，已自动顺延并继续启动；仅按**端口**比较，`:7580` → `[::]:7580` 不会误报 |
| HTTP / HTTPS 端口占用 | 启动会报「监听 HTTP/HTTPS …」明确错误，修改 `server.port` / `https_port` |
| `unable to open database file` | 2.1.37 起自动创建 `data/` 父目录；确认工作目录可写、未把 `data/` 设为只读 |
| Windows 设备列表为空 | 预期行为：QMI/MBIM/AT、VoWiFi、USB 音频、按网卡绑定、udev 热插拔都需要 Linux + USB 直通；模组部署请使用上游 HiDeck |
| MCP 返回 401 / 403 / 503 | 401：令牌无效或已撤销；403：反代 Host / 浏览器 Origin 不匹配（把上游 Host 设为 `localhost:7575`，勿关跨源校验）；503：开关未开或服务 / 存储异常 |
| HTTPS 语音无麦克风 | 浏览器仅允许 HTTPS 或 localhost 使用麦克风；反向代理场景需同时放行同号 TCP 与 UDP；内置 HTTPS 的私有 CA 需客户端信任 |
| 串口 / 模组识别异常 | 不要自行写 `QSIMDET`（会写 NV 且重启仍生效）；确认实际 AT 口与读写权限，见 [`modem-hardware.md`](modem-hardware.md) |
| EC20 在 macOS / OrbStack 报 QMI framing error / EOF | 设置 `VESIM_QMI_SERIAL_REQUESTS=1` 后重建服务，见 [`qmi-serial-requests.md`](qmi-serial-requests.md) |
| 旧版限流配置提示 | `server.sms_rate_limit_disabled` 已废弃，改用 `outbound_limits` |
