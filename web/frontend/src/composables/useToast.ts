import { reactive } from 'vue'

export type ToastTone = 'success' | 'danger' | 'info' | 'warning'

export interface Toast {
  id: number
  tone: ToastTone
  text: string
}

let seq = 0
export const toasts = reactive<Toast[]>([])

export function toast(text: string, tone: ToastTone = 'info', ms = 3600) {
  if (!text) return
  const id = ++seq
  toasts.push({ id, tone, text })
  setTimeout(() => dismissToast(id), ms)
}

export function dismissToast(id: number) {
  const i = toasts.findIndex((t) => t.id === id)
  if (i >= 0) toasts.splice(i, 1)
}

/** Show the outcome of a panel Msg: its text as success or error. */
export function toastMsg(m: { success: boolean; msg: string }) {
  if (m.msg) toast(m.msg, m.success ? 'success' : 'danger')
}

export function toastError(e: unknown) {
  toast((e as Error)?.message || String(e), 'danger')
}
