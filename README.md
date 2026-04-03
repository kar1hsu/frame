# Frame - Go 模块化后台管理框架

基于 **Gin + GORM + Casbin + JWT** 构建的模块化后台管理框架，内嵌 **Vue 3 + Element Plus** 前端管理面板。

## 技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| Web 框架 | Gin | 路由、中间件、请求处理 |
| ORM | GORM | 支持 MySQL / PostgreSQL，通过配置切换 |
| 鉴权 | golang-jwt | JWT Token 签发与验证 |
| 权限管理 | Casbin | RBAC 权限模型，支持路径参数匹配 |
| 配置管理 | Viper | YAML/ENV 多环境配置 |
| 日志 | Zap + Lumberjack | 结构化日志，自动轮转 |
| 缓存 | Redis (go-redis) | Token 黑名单、权限缓存、登录限流 |
| 任务队列 | Asynq | 分布式消息队列 + 定时任务，基于 Redis |
| 密码 | bcrypt | 密码哈希加密 |
| 前端 | Vue 3 + Element Plus + Vite | 管理面板，通过 Go embed 嵌入 |
| 部署 | Docker + docker-compose | 一键容器化部署 |

## 功能特性

- **模块化架构** - 通过 `Module` 接口扩展，自带 Admin 和 API 两个模块
- **JWT 鉴权** - Bearer Token 认证，支持退出登录 Token 黑名单
- **登录保护** - Redis 记录登录失败次数，5 次失败后锁定 15 分钟
- **RBAC 权限** - 基于 Casbin 的角色权限管理，支持 `keyMatch2` 路径参数匹配
- **按钮级权限** - 菜单关联 API，支持 list/query/add/edit/delete 细粒度控制
- **权限自动同步** - 分配菜单时自动生成对应的 Casbin API 策略，零手动配置
- **权限缓存** - Redis 缓存用户权限列表，角色变更时自动清除
- **用户管理** - 用户 CRUD、密码加密、角色分配
- **角色管理** - 角色 CRUD、菜单权限分配（含按钮粒度）
- **菜单管理** - 树形菜单管理（目录/菜单/按钮三级）、动态菜单
- **API 管理** - API 接口注册，与菜单/按钮关联，数据库驱动的权限配置
- **消息队列** - 基于 Asynq + Redis 的分布式任务队列，支持即时/延迟/唯一任务
- **定时任务** - Asynq Scheduler，cron 语法，多实例防重复投递
- **前端面板** - Element Plus 管理界面，通过 Go embed 内嵌到二进制
- **多数据库** - 通过配置切换 MySQL 或 PostgreSQL

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 20+（构建前端）
- MySQL 8.0+ 或 PostgreSQL 14+
- Redis 6+

### 1. 克隆项目

```bash
git clone <repo-url> frame
cd frame
```

### 2. 配置

编辑 `config/config.yaml`，修改数据库和 Redis 连接信息。

### 3. 构建前端

```bash
cd web/admin
npm install
npm run build
cd ../..
```

### 4. 启动服务

```bash
go run cmd/server/main.go
```

服务启动后：
- 后台管理面板: http://localhost:8080
- Admin API: http://localhost:8080/admin/*
- 对外 API: http://localhost:8080/api/*

### 5. 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 超级管理员 |

### 6. 前端开发模式

```bash
cd web/admin
npm run dev
```

Vite 开发服务器启动后访问 `http://localhost:5173`，API 自动代理到后端 `http://localhost:8080`。

## Docker 部署

```bash
cd deploy
docker-compose up -d
```

## 项目结构

```
frame/
├── cmd/
│   ├── server/main.go          # Web 服务入口
│   └── worker/main.go          # Worker 进程入口（消费者 + 定时任务）
├── config/
│   ├── config.yaml             # 应用配置
│   └── rbac_model.conf         # Casbin RBAC 模型
├── internal/
│   ├── app/                    # 应用初始化 (Config/Logger/DB/Redis/Casbin/Task)
│   ├── server/                 # HTTP 服务 & 路由注册 & 静态文件
│   ├── middleware/             # 中间件 (JWT/Casbin/CORS/Logger)
│   ├── model/                  # 数据模型 (User/Role/Menu/API)
│   ├── dao/                    # 数据访问层
│   ├── tasks/                  # 任务定义与注册（Handler + 定时任务）
│   ├── module/
│   │   ├── admin/              # Admin 后台模块
│   │   │   ├── handler/        # 请求处理 (Auth/User/Role/Menu/API)
│   │   │   ├── service/        # 业务逻辑
│   │   │   └── router.go       # Admin 路由注册
│   │   └── api/                # 对外 API 模块
│   └── pkg/                    # 内部公共包
│       ├── jwt/                # JWT 签发/解析
│       ├── cache/              # 缓存层 (Store接口/RedisStore/业务缓存)
│       ├── task/               # 任务系统 (Client/Worker/Scheduler/Manager)
│       ├── response/           # 统一响应格式
│       ├── errcode/            # 错误码
│       └── utils/              # 工具 (密码哈希/分页)
├── web/admin/                  # Vue 3 前端项目
├── embed.go                    # Go embed 嵌入前端
└── deploy/                     # Docker 部署文件
```

