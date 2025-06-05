# CloudDrop

云滴是一款webshell管理工具，使用前后端分离架构开发。

## 技术栈

### 后端
- **Go 1.23+** - 主要开发语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **SQLite/MySQL** - 数据库
- **JWT** - 身份认证
- **WebSocket** - 实时通信

### 前端
- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Ant Design** - UI组件库
- **Axios** - HTTP客户端
- **React Router** - 路由管理

## 功能特性

- 🔐 用户认证与权限管理
- 🌐 多webshell连接管理
- 📁 远程文件管理
- 💻 命令执行终端
- 📊 连接状态监控
- 🔄 实时日志查看
- 📦 插件扩展支持

## 项目结构

```
CloudDrop/
├── backend/          # Go后端服务
│   ├── cmd/         # 启动入口
│   ├── internal/    # 内部逻辑
│   ├── pkg/         # 公共包
│   └── configs/     # 配置文件
├── frontend/        # React前端
│   ├── src/         # 源代码
│   ├── public/      # 静态资源
│   └── build/       # 构建输出
└── docs/           # 文档
```

## 快速开始

### 后端启动
```bash
cd backend
go mod tidy
go run cmd/main.go
```

### 前端启动
```bash
cd frontend
npm install
npm start
```

## 开发计划

- [x] 项目初始化
- [ ] 后端API开发
- [ ] 前端界面开发
- [ ] webshell连接功能
- [ ] 文件管理功能
- [ ] 命令执行功能
- [ ] 部署与优化

