import { reactive } from 'vue'

export interface Toast {
  id: number
  /** ok done, warn nothing happened, err failed. */
  kind: 'ok' | 'warn' | 'err'
  text: string
}

/** Messages shown by XToaster in the corner of the screen. */
export const toasts = reactive<Toast[]>([])
let next = 0

export function toast(kind: Toast['kind'], text: string, ms = kind === 'ok' ? 3000 : 6000) {
  const id = ++next
  toasts.push({ id, kind, text })
  setTimeout(() => dismissToast(id), ms)
}

export function dismissToast(id: number) {
  const i = toasts.findIndex((t) => t.id === id)
  if (i >= 0) toasts.splice(i, 1)
}
