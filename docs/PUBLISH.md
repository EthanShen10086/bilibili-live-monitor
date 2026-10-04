# 发布到公开 GitHub 仓库

目标：EthanShen10086/bilibili-live-monitor，公开。上传包含 node/、go/、手册、测试与依赖锁文件；不包含安装依赖、密钥和运行数据。

`node_modules` 不上传。Go 不上传 GOPATH、模块下载缓存、构建缓存和 dist 二进制，保留 go.mod/go.sum。Node 的 vendor/bilibili-live-ws 是有许可证的源码构建依赖，需要保留。`.env.example` 是空值示例可以上传。

发布脚本检查已提交的实际清单及常见密钥格式，要求工作区干净、登录账号正确、无现有 remote，并创建新公开仓库。已有同名仓库时停止，不覆盖。自动工具无法读 Mac 钥匙串时，在已经 gh auth status 登录成功的普通终端执行：

```sh
cd /你的路径/bilibili-live-monitor
./scripts/publish.sh
```

不需发送 token，不需再次登录。脚本复用 gh keyring，并为 HTTPS Git 配置 gh 凭证助手。成功后核实 visibility 为 PUBLIC、远端最新提交和三个版本目录。未执行成功前只能说本地交付完成，不能说已上传。
