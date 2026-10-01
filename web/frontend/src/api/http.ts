import { panelUrl } from '@/env'

/** Envelope every panel endpoint answers with (web/entity.Msg). */
export interface Msg<T = unknown> {
  success: boolean
  msg: string
  obj: T
}

export class ApiError extends Error {
  constructor(message: string, public status = 0) {
    super(message)
  }
}

type Params = Record<string, string | number | boolean | undefined | null>

function encode(data?: Params): URLSearchParams {
  const body = new URLSearchParams()
  if (data) {
    for (const [k, v] of Object.entries(data)) {
      if (v !== undefined && v !== null) body.append(k, String(v))
    }
  }
  return body
}

let onUnauthorized: (() => void) | null = null
/** Called when the session expired (HTTP 401). */
export function setUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn
}

async function request<T>(method: 'GET' | 'POST', path: string, body?: BodyInit, query?: Params): Promise<Msg<T>> {
  let url = panelUrl(path)
  if (query) {
    const qs = encode(query).toString()
    if (qs) url += (url.includes('?') ? '&' : '?') + qs
  }
  const headers: Record<string, string> = { 'X-Requested-With': 'XMLHttpRequest' }
  if (body instanceof URLSearchParams) {
    headers['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8'
  }

  let resp: Response
  try {
    resp = await fetch(url, { method, body, headers, credentials: 'same-origin' })
  } catch (e) {
    throw new ApiError((e as Error).message || 'Network error')
  }

  if (resp.status === 401) {
    onUnauthorized?.()
    throw new ApiError('Unauthorized', 401)
  }

  let data: unknown = null
  const text = await resp.text()
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      throw new ApiError(`Unexpected response (${resp.status})`, resp.status)
    }
  }
  if (!resp.ok && !(data && typeof data === 'object' && 'success' in data)) {
    throw new ApiError(`HTTP ${resp.status}`, resp.status)
  }
  if (data && typeof data === 'object' && 'success' in data) {
    return data as Msg<T>
  }
  return { success: true, msg: '', obj: data as T }
}

export const http = {
  get<T>(path: string, query?: Params) {
    return request<T>('GET', path, undefined, query)
  },
  post<T>(path: string, data?: Params) {
    return request<T>('POST', path, encode(data))
  },
  upload<T>(path: string, form: FormData) {
    return request<T>('POST', path, form)
  },
  /** POST that returns a file; resolves to the blob and its suggested filename. */
  async download(path: string, data?: Params): Promise<{ blob: Blob; filename: string }> {
    const resp = await fetch(panelUrl(path), {
      method: 'POST',
      body: encode(data),
      credentials: 'same-origin',
      headers: {
        'X-Requested-With': 'XMLHttpRequest',
        'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
      },
    })
    if (resp.status === 401) {
      onUnauthorized?.()
      throw new ApiError('Unauthorized', 401)
    }
    if (!resp.ok) throw new ApiError(`HTTP ${resp.status}`, resp.status)
    const cd = resp.headers.get('content-disposition') || ''
    const filename = /filename="([^"]+)"/.exec(cd)?.[1] || 'download'
    return { blob: await resp.blob(), filename }
  },
}

/** Unwrap a Msg or throw its message. */
export async function unwrap<T>(p: Promise<Msg<T>>): Promise<T> {
  const m = await p
  if (!m.success) throw new ApiError(m.msg || 'Request failed')
  return m.obj
}

export function saveBlob(blob: Blob, filename: string) {
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  setTimeout(() => URL.revokeObjectURL(a.href), 1000)
}
