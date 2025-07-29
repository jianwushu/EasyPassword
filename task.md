本方案采用 Vite 作为构建工具，Vue 3 作为核心UI框架，并结合 TypeScript 提供类型安全。我们将构建一个包含 Popup 弹窗、Content Script (内容脚本) 和 Background Script (后台脚本) 的完整插件架构。

核心优势
极致的开发体验：Vite 提供闪电般的冷启动速度和亚秒级的热模块替换 (HMR)，让您在修改 Popup 或 Content Script 的 UI 时能立即看到效果。
现代化的框架：Vue 3 的 Composition API 使得逻辑复用和代码组织变得极其灵活和清晰，非常适合构建复杂的插件功能。
类型安全：TypeScript 的全面集成，从 Vue 组件到 Chrome API 调用，都能获得强大的类型提示和编译时错误检查，显著减少运行时 Bug。
CSS 隔离：通过为 Content Script 实施 Shadow DOM 策略，彻底解决了 Vue 组件样式与宿主页面样式冲突的业界难题。
高效构建：Vite 基于 Rollup 的生产环境构建，默认支持 Tree-shaking，确保最终打包的插件体积最小化。

基于现有的password-manager-frontend项目，构建一个password-manager-extension浏览器插件项目,样式参考现有项目

新增开发功能如下:
1.  实现主账号登录,退出功能
2.  登录成功后可以监听其他网页信息的登录表单提交请求，触发密码项新增弹窗（自动填写相应数据，账号、密码、站点名称、网址等信息）