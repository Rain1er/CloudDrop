# CloudDrop

CloudDrop（云滴）是一款前后端分离的网站管理与 WebShell 交互工作台。当前版本聚焦连接管理、存活验证、系统信息、命令执行、代码执行、文件管理与数据库管理。

> 本项目用于授权环境、内网靶场、CTF、沙箱与自有资产测试。

## 功能特性

- WebShell 连接管理：新增、编辑、删除、刷新、批量存活验证。
- 多类型支持：`php`、`java`、`c#`、`asp`。
- 系统信息：进入目标后自动获取基础环境信息。
- 命令执行：基于 `xterm.js` 的 Web 终端界面。
- 代码执行：独立 payload 输入与结果输出。
- 文件管理：路径输入、树状目录、文件查看、语法高亮、压缩、解压。
- 数据库管理：连接配置、库列表、SQL 查询与结构化结果展示。
- 动态密钥通信：入口脚本根据 `timezone` seed 派生密钥，并通过 session 向动态 payload 传递密钥上下文。

## 技术栈

### 后端

- Go 1.23+
- Gin
- GORM
- SQLite
- JWT

### 前端

- React 18
- TypeScript
- Vite
- Axios
- xterm.js
- react-arborist
- react-syntax-highlighter

## 项目结构

```text
CloudDrop/
├── backend/
│   ├── cmd/                  # 后端启动入口
│   ├── config/               # 配置与数据库初始化
│   ├── internal/
│   │   ├── handler/          # HTTP handler
│   │   ├── model/            # 数据模型
│   │   └── service/          # WebShell 交互服务
│   ├── pkg/
│   │   ├── api/              # PHP / Java / ASP / .NET payload
│   │   ├── middleware/       # 中间件
│   │   └── util/             # 加密、HTTP 工具
│   └── routes/               # 路由注册
├── frontend/
│   ├── src/                  # React 前端源码
│   ├── dist/                 # 构建输出
│   └── package.json
├── shell/                    # 目标端入口脚本 api.php / api.jsp
└── docs/                     # 开发文档
```

## 本地开发

### 启动后端

```bash
cd backend
go mod tidy
go run ./cmd/main.go
```

### 启动前端

```bash
cd frontend
npm install
npm run dev
```

## 开发状态

- [x] 后端基础 API
- [x] WebShell 连接管理
- [x] PHP / Java 目标端入口测试
- [x] React 前端工作台
- [x] 系统信息、命令执行、代码执行
- [x] 文件管理与压缩解压
- [x] 数据库管理界面
- [ ] 后端生产化部署配置
- [ ] 外部数据库适配
- [ ] 协议 stateless 化评估