## 路由层级设计

| 层级 | 中间件 | 说明 | 示例 |
|------|--------|------|------|
| 公开 | 无 | 无需登录 | `POST /admin/login` |
| 已登录 | JWT | 登录即可访问（个人信息、下拉选项） | `GET /admin/profile`, `/permissions`, `/roles/all` |
| 管理权限 | JWT + Casbin | 需要 RBAC 授权的管理操作 | `GET /admin/users`, `POST /admin/users` |

## API 概览

### 认证

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| POST | /admin/login | 登录（含登录失败限流） | 无 |
| POST | /admin/logout | 退出登录（Token 加入黑名单） | JWT |
| GET | /admin/profile | 获取当前用户信息 | JWT |
| GET | /admin/permissions | 获取当前用户权限标识列表（带缓存） | JWT |

### 用户管理

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | /admin/users | 用户列表 | JWT + RBAC |
| POST | /admin/users | 创建用户 | JWT + RBAC |
| GET | /admin/users/:id | 用户详情 | JWT + RBAC |
| PUT | /admin/users/:id | 更新用户 | JWT + RBAC |
| DELETE | /admin/users/:id | 删除用户 | JWT + RBAC |

### 角色管理

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | /admin/roles | 角色列表（分页） | JWT + RBAC |
| GET | /admin/roles/all | 全部角色（下拉选项） | JWT |
| POST | /admin/roles | 创建角色 | JWT + RBAC |
| GET | /admin/roles/:id | 角色详情 | JWT + RBAC |
| PUT | /admin/roles/:id | 更新角色 | JWT + RBAC |
| DELETE | /admin/roles/:id | 删除角色 | JWT + RBAC |
| PUT | /admin/roles/:id/menus | 分配菜单（自动同步 Casbin 策略） | JWT + RBAC |

### 菜单管理

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | /admin/menus/tree | 完整菜单树 | JWT |
| GET | /admin/menus/user | 当前用户菜单树 | JWT |
| POST | /admin/menus | 创建菜单（支持关联 API） | JWT + RBAC |
| GET | /admin/menus/:id | 菜单详情（含关联的 API） | JWT + RBAC |
| PUT | /admin/menus/:id | 更新菜单 | JWT + RBAC |
| DELETE | /admin/menus/:id | 删除菜单 | JWT + RBAC |

### API 管理

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | /admin/apis/all | 全部 API（下拉选项） | JWT |
| GET | /admin/apis | API 列表（分页） | JWT + RBAC |
| POST | /admin/apis | 创建 API | JWT + RBAC |
| PUT | /admin/apis/:id | 更新 API | JWT + RBAC |
| DELETE | /admin/apis/:id | 删除 API | JWT + RBAC |

## 权限系统

### 菜单类型

| 类型 | type 值 | 说明 |
|------|---------|------|
| 目录 | 0 | 菜单分组，如"系统管理" |
| 菜单 | 1 | 页面入口，如"用户管理" |
| 按钮 | 2 | 操作权限，如"新增用户"、"删除用户" |

### 权限标识命名规范

```
模块:资源:操作
```

每个菜单下的标准操作：

| 操作 | 标识示例 | 含义 |
|------|----------|------|
| list | system:user:list | 查看列表（菜单级，控制侧边栏显示） |
| query | system:user:query | 查看详情（只读查看单条记录） |
| add | system:user:add | 新增 |
| edit | system:user:edit | 编辑 |
| delete | system:user:delete | 删除 |

### 菜单字段填写规则

| 字段 | 目录 | 菜单 | 按钮 |
|------|------|------|------|
| 路由路径 | `/模块` | `/模块/资源` | 留空 |
| 组件路径 | 留空 | `模块/资源/index` | 留空 |
| 权限标识 | 留空 | `模块:资源:list` | `模块:资源:操作` |
| 图标 | 填写 | 填写 | 留空 |
| 关联API | 无 | 列表接口 | 对应的接口 |

### 权限配置流程（纯后台操作，零代码）

1. **API 管理** — 注册新的 API 接口记录
2. **菜单管理** — 创建菜单和按钮子项，关联对应的 API
3. **角色管理** — 给角色分配菜单/按钮，系统自动生成 Casbin 策略

### Redis 缓存

#### 架构设计

缓存层通过接口抽象，业务代码不直接依赖 go-redis：

```
业务代码 (cache.BlacklistToken / cache.GetUserPermissions ...)
    │
    ▼
cache.Store 接口 (String/Hash/List/Set 全类型支持)
    │
    ▼
cache.RedisStore 实现 (封装 go-redis，自动处理 key 前缀)
```

#### 全局 Key 前缀

通过 `config.yaml` 配置，多项目共用 Redis 时不会冲突：

```yaml
redis:
  key_prefix: "frame:"
```

实际存储的 key 示例：`frame:token:blacklist:eyJhb...`、`frame:perm:user:1`

#### Store 接口支持的数据类型

