# Linux 云服务器部署与运维

## 1. 资源和账号

准备支持 systemd 的 Linux、普通专用账号 `live-monitor`、SSH 密钥及出站访问 B站和飞书。建议 1 vCPU / 1 GB 起步。业务没有公网入站端口需求；tRPC 状态接口只监听回环地址。以下 sudo 操作由服务器管理员首次执行：

```sh
sudo useradd -m -s /bin/bash live-monitor
sudo mkdir -p /opt/live-monitor
sudo chown live-monitor:live-monitor /opt/live-monitor
sudo loginctl enable-linger live-monitor
```

配置该账号 authorized_keys，Mac 的 `~/.ssh/config` 添加：

```text
Host live-monitor
  HostName 服务器地址
  User live-monitor
  IdentityFile ~/.ssh/你的密钥
```

先验证 `ssh live-monitor true`，服务全部用该账号安装，禁止 sudo 执行业务命令。linger 让用户 systemd 在退出登录和重启后继续运行。

## 2. 构建并上传

推荐 Mac 交叉编译，避免小云主机编译 SQLite 消耗内存：

```sh
cd /你的路径/bilibili-live-monitor/go
./bin/setup
mkdir -p dist/linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o dist/linux-amd64/monitor ./cmd/monitor
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o dist/linux-amd64/trpc-monitor ./cmd/trpc-monitor
scp dist/linux-amd64/monitor dist/linux-amd64/trpc-monitor live-monitor:/opt/live-monitor/
scp config.yaml trpc_go.yaml live-monitor:/opt/live-monitor/
```

ARM 云主机改 GOARCH=arm64。服务器把文件放成管理命令约定的目录结构：

```sh
ssh live-monitor
cd /opt/live-monitor
mkdir -p dist bin var
mv monitor trpc-monitor dist/
chmod 700 var
# 从仓库 go/bin/ 上传 monitor 和 trpc-monitor 两个包装脚本至这里的 bin/
chmod 755 bin/monitor bin/trpc-monitor dist/monitor dist/trpc-monitor
```

也可以直接将整个仓库 clone 到服务器，进入 go/ 用 `./bin/setup` 构建；此时 install_dir 必须指向实际 go 目录。

## 3. 云端凭证与首次准备

在服务器安装目录独立创建 `.env`（内容参照 go/.env.example），chmod 600。凭证不用 scp 从 Mac 自动传输。配置里 cloud.install_dir 和 Mac 保持相同，云端设为 cloud：

```sh
./bin/monitor set-active cloud
./bin/monitor check-config --probe
./bin/monitor test-notification
./bin/monitor service install --side cloud
./bin/monitor service doctor --side cloud
```

此时先不启动，确保 Mac 仍是唯一监控实例。选择 tRPC-Go 时用 bin/trpc-monitor 执行 install，它写入该入口；远程管理 bin/monitor 仍须保留，切换会复用已安装的 unit。

## 4. 从 Mac 切换

Mac 也必须已安装本版本对应 launchd 服务，并批准本次开机启动。在 Mac 配置 ssh_host=live-monitor、install_dir=/opt/live-monitor。执行：

```sh
./bin/monitor switch cloud
./bin/monitor status
# 切回：
./bin/monitor confirm-start
./bin/monitor switch local
```

切换先预检目标，再禁用停止两端，确认旧进程退出后同步 SQLite 和非敏感 YAML。成功后更新 active，只有目标启用自动启动。无法证明旧实例停止会拒绝开始。目标启动失败先确认停止并带回最新队列再恢复源。不要手动同时 start 两端，分布式双机不是共享文件锁。

## 5. 自动恢复与验收

云端 unit 使用 Restart=always、20 秒等待、45 秒退出超时、UMask0077；enable + linger 实现服务器重启后的无人值守启动，不使用 Mac 的图形确认机制。

```sh
./bin/monitor service health --side cloud
./bin/monitor service verify-recovery --side cloud
systemctl --user status live-monitor.service
journalctl --user -u live-monitor.service -n 100
```

服务实际 unit 名以 doctor 和安装输出为准。verify-recovery 会杀掉当前受管理进程，验证新 PID 和检测器恢复；须在可接受短暂停顿时执行。服务器 reboot 由管理员操作，随后检查开机后自动恢复、群消息和两端唯一运行。尚无云主机资源时，这些命令与源码测试不能代替真实云端验收。

更新二进制前 stop，保留数据库和 .env，替换后 start/health。tRPC 的 19028/19029 不打开安全组，可用 SSH 端口转发访问。Node 云端详细手册另见 node/CLOUD_DEPLOYMENT.md。
