/// <reference types="vite/client" />

declare module 'virtual:locale/*' {
  import type { MessageSchema } from './i18n/schema.gen'
  const messages: MessageSchema
  export default messages
}
