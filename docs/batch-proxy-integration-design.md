# Batch 代理分配对接设计文档

这份文档说明：如果你的业务系统希望“每 50 个用户分配一个代理 URL，并且每个组使用的节点隔离”，应该如何和当前 Resin 服务对接。

## 当前服务信息

- 局域网 IP：`192.168.20.204`
- 代理端口：`12260`
- 代理 Token：`zxc13875517127`

统一入口格式：

```text
http://平台名:zxc13875517127@192.168.20.204:12260
```

## 1. 目标

你的业务侧有一个配置项：
- `proxy_url`

用户支持填写：
- `http://...`
- `https://...`
- `socks5://...`

你的目标是：
- 每 50 个用户共用一个代理 URL
- 不同组之间尽量不要混用节点
- 希望每组是隔离的固定代理池

## 2. 推荐方案

推荐使用 Resin 的 **固定节点平台**，而不是 `platform.batch001` 这种 account 分组。

原因：
- account 分组只是逻辑身份
- 不能保证不同 batch 不会混到同一个节点
- 你真正想要的是“组和组隔离”

所以正确做法是：
- 一个 batch 对应一个固定节点平台

## 3. 当前已经配置好的固定节点平台

### clean 固定节点平台
- `clean01` → 美西直连-80
- `clean02` → 美西直连-74
- `clean03` → 美西直连-74tunnel
- `clean04` → 美西家宽-192

### mmxz 固定节点平台
- `mmxz01` → v3【hk】香港1-nf
- `mmxz02` → 【直连】日本gmo原生
- `mmxz03` → v4【a】新加坡pccw
- `mmxz04` → 【EU】德国
- `mmxz05` → v3【a】新加坡-01
- `mmxz06` → v3【a】新加坡-02
- `mmxz07` → v3【a】香港-02
- `mmxz08` → v3【hk】香港2-nf
- `mmxz09` → v3【hk】香港G口
- `mmxz10` → v3【hk】香港G口-1
- `mmxz11` → v3【hk】香港G口-2
- `mmxz12` → v3【jp】日本2G-02
- `mmxz13` → v3【sg】新加坡-1
- `mmxz14` → v3【sg】新加坡-2
- `mmxz15` → v3【sg】新加坡-3
- `mmxz16` → v4【a】台湾
- `mmxz17` → v4【sg】新加坡-3
- `mmxz18` → 【EU】土耳其
- `mmxz19` → 【EU】英国伦敦原生
- `mmxz20` → 【a】香港hkt
- `mmxz21` → 【chatgpt】美国原生-01
- `mmxz22` → 【chatgpt】美国原生-02
- `mmxz23` → 【chatgpt】美国原生-03
- `mmxz24` → 【直连】新加坡-01
- `mmxz25` → 【直连】新加坡原生
- `mmxz26` → 【直连】美国cera-02
- `mmxz27` → 【直连】香港
- `mmxz28` → 【直连】香港hkt-NF

### 公共池（不隔离）
- `clean`
- `mmxz`

## 4. 代理 URL 模板

统一入口：

```text
192.168.20.204:12260
```

统一 token：

```text
zxc13875517127
```

### clean 固定节点池

```text
http://clean01:zxc13875517127@192.168.20.204:12260
http://clean02:zxc13875517127@192.168.20.204:12260
http://clean03:zxc13875517127@192.168.20.204:12260
http://clean04:zxc13875517127@192.168.20.204:12260
```

### mmxz 固定节点池

```text
http://mmxz01:zxc13875517127@192.168.20.204:12260
http://mmxz02:zxc13875517127@192.168.20.204:12260
http://mmxz03:zxc13875517127@192.168.20.204:12260
http://mmxz04:zxc13875517127@192.168.20.204:12260
```

## 5. 用户分配逻辑

推荐分配方式：

### clean
- 用户 1-50 → `clean01`
- 用户 51-100 → `clean02`
- 用户 101-150 → `clean03`
- 用户 151-200 → `clean04`

### mmxz
- 用户 1-50 → `mmxz01`
- 用户 51-100 → `mmxz02`
- 用户 101-150 → `mmxz03`
- 用户 151-200 → `mmxz04`
- 之后继续顺延到 `mmxz05` ~ `mmxz28`

## 6. 为什么这样才符合你的诉求

你的诉求不是“给不同组起不同名字”，而是：

- 每组真正使用不同的节点
- 不同组之间不要混流量
- 避免全部挤在同一个节点上

如果只是用：

```text
clean.batch001
clean.batch002
```

那只是不同 account，不能保证 Resin 不把它们路由到同一个 clean 节点池里的节点。

而现在这种固定平台方案：
- `clean01` 只包含一个节点
- `clean02` 只包含一个节点
- `clean03` 只包含一个节点
- `clean04` 只包含一个节点

所以组和组之间天然隔离。

## 7. 业务系统应该存什么

建议你的业务系统维护：

| group_id | pool_name | user_count | proxy_url | status |
|---|---|---:|---|---|
| group001 | clean01 | 50 | `http://clean01:...` | active |
| group002 | clean02 | 50 | `http://clean02:...` | active |
| group003 | clean03 | 13 | `http://clean03:...` | active |

