# 招聘平台 2.0 — 规划与实施文档

**日期**: 2026-06-22
**版本**: 2.0
**当前状态**: 1.0 已完成（47 commits，GitHub 已推送）

---

## 一、1.0 交付回顾

| 模块 | 状态 | 核心能力 |
|------|------|----------|
| 基础设施 | ✅ | Go + Gin + PostgreSQL + Redis + JWT + COS + SMS |
| 用户与简历 | ✅ | 验证码登录、Profile CRUD、HR 认证、简历管理 |
| 职位系统 | ✅ | 发布/编辑/下架、PG 全文搜索、城市过滤 |
| 核心闭环 | ✅ | 投递、腾讯云 IM UserSig、面试邀约、地图导航 |
| 智能推荐 | ✅ | Jaccard 标签匹配、时间衰减排序 |
| Web 前端 | ✅ | React + Tailwind + Ant Design，3 套仪表盘 + 登录页 |
| 日志追踪 | ✅ | slog 结构化日志 + 请求级 Trace ID |
| 安全 | ✅ | config.yaml 脱敏、.gitignore 覆盖 |

**技术栈**: Go 1.22+ / Gin / PostgreSQL 16 / Redis 7 / React 19 / TypeScript / Tailwind CSS v4 / Ant Design 6

---

## 二、2.0 总览

```
Phase 1: 工程化（Docker + 种子数据 + API 文档）   预计 3 天
Phase 2: 质量（测试 + CI/CD + 错误码）             预计 3 天
Phase 3: 进阶（缓存 + 搜索 + 实时 + 移动端）       预计 2 周
```

---

## 三、Phase 1：工程化（立即可做）

### 3.1 Docker 一键部署 ⭐ 最高优先级

**目标**: 新机器上 `docker compose up` 一条命令启动全栈

**实施**:

```
项目根目录添加:
├── Dockerfile              # Go 后端多阶段构建
├── Dockerfile.frontend     # Nginx + 前端静态文件
├── docker-compose.yml      # PostgreSQL + Redis + Go + Nginx
├── nginx.conf              # 反向代理 /api → Go，/ → 前端
└── .env.docker             # Docker 环境变量
```

```dockerfile
# Dockerfile (Go)
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server/

FROM alpine:3.20
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/config/config.example.yaml /app/config/config.yaml
EXPOSE 8080
CMD ["/app/server"]
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: devuser
      POSTGRES_PASSWORD: devpass
      POSTGRES_DB: recruitment
    ports: ["5432:5432"]
    volumes: ["pgdata:/var/lib/postgresql/data"]

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]

  backend:
    build: .
    ports: ["8080:8080"]
    depends_on: [postgres, redis]
    environment:
      DB_HOST: postgres
      DB_USER: devuser
      DB_PASSWORD: devpass
      DB_NAME: recruitment
      REDIS_ADDR: redis:6379
      JWT_SECRET: docker-secret-change-me
    command: >
      sh -c "migrate -path /app/migrations -database postgres://... up && /app/server"

  frontend:
    build:
      context: ./web
      dockerfile: Dockerfile.frontend
    ports: ["80:80"]
    depends_on: [backend]

volumes:
  pgdata:
```

### 3.2 种子数据脚本

**目标**: `go run cmd/seed/main.go` 一键生成测试数据

**Mock 数据量**:
- 20 个职位（5 个公司 × 4 个岗位，分布在北上广深杭）
- 10 个求职者用户 + 5 份简历
- 3 个 HR 用户
- 5 条投递记录 + 2 条面试邀约

```
cmd/seed/main.go:
├── 创建 HR 用户（3人）
├── 创建求职者用户（10人）
├── 创建简历（5份，含技能标签）
├── 创建职位（20个，含全文搜索向量）
├── 创建投递记录（5条）
└── 创建面试邀约（2条）
```

### 3.3 Swagger API 文档

**目标**: 访问 `http://localhost:8080/swagger` 查看交互式文档

**实施**:
```bash
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/swag
```

在 `main.go` 顶部添加注释:
```go
// @title           放心招聘平台 API
// @version         2.0
// @description     放心招聘平台后端接口文档
// @host            localhost:8080
// @BasePath        /api
```

