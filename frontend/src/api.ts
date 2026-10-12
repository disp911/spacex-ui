// Thin client for the panel's JSON API. Every endpoint answers
// { success, msg, obj }; a failed call throws ApiError with the server message.

import { config, panelUrl } from './config'

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

/** Sends the browser to the login page; the session is gone. */
export function toLogin(): never {
  window.location.assign(config.base)
  throw new ApiError('Not logged in', 401)
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

  const res = await fetch(panelUrl(path), init)
  // The API hides itself behind 404 from requests without a session.
  if (res.status === 401 || (res.status === 404 && path.startsWith('panel/api/'))) toLogin()
  if (!res.ok) throw new ApiError(res.statusText || `HTTP ${res.status}`, res.status)

  const reply = (await res.json()) as ApiReply<T>
  if (!reply.success) throw new ApiError(reply.msg, res.status)
  return reply.obj
}

export const api = {
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body?: unknown) => request<T>('POST', path, body),
}
