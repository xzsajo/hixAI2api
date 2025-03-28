# HIX2API 前端管理界面

这是 HIX2API 的前端管理界面，提供了 API Key 和 Cookie 的管理功能。

## 特性

- Apple 风格的设计
- 响应式布局，支持各种设备
- 支持亮色/暗色模式
- 完整的 API Key 和 Cookie 管理
- 与后端无缝集成

## 开发环境设置

### 前提条件

- Node.js (v14 或更高版本)
- npm (v6 或更高版本)

### 安装步骤

1. 进入 frontend 目录：

```bash
cd frontend
```

2. 安装依赖：

```bash
npm install
```

3. 启动开发服务器：

```bash
npm run dev
```

现在，您可以通过浏览器访问 `http://localhost:5173` 来查看应用程序。

## 构建部署

### 构建项目

执行以下命令构建生产环境版本：

```bash
npm run build
```

或者直接执行构建脚本：

```bash
sh build.sh
```

构建完成后，将在 `dist` 目录中生成静态文件。

### 部署说明

由于项目使用了 Golang 的 embed 功能，构建后的前端文件会被嵌入到后端程序中。无需单独部署前端。

1. 确保 `dist` 目录中的文件可以被 Golang 的 embed 访问
2. 在后端程序中运行正常的构建流程

## 项目结构

```
frontend/
├── public/              # 静态资源
├── src/                 # 源代码
│   ├── assets/          # 资源文件
│   ├── components/      # 共享组件
│   ├── pages/           # 页面组件
│   ├── services/        # API 服务
│   ├── styles/          # 样式文件
│   ├── utils/           # 实用工具函数
│   ├── App.jsx          # 主应用组件
│   ├── App.css          # 应用样式
│   ├── main.jsx         # 入口文件
│   └── index.css        # 全局样式
├── index.html           # HTML 模板
├── package.json         # 项目配置
├── vite.config.js       # Vite 配置
└── README.md            # 本文档
```

## 功能说明

### 身份验证

- 登录页面验证管理密钥
- 安全存储认证信息在 Cookie 中
- 自动处理认证过期和未授权错误

### API Key 管理

- 查看所有 API Key
- 添加新的 API Key
- 编辑现有 API Key
- 删除 API Key

### Cookie 管理

- 查看所有 Cookie 及其额度情况
- 添加新的 Cookie
- 编辑现有 Cookie
- 删除 Cookie
- 刷新所有 Cookie 的额度信息
