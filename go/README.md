# Go 与 tRPC-Go

```sh
./bin/setup
cp .env.example .env
chmod 600 .env
# 填写凭证后：
./bin/monitor check-config --probe
./bin/monitor test-notification
# tRPC-Go 入口：
./bin/trpc-monitor run
```

不要同时运行两个入口，二者共享 SQLite、配置和实例锁。常驻安装参见 [Mac 手册](../docs/MAC.md)、[云端手册](../docs/CLOUD.md)。测试见 [验收记录](../docs/VALIDATION.md)。无需手动 vendor 模块，go.mod/go.sum 记录依赖。

若默认模块代理下载缓慢，可为本次构建设置 `GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct`，保持 Go checksum 校验开启。这里不修改全局 Go 设置。

轮询间隔设置：[分钟级配置](../docs/POLLING_INTERVAL.md)。
