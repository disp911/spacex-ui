// Thin client for the panel's JSON API. Every endpoint answers
// { success, msg, obj }; a failed call throws ApiError with the server message.

import { serverUrl } from './config'

export interface ApiReply<T> {
  success: boolean
  msg: string
  obj: T
}

export class ApiError extends Error {
  constructor(message: string, readonly status: number) {
    super(message)
  }
}

let onSessionLost = () => {}

/** Called when the server no longer knows the session; the app shows the login screen. */
export function handleSessionLost(handler: () => void) {
  onSessionLost = handler
}

async function request<T>(method: 'GET' | 'POST', path: string, body?: unknown): Promise<T> {
  const init: RequestInit = {
    method,
    credentials: 'same-origin',
    // The panel answers 401 instead of a redirect to such requests.
    headers: { 'X-Requested-With': 'XMLHttpRequest', Accept: 'application/json' },
  }
  if (body instanceof FormData) {
    init.body = body
  } else if (body !== undefined) {
    ;(init.headers as Record<string, string>)['Content-Type'] = 'application/json'
    init.body = JSON.stringify(body)
  }

  const res = await fetch(serverUrl(path), init)
  // The API hides itself behind 404 from requests without a session.
  if (res.status === 401 || (res.status === 404 && path.startsWith('panel/api/'))) {
    onSessionLost()
    throw new ApiError(res.statusText || 'Unauthorized', 401)
  }
  if (!res.ok) throw new ApiError(res.statusText || `HTTP ${res.status}`, res.status)

  const reply = (await res.json()) as ApiReply<T>
  if (!reply.success) throw new ApiError(reply.msg, res.status)
  return reply.obj
}

export const api = {
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body?: unknown) => request<T>('POST', path, body),
}
