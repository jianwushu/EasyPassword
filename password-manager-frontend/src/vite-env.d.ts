/// <reference types="vite/client" />
/// <reference types="pinia" />

declare module '*.vue' {
  import { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
