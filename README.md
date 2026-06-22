# 招聘平台 (Recruitment Platform)

类似 BOSS 直聘的在线招聘平台，支持求职者与 HR 实时直聊。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go + Gin |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 前端 | React + TypeScript + Ant Design 5 |
| 即时通讯 | 腾讯云 IM SDK |

## 功能

### 求职端
- 手机验证码登录
- 简历管理（多份简历，附件上传）
- 职位搜索（全文检索 + 城市过滤）
- 投递简历
- 接收面试邀约（接受/婉拒/导航）

### 招聘端
- HR 企业认证
- 职位发布/编辑/下架
- 候选人管理（查看简历/标记状态）
- 发起面试邀约（时间/地址/导航）
- 智能推荐（技能匹配）

## 快速开始

### 前置条件
- Go 1.22+
- Node.js 18+
- PostgreSQL 16
- Redis 7

### 1. 配置

```bash
cp config/config.example.yaml config/config.yaml
# 编辑 config/config.yaml，填入你的数据库和 Redis 连接信息
```

### 2. 数据库迁移

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "postgres://user:pass@localhost:5432/recruitment?sslmode=disable" up
```

### 3. 启动后端

```bash
go build -o server.exe ./cmd/server/
./server.exe
# API 运行在 http://localhost:8080
```

### 4. 启动前端

```bash
cd web
npm install
npx vite
# 前端运行在 http://localhost:5173
```

### 5. 开发登录

项目使用 Mock SMS，验证码会打印在后端终端日志中：
```
[MOCK SMS] To: 13800138000, Code: 123456
```

## 项目结构

```
├── cmd/server/          # 入口
├── internal/
│   ├── auth/            # 认证模块
│   ├── user/            # 用户模块
│   ├── resume/          # 简历模块
│   ├── job/             # 职位模块
│   ├── application/     # 投递模块
│   ├── interview/       # 面试模块
│   ├── chat/            # IM 桥接
│   ├── recommend/       # 推荐系统
│   ├── middleware/       # HTTP 中间件
│   ├── config/          # 配置管理
│   └── pkg/             # 基础设施包
│       ├── auth/        # JWT
│       ├── broker/      # 消息代理
│       ├── cos/         # 对象存储
│       ├── maps/        # 地图服务
│       ├── redis/       # Redis 客户端
│       └── sms/         # 短信客户端
├── migrations/          # 数据库迁移
├── config/              # 配置文件
└── web/                 # React 前端
    └── src/
        ├── api/         # API 客户端
        ├── pages/       # 页面组件
        ├── layouts/     # 布局组件
        └── stores/      # 状态管理
```

## License

MIT
