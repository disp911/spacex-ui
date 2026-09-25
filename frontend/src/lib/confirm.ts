import { reactive } from 'vue'

export interface ConfirmOptions {
  title: string
  text?: string
  cta: string
  /** Label shown next to the spinner while action runs. */
  busy?: string
  tone?: 'warn' | 'danger'
  /** Runs when confirmed; the dialog stays open, busy, until it settles. */
  action?: () => Promise<unknown> | unknown
}

interface State extends ConfirmOptions {
  open: boolean
  closing: boolean
  running: boolean
  resolve?: (ok: boolean) => void
}

export const confirmState = reactive<State>({ open: false, closing: false, running: false, title: '', cta: '' })

/** Opens the confirmation dialog; resolves true once confirmed (and the action finished). */
export function confirm(opts: ConfirmOptions): Promise<boolean> {
  confirmState.resolve?.(false)
  Object.assign(confirmState, { text: '', busy: '', tone: 'warn', action: undefined }, opts, {
    open: true,
    closing: false,
    running: false,
  })
  return new Promise((resolve) => (confirmState.resolve = resolve))
}

function close(ok: boolean) {
  const resolve = confirmState.resolve
  confirmState.resolve = undefined
  confirmState.closing = true
  setTimeout(() => {
    confirmState.open = false
    confirmState.closing = false
  }, 220)
  resolve?.(ok)
}

export function cancelConfirm() {
  if (!confirmState.running) close(false)
}

export async function runConfirm() {
  if (confirmState.running) return
  if (confirmState.action) {
    confirmState.running = true
    try {
      await confirmState.action()
    } finally {
      confirmState.running = false
    }
  }
  close(true)
}
