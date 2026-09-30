import { env } from '@/env'

type Handler = (payload: any) => void

/**
 * Minimal client for the panel's /ws hub. Messages are
 * `{ type, payload, time }`; reconnects with backoff and reports
 * `down` after too many failures so callers can fall back to polling.
 */
export class PanelSocket {
  private ws: WebSocket | null = null
  private handlers = new Map<string, Set<Handler>>()
  private attempts = 0
  private timer: number | undefined
  private closed = false

  static readonly maxAttempts = 6

  connect() {
    this.closed = false
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
    let socket: WebSocket
    try {
      socket = new WebSocket(`${proto}//${location.host}${env.basePath}ws`)
    } catch {
      this.scheduleReconnect()
      return
    }
    this.ws = socket
    socket.onopen = () => {
      this.attempts = 0
      this.emit('open', null)
    }
    socket.onmessage = (ev) => {
      if (typeof ev.data !== 'string' || ev.data.length > 4_000_000) return
      try {
        const m = JSON.parse(ev.data)
        if (m && typeof m.type === 'string') this.emit(m.type, m.payload)
      } catch {
        /* ignore malformed frames */
      }
    }
    socket.onclose = () => {
      this.ws = null
      if (!this.closed) this.scheduleReconnect()
    }
  }

  private scheduleReconnect() {
    this.attempts++
    if (this.attempts > PanelSocket.maxAttempts) {
      this.emit('down', null)
      return
    }
    const delay = Math.min(30_000, 1000 * 2 ** (this.attempts - 1))
    this.timer = window.setTimeout(() => this.connect(), delay)
  }

  on(type: string, fn: Handler) {
    if (!this.handlers.has(type)) this.handlers.set(type, new Set())
    this.handlers.get(type)!.add(fn)
    return () => this.handlers.get(type)?.delete(fn)
  }

  private emit(type: string, payload: unknown) {
    this.handlers.get(type)?.forEach((fn) => fn(payload))
  }

  get connected() {
    return this.ws?.readyState === WebSocket.OPEN
  }

  close() {
    this.closed = true
    clearTimeout(this.timer)
    this.ws?.close()
    this.ws = null
  }
}