| 类型 | 方法 |
|------|------|
| String | `Get` `Set` `Del` `Exists` `Incr` `Decr` `Expire` `TTL` `Scan` |
| Hash | `HGet` `HSet` `HDel` `HGetAll` `HExists` `HIncrBy` `HKeys` `HLen` `HMGet` |
| List | `LPush` `RPush` `LPop` `RPop` `LRange` `LLen` `LRem` `LIndex` `LTrim` |
| Set | `SAdd` `SRem` `SMembers` `SIsMember` `SCard` |

#### 内置缓存功能

| 功能 | Key 格式 | TTL | 说明 |
|------|----------|-----|------|
| Token 黑名单 | `token:blacklist:{token}` | 与 JWT 过期时间一致 | 退出登录后 Token 立即失效 |
| 权限缓存 | `perm:user:{userID}` | 10 分钟 | 减少权限查询的数据库压力 |
| 登录限流 | `login:fail:{username}` | 15 分钟 | 5 次失败后锁定 |

#### 业务中使用缓存

```go
// 直接调用封装好的业务方法
cache.BlacklistToken(token, expiration)
cache.IsTokenBlacklisted(token)
cache.SetUserPermissions(userID, perms)

// 或通过 Store 接口使用任意 Redis 操作
store := cache.GetStore()
store.HSet("user:profile:1", "name", "张三", "age", "25")
store.LPush("task:queue", taskJSON)
store.SAdd("online:users", "user_1")
```

#### 自定义实现

测试或切换缓存方案时，实现 `cache.Store` 接口即可：

```go
cache.InitStore(myMemoryStore)  // 替换为内存实现
cache.InitStore(myRedisCluster) // 替换为集群实现
```

## 任务系统（消息队列 + 定时任务）

基于 Asynq + Redis，支持分布式部署，多 Worker 实例自动负载均衡。

### 架构

```
Web 服务 (生产者)                    Worker 进程 (消费者)
cmd/server/main.go                  cmd/worker/main.go
  │ app.TaskMgr.Client.Enqueue()       │ tasks.RegisterHandlers()
  ▼                                    ▼
┌──────────────────────────────────────────┐
│                  Redis                   │
│  队列: critical / default / low          │
│  Scheduler: cron 定时投递（内置分布式锁） │
└──────────────────────────────────────────┘
```

### 启动 Worker

Web 服务和 Worker 是独立进程，可以分开部署：

```bash
# 终端 1: Web 服务（生产者）
go run cmd/server/main.go

# 终端 2: Worker（消费者 + 定时任务）
go run cmd/worker/main.go
```

多实例部署时，启动多个 Worker 即可水平扩展。

### 投递任务（生产者）

在任意 Handler / Service 中调用：

```go
// 即时任务
app.TaskMgr.Client.Enqueue("email:send", EmailPayload{To: "user@example.com", Subject: "Welcome"})

// 延迟任务（10 分钟后执行）
app.TaskMgr.Client.EnqueueDelay("email:send", payload, 10*time.Minute)

// 去重任务（1 小时内同样的任务只投递一次）
app.TaskMgr.Client.EnqueueUnique("report:generate", payload, 1*time.Hour)

// 指定队列（高优先级）
app.TaskMgr.Client.EnqueueToQueue("order:notify", payload, "critical")
```

### 定义任务处理器（消费者）

在 `internal/tasks/` 中创建：

```go
// internal/tasks/types.go — 定义任务类型名
const TypeOrderNotify = "order:notify"

// internal/tasks/order.go — 实现处理逻辑
func HandleOrderNotify(ctx context.Context, payload []byte) error {
    var p OrderPayload
    json.Unmarshal(payload, &p)
    // 处理逻辑...
    return nil
}

// internal/tasks/register.go — 注册
func RegisterHandlers(w *task.Worker) {
    w.Handle(TypeOrderNotify, HandleOrderNotify)
}
```

### 定时任务

在 `internal/tasks/register.go` 中注册：

```go
func RegisterCronJobs(s *task.Scheduler) {
    // 每天凌晨 2 点清理
    s.Register(task.CronTask{Cron: "0 2 * * *", TypeName: TypeCleanup})
    // 每 5 分钟执行
    s.Register(task.CronTask{Cron: "@every 5m", TypeName: TypeSyncData})
    // 指定队列
    s.Register(task.CronTask{Cron: "0 8 * * 1", TypeName: TypeWeeklyReport, Queue: "low"})
}
```

### 队列优先级

在 `config.yaml` 中配置，排在前面的优先级更高：

```yaml
task:
  concurrency: 10
  queues:
    - critical    # 权重 3（最高）
    - default     # 权重 2
    - low         # 权重 1（最低）
```

## 扩展模块

实现 `Module` 接口即可添加新模块：

```go
type Module interface {
    Name() string
    RegisterRoutes(rg *gin.RouterGroup)
}
```

在 `main.go` 中注册：

```go
router := server.NewRouter(
    frame.AdminDist,
    admin.New(),
    api.New(),
    yourmodule.New(), // 新模块
)
```

## License

MIT
