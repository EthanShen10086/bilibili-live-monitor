# B站开播提醒：Node、Go、tRPC-Go

监控直播间 1616，默认北京时间周三、周五、周六、周日 18:00–24:00 每 1 分钟轮询，通过飞书签名群机器人提醒。每场直播去重，不限制每周通知次数。

| 实现 | 目录 | 用途 |
|---|---|---|
| Node.js 24 + TypeScript | `node/` | 已有默认实现 |
| Go 1.26.3 | `go/cmd/monitor` | 推荐的新部署选择，纯 Go SQLite、单文件运行 |
| tRPC-Go 1.1.0 | `go/cmd/trpc-monitor` | 复用 Go 业务，使用真实 tRPC-Go 框架提供本机状态服务 |

三种实现选一种运行。相同服务名、SQLite 格式和锁协议便于迁移；不同目录不会自动共享去重状态，迁移须停止源实例并复制状态。tRPC-Go 不是包装名称：实际创建框架 Server、注册 HTTP Service，并协调后台任务退出。

- [Mac 配置、消息测试和版本迁移](docs/MAC.md)
- [Linux 云服务器部署、恢复与切换](docs/CLOUD.md)
- [实现原理、配置与依赖](docs/ARCHITECTURE.md)
- [测试和验收边界](docs/VALIDATION.md)
- [GitHub 发布与排除文件](docs/PUBLISH.md)

凭证放在各部署目录 `.env`，权限 600，绝不提交。仓库保留 `package-lock.json`、`go.mod`、`go.sum`，排除 node_modules、Go 缓存、二进制、数据库和日志。第三方 Node 长连接源码及其许可证位于 `node/vendor/`，是构建所需源码，并非安装产生的依赖目录。

官方事件与飞书私聊可通过 YAML 配置启用；需要各平台凭证和目标房间授权。没有授权时不能保证订阅任意房间，也不会自动降级为其他接入方式。

轮询间隔设置：[分钟级配置和重启步骤](docs/POLLING_INTERVAL.md)。
