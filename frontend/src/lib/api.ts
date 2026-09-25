import { panelUrl } from './config'
import { toast } from './toast'

// The panel answers with { success, msg, obj }. Requests go out as
// form-urlencoded POSTs, which the Go controllers bind with PostForm.

export interface Msg<T = unknown> {
  success: boolean
  msg: string
  obj: T
}

type Body = Record<string, string | number | boolean | undefined | null> | FormData

interface Options {
  /** Show msg as a toast (the old panel does this for every response). */
  notify?: boolean
}

function encode(body?: Body): BodyInit | undefined {
  if (!body) return undefined
  if (body instanceof FormData) return body
  const p = new URLSearchParams()
  for (const [k, v] of Object.entries(body)) {
    if (v !== undefined && v !== null) p.append(k, String(v))
  }
  return p
}

async function request<T>(method: 'GET' | 'POST', path: string, body?: Body, opts: Options = {}): Promise<Msg<T>> {
  const notify = opts.notify ?? true
  let msg: Msg<T>
  try {
    const resp = await fetch(panelUrl(path), {
      method,
      body: method === 'POST' ? encode(body) : undefined,
      headers: { 'X-Requested-With': 'XMLHttpRequest' },
      credentials: 'same-origin',
    })
    if (resp.status === 401) {
      location.reload()
      return new Promise(() => {})
    }
    if (!resp.ok) throw new Error(`${resp.status} ${resp.statusText}`)
    const data = await resp.json()
    msg = data && typeof data === 'object' && 'success' in data
      ? { success: !!data.success, msg: data.msg ?? '', obj: data.obj as T }
      : { success: false, msg: 'Unexpected response', obj: null as T }
  } catch (e) {
    msg = { success: false, msg: e instanceof Error ? e.message : String(e), obj: null as T }
  }
  if (notify && msg.msg) toast(msg.msg, msg.success ? 'success' : 'error')
  return msg
}

export const api = {
  get: <T>(path: string, opts?: Options) => request<T>('GET', path, undefined, opts),
  post: <T>(path: string, body?: Body, opts?: Options) => request<T>('POST', path, body, opts),
}

/** POSTs a form and saves the response body as a file. */
export async function download(path: string, body: Body, fallbackName: string) {
  const resp = await fetch(panelUrl(path), {
    method: 'POST',
    body: encode(body),
    headers: { 'X-Requested-With': 'XMLHttpRequest' },
    credentials: 'same-origin',
  })
  if (!resp.ok) throw new Error(`${resp.status} ${resp.statusText}`)
  const name = /filename="([^"]+)"/.exec(resp.headers.get('content-disposition') ?? '')?.[1] ?? fallbackName
  saveBlob(await resp.blob(), name)
}

export function saveBlob(blob: Blob, name: string) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = name
  document.body.appendChild(a)
  a.click()
  a.remove()
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))
