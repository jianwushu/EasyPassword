# EasyPassword 前端开发任务清单

本文档根据 `EasyPassword.md` 中的总体设计方案，为前端开发团队规划了详细的任务清单。

## 阶段一：MVP (最小可行产品) - 核心功能实现

目标：构建应用的核心骨架，实现端到端加密的密码管理闭环。

- [x] **1. 项目初始化与基础设置**
    - [x] 使用 Vite 初始化 Vue 3 + TypeScript 项目。
    - [x] 集成 Naive UI (或 Element Plus) 作为 UI 组件库。
    - [x] 配置 ESLint, Prettier, 和 EditorConfig 以确保代码规范。
    - [x] 搭建 `EasyPassword.md` 中建议的项目目录结构 (`/api`, `/components`, `/crypto`, `/router`, `/store`, `/views` 等)。

- [x] **2. 核心加密模块 (`/src/crypto`)**
    - [x] 实现 `deriveKey` 函数，用于从主密码和盐 (salt) 生成加密密钥。优先考虑使用 Web Crypto API (`SubtleCrypto`) 的 `PBKDF2`。
    - [x] 实现 `encryptVaultItem` 函数，使用派生的密钥加密密码项对象 (JSON -> AES-GCM -> Base64)。
    - [x] 实现 `decryptVaultData` 函数，解密从后端获取的数据 (Base64 -> AES-GCM -> JSON)。
    - [x] 编写单元测试，确保加密/解密逻辑的正确性和可逆性。

- [x] **3. 状态管理 (`/src/store`)**
    - [x] 使用 Pinia 创建 `authStore`，用于管理用户认证状态、JWT Token 和 `master_salt`。
    - [x] 创建 `vaultStore`，用于管理加密的密码库数据和解密后的数据。
    - [x] 配置 Pinia 持久化插件，将 `authStore` 中的必要信息 (如 JWT) 安全地存储在 `localStorage` 中。

- [x] **4. API 服务封装 (`/src/api`)**
    - [x] 封装 Axios 实例，统一处理请求头 (Authorization: Bearer {token})、基础 URL 和错误响应。
    - [x] 创建 `auth.js` 模块，包含 `register` 和 `login` 请求函数。
    - [x] 创建 `vault.js` 模块，包含 `getVault`, `addVaultItem`, `updateVaultItem`, `deleteVaultItem` 等函数。

- [x] **5. 路由与页面 (`/src/router`, `/src/views`)**
    - [x] 配置 Vue Router，设置路由守卫 (Navigation Guards) 以保护需要认证的页面。
    - [x] 创建 `Login.vue` 页面：包含登录表单，调用 `authStore` 执行登录逻辑。
    - [x] 创建 `Register.vue` 页面：包含注册表单，处理用户注册流程。
    - [x] 创建 `Dashboard.vue` 主控台页面，作为登录后的主界面，并包含创建、读取、更新、删除密码项的UI。

- [x] **6. 核心业务流程实现**
    - [x] **注册流程**: 用户输入用户名和主密码 -> 前端调用 `authStore` 的注册 action -> 发送请求到后端 -> 完成注册。
    - [x] **登录流程**: 用户输入用户名和主密码 -> 前端调用 `authStore` 的登录 action -> 成功后，保存 JWT 和 `master_salt` 到 Pinia Store -> 跳转到主控台。
    - [x] **密码库 CRUD**:
        - [x] **创建**: 在前端表单中输入明文 -> 使用主密码和 `master_salt` 加密 -> 发送密文到后端保存。
        - [x] **读取**: 从后端获取所有加密项 -> 用户输入主密码（如果尚未输入） -> 在前端解密并展示。
        - [x] **更新**: 类似创建流程，发送更新后的密文。
        - [x] **删除**: 调用后端 API 删除指定项。

## 阶段二：功能增强与体验优化

目标：在核心功能基础上，增加实用功能，提升用户体验。

- [x] **1. 密码生成器 (`/src/components/PasswordGenerator.vue`)**
    - [x] 创建一个可复用组件，用于生成高强度随机密码。
    - [x] 提供可配置选项，如密码长度、是否包含数字、特殊字符等。

- [x] **2. 搜索功能**
    - [x] 在主控台页面添加搜索框。
    - [x] 实现客户端搜索逻辑：对 `vaultStore` 中已解密的密码项列表进行实时过滤。

- [x] **3. UI/UX 优化**
    - [x] 实现密码项的分类或标签功能。
    - [x] 优化表单验证和用户反馈（如加载状态、成功/错误提示）。
    - [x] 提升应用的响应式布局，适配不同屏幕尺寸。

## 阶段三：安全加固与部署

目标：审查并加固应用安全，为正式上线做准备。

- [ ] **1. 安全审计与加固**
    - [ ] 审查代码，防止 XSS 攻击（例如，确保不使用 `v-html` 渲染用户输入内容）。
    - [ ] 配合后端，确保 CSRF 防护措施到位。
    - [ ] 实现剪贴板自动清除功能，以增强安全性。

- [ ] **2. 构建与部署**
    - [x] 编写 `Dockerfile` 用于构建前端生产环境的 Docker 镜像。
    - [x] 在 `docker-compose.yml` 中集成前端服务，方便本地一键启动。
    - [ ] 配置 CI/CD (如 GitHub Actions) 流程，实现自动化测试、构建和部署。