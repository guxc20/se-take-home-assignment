# 烹饪机器人设计
## 功能点
* 支持“新建普通订单”、“新建VIP订单”
* 支持“+bot”、“-bot”
* 机器人一次处理1个订单，10秒完成
* 优先队列逻辑：VIP > 普通，且新 VIP 在现有 VIP 之后
* 减少机器人时，正在处理的订单回到待处理队列
* status 命令查看订单与机器人状态
* 并发安全

## 运行
```azure
./cookrobot
```
如果是其他架构的系统，请编译对应版本，比如linux/amd64的编译命令：
```azure
make build GOOS=linux GOARCH=amd64
```
## 测试
### 可用命令
- `new normal (nn)` → 新建普通订单
- `new vip (nv)` → 新建 VIP 订单
- `+bot (+b)` → 增加机器人
- `-bot (-b)` → 移除机器人
- `status (st)` → 查看当前订单和机器人状态
- `exit (ex)` → 退出程序

### 执行样例

```azure
=== 麦当劳烹饪机器人系统 ===
> 指令: new normal (nn) | new vip (nv) | +bot (+b) | -bot (-b) | status (st) | exit (ex)
> nn
[系统] 新订单 1 已加入 (VIP: false)

> nn
[系统] 新订单 2 已加入 (VIP: false)

> nn
[系统] 新订单 3 已加入 (VIP: false)

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 1 (VIP: false)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 已完成订单:
> 机器人状态:
> ====================

> nv
[系统] 新订单 4 已加入 (VIP: true)

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 4 (VIP: true)
> - 订单 1 (VIP: false)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 已完成订单:
> 机器人状态:
> ====================

> nv
[系统] 新订单 5 已加入 (VIP: true)

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 4 (VIP: true)
> - 订单 5 (VIP: true)
> - 订单 1 (VIP: false)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 已完成订单:
> 机器人状态:
> ====================

> +b
[系统] 添加机器人 1

> [Robot 1] 开始处理订单 4 (VIP: true)

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 5 (VIP: true)
> - 订单 1 (VIP: false)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 已完成订单:
> 机器人状态:
> - 机器人 1: 正在处理订单 4
> ====================

> +b
[系统] 添加机器人 2

> [Robot 2] 开始处理订单 5 (VIP: true)

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 1 (VIP: false)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 已完成订单:
> 机器人状态:
> - 机器人 1: 正在处理订单 4
> - 机器人 2: 正在处理订单 5
> ====================

> [Robot 1] 完成订单 4

> [Robot 1] 开始处理订单 1 (VIP: false)

> -b
[系统] 移除机器人 2

> [Robot 2] 停止处理订单 5

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 5 (VIP: true)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 已完成订单:
> - 订单 4 (VIP: true)
> 机器人状态:
> - 机器人 1: 正在处理订单 1
> ====================

> [Robot 1] 完成订单 1

> [Robot 1] 开始处理订单 5 (VIP: true)

> [Robot 1] 完成订单 5

> [Robot 1] 开始处理订单 2 (VIP: false)

> st
====== 当前状态 ======
> 待处理订单:
> - 订单 3 (VIP: false)
> 已完成订单:
> - 订单 4 (VIP: true)
> - 订单 1 (VIP: false)
> - 订单 5 (VIP: true)
> 机器人状态:
> - 机器人 1: 正在处理订单 2
> ====================

> [Robot 1] 完成订单 2

> [Robot 1] 开始处理订单 3 (VIP: false)

> [Robot 1] 完成订单 3

> st
====== 当前状态 ======
> 待处理订单:
> 已完成订单:
> - 订单 4 (VIP: true)
> - 订单 1 (VIP: false)
> - 订单 5 (VIP: true)
> - 订单 2 (VIP: false)
> - 订单 3 (VIP: false)
> 机器人状态:
> - 机器人 1: 空闲
> ====================

> ex

```
