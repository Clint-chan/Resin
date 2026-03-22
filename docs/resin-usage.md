# Resin 使用说明

## 1. 当前状态

本机 Resin 已启动：
- 管理地址：`http://127.0.0.1:12260`
- 管理 Token：`zxc13875517127`
- 代理 Token：`zxc13875517127`

已导入两个本地订阅：
- `clean-local-nodes`：4 个本地 clean 出口（通过宿主机 Xray 暴露给 Resin）
- `mmxz-local-nodes`：28 个 mmxz 节点

当前 clean 平台已有 4 个可路由节点；mmxz 这批已有 18 个可路由节点，已探测出 9 个健康出口 IP。

## 2. 目录里的关键文件

项目根目录下：
- `.env`：运行参数
- `docker-compose.yml`：容器启动配置
- `docs/resin-usage.md`：这份说明

你自己的节点来源：
- `~/.config/xray/config_3.json`
- `~/.config/xray/config_4.json`
- `~/.config/xray/config_5.json`
- `~/.config/xray/mmxz/订阅地址.txt`

## 3. 启动和停止

启动：

```bash
docker compose up -d
```

查看状态：

```bash
docker compose ps
```

查看日志：

```bash
docker compose logs -f resin
```

停止：

```bash
docker compose down
```

## 4. 打开管理界面

浏览器访问：

```text
http://127.0.0.1:12260
```

登录时使用：
- Password / Admin Token：`zxc13875517127`

## 5. 这个项目能怎么用

Resin 本质上是一个代理池网关。

它能帮你：
- 把多个节点统一收口到一个入口
- 自动做节点探测、延迟测试、故障切换
- 按节点标签、地区、订阅来源来筛选节点
- 给业务提供固定代理入口，而不是手动切节点
- 做 sticky session（按账号尽量绑定到同一个出口 IP）

对你现在这个场景，最直接的用途就是：
- 把 `clean-local-nodes` 作为“干净节点池”
- 把 `mmxz-local-nodes` 作为“机场节点池”
- 后面按不同业务分平台使用

## 6. 你现在已经导入的节点

### 6.1 干净节点

来源：
- `~/.config/xray/config.json`
- `~/.config/xray/config_3.json`
- `~/.config/xray/config_4.json`
- `~/.config/xray/config_5.json`

当前不是让 Resin 直接吃这几个 Reality 节点，而是：
- 先由宿主机上的 Xray 分别起 4 个本地代理出口
- 再让 Docker 里的 Resin 通过 `host.docker.internal` 接入这些本地 socks5 出口

当前 `clean-local-nodes` 里接入的是：
- `socks5://host.docker.internal:10809#美西直连-80`
- `socks5://host.docker.internal:10810#美西直连-74`
- `socks5://host.docker.internal:10811#美西直连-74tunnel`
- `socks5://host.docker.internal:10812#美西家宽-192`

说明：
- 这 4 个节点已经进入 Resin 的 clean 池
- 当前 `clean` 平台已有 4 个可路由节点
- 这是目前最稳定的接法，因为你本地 Xray 已验证可用，而 Resin 容器直接访问宿主机 Xray 出口即可

### 6.2 mmxz 节点

来源：`~/.config/xray/mmxz/订阅地址.txt`

我已经把它导入到：
- `mmxz-local-nodes`

说明：
- 共导入 28 个节点
- 当前已有 18 个可路由节点
- 已探测到 9 个健康出口 IP
- 已识别出多个地区出口，例如 `hk`、`jp`、`sg`、`us`、`de`、`tr`

## 7. 最常用的使用方式

### 7.1 当普通 HTTP 代理使用

```bash
curl -x http://127.0.0.1:12260 -U ":zxc13875517127" https://api.ipify.org
```

说明：
- 代理地址就是 `127.0.0.1:12260`
- 用户名可以留空
- 密码就是 `RESIN_PROXY_TOKEN`

### 7.2 指定平台使用

后面你在 Web UI 里建好平台后，可以这样用：

```bash
curl -x http://127.0.0.1:12260 -U "clean:zxc13875517127" https://api.ipify.org
```

或者：

```bash
curl -x http://127.0.0.1:12260 -U "mmxz:zxc13875517127" https://api.ipify.org
```

这里的 `clean` / `mmxz` 是你自己在 UI 里创建的平台名。

## 8. 推荐你下一步怎么配

我已经帮你建好了两个 Platform：
- `clean`
- `mmxz`

如果后面你想细分，也可以继续在 UI 里新增平台。


### Platform 1：`clean`
用途：专门走你自己的干净节点

当前状态：
- 已绑定到 4 个宿主机 Xray 出口
- 通过 Resin 中的 `clean-local-nodes` 进入 clean 池

建议做法：
- 关键业务优先走这里
- 如果你后面替换本地 Xray 出口配置，只要保持端口不变，Resin 这边不用改

### Platform 2：`mmxz`
用途：机场通用节点池

建议做法：
- 用地区、标签、名字筛掉你不想要的节点
- 比如只保留 `hk` / `jp` / `sg` / `us`

## 9. 如果你要新增节点

你后面有新节点时，有两种最省事方式：

### 方式 A：UI 里直接新增本地订阅
进入 `Subscriptions`，新增：
- `source_type = local`
- 把节点内容直接贴进去

支持的内容格式包括：
- `vmess://`
- `vless://`
- `trojan://`
- `ss://`
- sing-box JSON
- Clash YAML / JSON

### 方式 B：直接替换现有本地订阅内容
如果你只是更新同一批节点，可以把新的订阅内容贴回原订阅里，然后刷新。

## 10. 常见排查

### 节点导入了但不可用
先看：
- 节点本身是否还活着
- 服务器是否能连通目标节点
- Reality / VMess / VLESS 参数是否已变更

### 页面打不开
检查：

```bash
docker compose ps
```

再看日志：

```bash
docker compose logs -f resin
```

### 想重新导入
可以在 UI 里删掉旧订阅后重新建，或者更新已有订阅内容再刷新。

## 11. 你现在可以直接做的事

最推荐先试这几步：

1. 打开 `http://127.0.0.1:12260`
2. 看 `Subscriptions` 里是不是已经有：
   - `clean-local-nodes`
   - `mmxz-local-nodes`
3. 再去 `Nodes` 页面看哪些节点健康
4. 查看我已建好的平台：`clean` 和 `mmxz`
5. 用 curl 测试代理出口 IP

## 12. 一条最简单的测试命令

```bash
curl -x http://127.0.0.1:12260 -U ":zxc13875517127" https://api.ipify.org
```

如果你想指定某个平台：

```bash
curl -x http://127.0.0.1:12260 -U "mmxz:zxc13875517127" https://api.ipify.org
```

clean 平台测试：

```bash
curl -x http://127.0.0.1:12260 -U "clean:zxc13875517127" https://api.ipify.org
```
