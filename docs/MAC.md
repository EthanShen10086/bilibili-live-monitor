# Mac 复现手册

## 1. 选择实现并构建

现有 Node 后台服务无需重装。先测试 Go 单次命令，确认后再迁移常驻服务。

```sh
cd /你的路径/bilibili-live-monitor/go
# 安装 Go 1.26.3 或兼容的更高版本后：
./bin/setup
cp .env.example .env
chmod 600 .env
```

编辑 `.env`：填写飞书群自定义机器人 Webhook 和签名密钥。创建群机器人时开启签名校验，在飞书群设置中添加自定义机器人；Webhook 只能写入本地文件，不粘贴聊天、GitHub 或截图。配置 `config.yaml`，默认 polling + feishu_group。

```sh
./bin/monitor check-config --probe
./bin/monitor test-notification
```

第二条会真实发一条测试消息。检查群内消息和手机飞书通知；API 成功不能证明手机弹窗成功。tRPC-Go 使用 `./bin/trpc-monitor` 替换上述命令，通知方式相同。

## 2. 安装常驻服务

若 Node 已在运行，先在原 Node 目录执行 `./bin/monitor service stop --side local`，确认停止，再将其 `var/state.sqlite` 复制到 Go 目录 `var/`。实际数据库文件名以源目录为准；不要复制活动数据库。可以通过停止后的 `state-export` / `state-import` 命令迁移（参见下文）。`.env` 单独配置，两边环境变量名称相同。

```sh
./bin/monitor confirm-start
./bin/monitor service install --side local
./bin/monitor service start --side local
./bin/monitor service health --side local
./bin/monitor status
```

使用 tRPC-Go 常驻时，上述命令改为 `./bin/trpc-monitor`，install 会写入对应二进制的绝对路径。该版本额外读取 `trpc_go.yaml`，默认只绑定 127.0.0.1，状态 19029、管理 19028。

```sh
curl -f http://127.0.0.1:19029/healthz
curl http://127.0.0.1:19029/status
```

这两个接口只读；检测器健康或处于窗口外时 healthz 返回 200，等待确认或检测失败时返回 503。

## 3. 登录、重启确认与崩溃恢复

macOS launchd 在登录后启动。`deployment.local.confirm_each_boot: true`：每次电脑重启后的首次启动弹出确认；确认后保存本次开机标识，同一开机周期再次登录或崩溃恢复不会重复确认。选择暂不启用会保持等待确认，不执行检测；可手动 confirm-start 或停止服务。也可手动 `confirm-start` 批准本次启动。

```sh
./bin/monitor service doctor --side local
# 会主动终止本项目被管理的进程，再检查 PID 变化和健康；只在确认允许测试时执行：
./bin/monitor service verify-recovery --side local
```

跨重启确认需实际重启电脑验收：重启、登录、确认、查看 status 和群消息。不要把源码测试当成这项验收。监控时段需保持 Mac 唤醒和网络可用。

## 4. 停止与版本迁移

```sh
./bin/monitor service stop --side local
./bin/monitor state-export > /安全路径/state.base64
# 目标也必须停止，用目标版本命令导入：
./bin/trpc-monitor state-import < /安全路径/state.base64
./bin/trpc-monitor service install --side local
./bin/trpc-monitor service start --side local
```

导出的状态含直播记录和待发任务，属于本地运行数据，不上传仓库。保留源配置及备份，迁移失败时先确认目标已停止再恢复源服务。Go 原生与 tRPC-Go 共用同一 `go/var` 时不必导入，仍须先停再安装新入口。Node 详细操作保留在 `node/REPRODUCE_MAC.md`。