用户表：

| user_id | group_id | proxy_url |
|---|---|---|
| u001 | group001 | `http://clean01:...` |
| u051 | group002 | `http://clean02:...` |

## 8. 分配逻辑伪代码

```python
def allocate_proxy(mode: str = "clean"):
    pool_names = {
        "clean": ["clean01", "clean02", "clean03", "clean04"],
        "mmxz": ["mmxz01", "mmxz02", "mmxz03", "mmxz04", "mmxz05", "mmxz06", "mmxz07", "mmxz08", "mmxz09", "mmxz10", "mmxz11", "mmxz12", "mmxz13", "mmxz14", "mmxz15", "mmxz16", "mmxz17", "mmxz18", "mmxz19", "mmxz20", "mmxz21", "mmxz22", "mmxz23", "mmxz24", "mmxz25", "mmxz26", "mmxz27", "mmxz28"],
    }[mode]

    group = find_first_group_with_capacity(pool_names, max_users=50)
    if not group:
        raise Exception("all isolated pools are full")

    group.user_count += 1
    save(group)

    return f"http://{group.pool_name}:zxc13875517127@192.168.20.204:12260"
```

## 9. Resin 提供的接口能帮你做什么

### 查看所有平台
```http
GET /api/v1/platforms
```

用途：
- 看当前有哪些平台
- 看每个平台的 `routable_node_count`

### 查看某个平台下 lease
```http
GET /api/v1/platforms/{id}/leases
```

用途：
- 看这个固定平台当前是否有活跃 account
- 用于运营观察

### 查看某个平台 IP 负载
```http
GET /api/v1/platforms/{id}/ip-load
```

用途：
- 看固定池当前出口负载情况

## 10. 固定池不可用时的处理策略（推荐方案 C）

如果某个固定节点平台不可用，推荐不要把整组用户永久卡死在这个池上，而是采用：

### 方案 C：动态挑可用固定池

规则：
- 每个组默认绑定一个固定池
- 如果这个固定池当前不可用，就从同类型的“当前可用固定池列表”里重新挑一个可用池
- 业务系统更新这个组的 `pool_name` 和 `proxy_url`
- 新的代理地址发给这个组

### 例子

原本：
- `group017` → `mmxz05`

如果 `mmxz05` 当前：
- `routable_node_count = 0`

那么业务系统可以自动改成：
- `group017` → `mmxz11`

然后给这个组新的 URL：

```text
http://mmxz11:zxc13875517127@192.168.20.204:12260
```

### 为什么推荐方案 C

优点：
- 可用性最高
- 某个固定池坏了不会导致整组直接不可用
- 仍然保留“组级隔离”，只是这个组会迁移到另一个可用固定池

缺点：
- 组和具体节点的绑定不是永久不变的
- 业务系统需要维护当前组映射到哪个池

### 如何判断某个池是否不可用

最直接方式：

```http
GET /api/v1/platforms
```

看平台字段：
- `routable_node_count`

判断规则：
- `routable_node_count > 0` → 当前可用
- `routable_node_count == 0` → 当前不可用

### 推荐的业务系统逻辑

建议维护：

| group_id | pool_name | user_count | proxy_url | status |
|---|---|---:|---|---|
| group001 | clean01 | 50 | `http://clean01:...` | active |
| group002 | clean02 | 50 | `http://clean02:...` | active |
| group017 | mmxz11 | 32 | `http://mmxz11:...` | migrated |

定时任务：
1. 定时查询 `GET /api/v1/platforms`
2. 找出 `routable_node_count == 0` 的平台
3. 如果某组绑定的平台不可用：
   - 在同类型平台中找下一个可用池
   - 更新该组 `pool_name`
   - 重新生成 `proxy_url`

### 伪代码

```python
def choose_available_pool(mode: str, current_pool: str | None = None):
    pools = get_platforms_from_resin()
    candidates = [p for p in pools if p['name'].startswith(mode) and p['routable_node_count'] > 0]

    if current_pool and any(p['name'] == current_pool for p in candidates):
        return current_pool

    if not candidates:
        raise Exception(f'no available {mode} fixed pools')

    return candidates[0]['name']


def ensure_group_pool(group, mode='mmxz'):
    pool = choose_available_pool(mode, group.pool_name)
    if pool != group.pool_name:
        group.pool_name = pool
        group.proxy_url = f"http://{pool}:zxc13875517127@192.168.20.204:12260"
        group.status = 'migrated'
        save(group)
```

## 11. 最终推荐结论

如果你的真实诉求是：
- 每 50 人一个组
- 每组尽量固定且隔离
- 避免所有人挤在同一个节点
- 某个固定池挂了还能自动迁移到可用池

那就应该：

- 不要再用 `clean.batch001` 这种 account 分组方案
- 直接用已经配置好的固定节点平台：
  - `clean01`~`clean04`
  - `mmxz01`~`mmxz28`
- 业务系统采用方案 C：动态挑可用固定池

这样最符合你的业务目标，也最容易长期维护。
