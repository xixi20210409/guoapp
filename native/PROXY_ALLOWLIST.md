# duanju_core_server — 按主机代理 allowlist

盒子上插件（`plugin.video.hongguo`）走 duanju_core_server 这个 Go 核心取流。
部分上游源（如花果/观果/网果一类）有 IP 限流或 Cloudflare 挑战，
需要走代理才能稳定访问；但主源红果若也走代理会被拖慢（实测目录拉取飙到 60s+）。

为此核心支持「按主机代理」：

- 环境变量 `DUANJU_CORE_PROXY`：手动代理地址（如 `http://192.168.68.3:7893`）。
- 环境变量 `DUANJU_CORE_PROXY_HOSTS`：逗号分隔的「允许走代理」主机后缀白名单。
  仅当请求主机匹配白名单时才走 `DUANJU_CORE_PROXY`；其余（含红果、localhost）一律直连。

实现位置：

- `core/app_settings.go`：`proxyRouter.allow` 字段 + `proxy()` 方法的后缀匹配逻辑。
- `core/app_runtime.go`：`newNativeEngine` 读取 `DUANJU_CORE_PROXY_HOSTS` 注入到 router.allow。

插件侧（`resources/lib/zjcore.py`）在拉起核心时把代理地址与
`PROXY_HOSTS_ALLOWLIST`（8 个 duanju 源的 Base URL 域名）注入上述环境变量。

> 构建开关：`go build -ldflags="-X duanjuapp/native/core.buildAllSources=true"`
> 放开全站源；CI（`.github/workflows/build_core.yml`）每次 `native/**` 变更时
> 交叉编译 linux/arm64 并把二进制提交到 `core_builds/`，插件 `updater.py` 据此自动更新盒子核心。