在每个 Handler 上添加 Swagger 注解:
```go
// @Summary      发送验证码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body  sendCodeReq  true  "手机号"
// @Success      200   {object}  gin.H
// @Router       /auth/send-code [post]
func (h *Handler) SendCode(c *gin.Context) { ... }
```

---

## 四、Phase 2：质量（迭代夯实）

### 4.1 测试体系

| 层级 | 工具 | 目标 |
|------|------|------|
| 单元测试 | Go `testing` + `testify` | 覆盖 Repository、Service 层 |
| 集成测试 | `httptest` | Mock DB，测 Handler |
| 前端测试 | Vitest + React Testing Library | 组件渲染 |

```
测试目标:
internal/auth/service_test.go      # 验证码 + 登录逻辑
internal/job/repository_test.go    # 搜索 + CRUD
internal/recommend/service_test.go # 推荐算法
web/src/__tests__/                  # 组件测试
```

### 4.2 GitHub Actions CI/CD

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env: {POSTGRES_USER: test, POSTGRES_PASSWORD: test}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test ./...
      - run: go vet ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: cd web && npm ci && npm run build
```

### 4.3 统一错误码

```go
// internal/pkg/errors/codes.go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

var (
    ErrInvalidPhone    = &AppError{Code: 10001, Message: "手机号格式不正确"}
    ErrCodeExpired     = &AppError{Code: 10002, Message: "验证码已过期"}
    ErrUnauthorized    = &AppError{Code: 20001, Message: "未授权"}
    ErrJobNotFound     = &AppError{Code: 30001, Message: "职位不存在"}
    ErrResumeNotFound  = &AppError{Code: 30002, Message: "简历不存在"}
    ErrFileTooLarge    = &AppError{Code: 40001, Message: "文件大小超过限制"}
    ErrInvalidFileType = &AppError{Code: 40002, Message: "不支持的文件格式"}
)
```

---

## 五、Phase 3：进阶（产品进化）

### 5.1 Redis 缓存

**目标**: 热门职位列表缓存 5 分钟，减少 DB 查询

```go
// internal/job/service.go
func (s *Service) ListWithCache(ctx context.Context, limit, offset int) ([]Job, error) {
    key := fmt.Sprintf("jobs:list:%d:%d", limit, offset)
    // 1. 查 Redis
    // 2. 命中 → 返回
    // 3. 未命中 → 查 DB → 存 Redis (5min TTL)
}
```

### 5.2 Elasticsearch 搜索升级

**时机**: 职位数 > 1万 时替换 PG 全文搜索

```
PG 中文分词 → Elasticsearch IK Analyzer
单关键字搜索 → 多字段加权搜索（title^3, description^1, company^2）
```

### 5.3 WebSocket 实时通知

**目标**: 投递状态变更、面试邀约、聊天 → 实时推送

```
方案: gorilla/websocket
事件: application.updated, interview.created, message.received
```

### 5.4 移动端 PWA

```
web/ → 添加 manifest.json + Service Worker
  安装到手机桌面，离线可用，推送通知
```

---

## 六、排期表

| 时间段 | Phase | 任务 | 交付物 |
|--------|-------|------|--------|
| Day 1-2 | 1 | Docker 一键部署 | `docker compose up` 全栈启动 |
| Day 2-3 | 1 | 种子数据 + API 文档 | 20 职位 10 用户，Swagger 可访问 |
| Day 4-5 | 2 | 单元测试 + 集成测试 | 覆盖率 > 60% |
| Day 5-6 | 2 | CI/CD + 错误码 | PR 自动跑测试 |
| Week 3-4 | 3 | 缓存 + ES + 实时通知 | 按需启动 |

---

## 七、里程碑

```
✅ V1.0  2026-06-22  全栈 MVP 可用（47 commits）
⬜ V2.0  2026-07-01  Docker + 种子数据 + 文档 + CI/CD
⬜ V2.5  2026-07-15  缓存 + 测试 + 错误码完善
⬜ V3.0  2026-08-01  ES 搜索 + 实时通知 + PWA
```

---

## 八、下一步行动

按 Enter 开始 Phase 1 第一项：**Docker Compose 一键部署**。
