# EasyPassword - 安全、开源的密码管理器

EasyPassword 是一个基于零知识架构构建的现代化密码管理器。它旨在提供一个安全、可靠且用户友好的方式来存储和管理您的所有密码。由于采用了端到端加密，您的主密码和存储的任何密码数据都不会以明文形式传输或存储在服务器上，确保只有您自己才能访问您的敏感信息。

## ✨ 核心特性

*   **零知识架构**: 服务器无法访问您的明文密码数据。所有加密和解密都在您的设备上（客户端）完成。
*   **端到端加密 (E2EE)**: 您的密码库使用从您的主密码派生的密钥进行加密，确保数据在传输和静态存储时的绝对安全。
*   **前后端分离**: 采用现代化的前后端分离架构，易于开发、扩展和维护。
*   **Docker 化部署**: 整个应用（前端、后端、数据库）可以通过 Docker Compose 一键启动，极大简化了部署和本地开发流程。
*   **开源**: 完全开源，欢迎社区贡献和审计。

## 🛠️ 技术栈

本项目采用了一系列现代化、高性能的技术构建：

| 类别 | 技术 |
| :--- | :--- |
| **前端** | Vue 3, Vite, TypeScript, Pinia, Naive UI, Axios, SubtleCrypto |
| **后端** | Golang, Gin, GORM, JWT |
| **数据库** | PostgreSQL, BoltDB |
| **容器化** | Docker, Docker Compose |

## 🚀 如何运行

我们强烈推荐使用 Docker 来运行此项目，因为它可以保证环境的一致性并简化设置。

### 先决条件

*   [Docker](https://www.docker.com/get-started)
*   [Docker Compose](https://docs.docker.com/compose/install/)

### 一键启动

1.  克隆本仓库到本地：
    ```bash
    git clone https://github.com/your-username/EasyPassword.git
    cd EasyPassword
    ```

2.  **重要**: 在生产环境部署前，请务必修改 `docker-compose.yml` 文件中的 `JWT_SECRET_KEY`，将其替换为一个强大且唯一的密钥。
    ```yaml
    # docker-compose.yml
    environment:
      - JWT_SECRET_KEY=your-super-secret-key-change-me # <--- 修改这里
    ```

3.  使用 Docker Compose 构建并启动所有服务：
    ```bash
    docker-compose up --build -d
    ```
    该命令会在后台构建并启动前端、后端和数据库服务。

4.  服务访问：
    *   **前端应用**: 打开浏览器访问 `http://localhost:80`
    *   **后端 API**: 服务运行在 `http://localhost:8080`

### 停止服务

要停止所有正在运行的服务，请执行：
```bash
docker-compose down
```

## 🏗️ 项目结构

本项目是一个单体仓库 (monorepo)，包含两个主要的子项目：

*   `password-manager-frontend/`: 包含所有前端代码。这是一个基于 Vue 3 的单页面应用 (SPA)，负责用户界面、交互以及客户端的加解密逻辑。
*   `password-manager-backend/`: 包含所有后端代码。这是一个基于 Golang Gin 框架的 API 服务，负责用户认证、处理和存储加密后的数据。

## 🤝 贡献

欢迎对 EasyPassword 项目做出贡献！我们鼓励任何形式的贡献，包括但不限于：

*   报告 Bug
*   提交功能需求
*   编写和改进文档
*   提交拉取请求 (Pull Requests)
