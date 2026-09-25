import { config } from './config'

// Real-time updates from the panel's /ws hub: { type, payload, time }.
// Reconnects with exponential backoff, then once a minute.

type Handler = (payload: any) => void

const handlers = new Map<string, Set<Handler>>()
let socket: WebSocket | null = null
let attempts = 0
let timer: number | undefined
let connected = false

function emit(type: string, payload?: unknown) {
  handlers.get(type)?.forEach((h) => {
    try {
      h(payload)
    } catch (e) {
      console.error(`ws handler "${type}" failed`, e)
    }
  })
}

function open() {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  try {
    socket = new WebSocket(`${proto}//${location.host}${config.basePath}ws`)
  } catch {
    schedule()
    return
  }
  socket.onopen = () => {
    attempts = 0
    connected = true
    emit('connected')
  }
  socket.onmessage = (e) => {
    try {
      const m = JSON.parse(e.data)
      if (m && typeof m.type === 'string') emit(m.type, m.payload)
    } catch {
      /* ignore malformed frames */
    }
  }
  socket.onclose = () => {
    socket = null
    if (connected) emit('disconnected')
    connected = false
    schedule()
  }
}

function schedule() {
  clearTimeout(timer)
  attempts++
  const base = attempts <= 10 ? Math.min(30_000, 1000 * 2 ** (attempts - 1)) : 60_000
  timer = window.setTimeout(open, base * (0.75 + Math.random() * 0.5))
}

export const ws = {
  connect() {
    if (!socket) open()
  },
  get connected() {
    return connected
  },
  on(type: string, h: Handler) {
    if (!handlers.has(type)) handlers.set(type, new Set())
    handlers.get(type)!.add(h)
    return () => handlers.get(type)?.delete(h)
  },
}
