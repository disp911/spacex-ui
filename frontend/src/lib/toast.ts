import { reactive } from 'vue'

export type ToastKind = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  text: string
  kind: ToastKind
}

export const toasts = reactive<Toast[]>([])
let seq = 0

export function toast(text: string, kind: ToastKind = 'info', ms = 3200) {
  const id = ++seq
  toasts.push({ id, text, kind })
  setTimeout(() => {
    const i = toasts.findIndex((x) => x.id === id)
    if (i >= 0) toasts.splice(i, 1)
  }, ms)
}
