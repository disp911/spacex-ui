// The panel's WebSocket: the server pushes the status every 2 seconds and
// tells about changes. One connection serves every subscriber and reconnects
// after a drop.

import { serverUrl } from './config'

type Handler = (payload: unknown) => void

const handlers = new Map<string, Set<Handler>>()
let socket: WebSocket | null = null
let retry = 0
let timer = 0
let dropped = false

function connect() {
  const url = new URL(serverUrl('ws'), location.href)
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
  socket = new WebSocket(url)
  socket.onopen = () => {
    retry = 0
    timer = 0
    // Messages sent while the connection was down are lost.
    if (dropped) handlers.get('reconnected')?.forEach((h) => h(null))
    dropped = false
  }
  socket.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data as string) as { type: string; payload: unknown }
      handlers.get(msg.type)?.forEach((h) => h(msg.payload))
    } catch {
      // not a message for us
    }
  }
  socket.onclose = () => {
    socket = null
    if (!handlers.size) return
    // 1 s, 2 s, 4 s … up to 30 s between attempts.
    dropped = true
    timer = window.setTimeout(connect, Math.min(30_000, 1000 * 2 ** retry++))
  }
}

/** Calls handler for every message of the type; returns the unsubscribe function. */
export function subscribe(type: string, handler: Handler) {
  if (!handlers.has(type)) handlers.set(type, new Set())
  handlers.get(type)!.add(handler)
  if (!socket && !timer) connect()
  return () => {
    handlers.get(type)?.delete(handler)
    if (!handlers.get(type)?.size) handlers.delete(type)
    if (!handlers.size) {
      clearTimeout(timer)
      timer = 0
      socket?.close()
      socket = null
    }
  }
}
