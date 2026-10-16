// WireGuard keys made in the browser: X25519 in plain arithmetic, since the
// panel is often opened over HTTP, where crypto.subtle is not available.
// Ported from the old panel's Wireguard class (TweetNaCl's field code).

type Field = Float64Array

const gf = (init?: number[]): Field => {
  const r = new Float64Array(16)
  if (init) init.forEach((v, i) => (r[i] = v))
  return r
}

function carry(o: Field) {
  for (let i = 0; i < 16; i++) {
    o[(i + 1) % 16]! += (i < 15 ? 1 : 38) * Math.floor(o[i]! / 65536)
    o[i]! &= 0xffff
  }
}

function cswap(p: Field, q: Field, b: number) {
  const c = ~(b - 1)
  for (let i = 0; i < 16; i++) {
    const t = c & (p[i]! ^ q[i]!)
    p[i]! ^= t
    q[i]! ^= t
  }
}

function pack(o: Uint8Array, n: Field) {
  const m = gf()
  const t = gf()
  for (let i = 0; i < 16; i++) t[i] = n[i]!
  carry(t)
  carry(t)
  carry(t)
  for (let j = 0; j < 2; j++) {
    m[0] = t[0]! - 0xffed
    for (let i = 1; i < 15; i++) {
      m[i] = t[i]! - 0xffff - ((m[i - 1]! >> 16) & 1)
      m[i - 1]! &= 0xffff
    }
    m[15] = t[15]! - 0x7fff - ((m[14]! >> 16) & 1)
    const b = (m[15]! >> 16) & 1
    m[14]! &= 0xffff
    cswap(t, m, 1 - b)
  }
  for (let i = 0; i < 16; i++) {
    o[2 * i] = t[i]! & 0xff
    o[2 * i + 1] = t[i]! >> 8
  }
}

function add(o: Field, a: Field, b: Field) {
  for (let i = 0; i < 16; i++) o[i] = (a[i]! + b[i]!) | 0
}

function subtract(o: Field, a: Field, b: Field) {
  for (let i = 0; i < 16; i++) o[i] = (a[i]! - b[i]!) | 0
}

function multmod(o: Field, a: Field, b: Field) {
  const t = new Float64Array(31)
  for (let i = 0; i < 16; i++) for (let j = 0; j < 16; j++) t[i + j]! += a[i]! * b[j]!
  for (let i = 0; i < 15; i++) t[i]! += 38 * t[i + 16]!
  for (let i = 0; i < 16; i++) o[i] = t[i]!
  carry(o)
  carry(o)
}

function invert(o: Field, i: Field) {
  const c = gf()
  for (let a = 0; a < 16; a++) c[a] = i[a]!
  for (let a = 253; a >= 0; a--) {
    multmod(c, c, c)
    if (a !== 2 && a !== 4) multmod(c, c, i)
  }
  for (let a = 0; a < 16; a++) o[a] = c[a]!
}

function clamp(z: Uint8Array) {
  z[31] = (z[31]! & 127) | 64
  z[0]! &= 248
}

function publicKey(privateKey: Uint8Array): Uint8Array {
  const z = new Uint8Array(privateKey)
  const a = gf([1])
  const b = gf([9])
  const c = gf()
  const d = gf([1])
  const e = gf()
  const f = gf()
  const k121665 = gf([0xdb41, 1])
  const nine = gf([9])
  clamp(z)
  for (let i = 254; i >= 0; i--) {
    const r = (z[i >>> 3]! >>> (i & 7)) & 1
    cswap(a, b, r)
    cswap(c, d, r)
    add(e, a, c)
    subtract(a, a, c)
    add(c, b, d)
    subtract(b, b, d)
    multmod(d, e, e)
    multmod(f, a, a)
    multmod(a, c, a)
    multmod(c, b, e)
    add(e, a, c)
    subtract(a, a, c)
    multmod(b, a, a)
    subtract(c, d, f)
    multmod(a, c, k121665)
    add(a, a, d)
    multmod(c, c, a)
    multmod(a, d, f)
    multmod(d, b, nine)
    multmod(b, e, e)
    cswap(a, b, r)
    cswap(c, d, r)
  }
  invert(c, c)
  multmod(a, a, c)
  const out = new Uint8Array(32)
  pack(out, a)
  return out
}

const toBase64 = (bytes: Uint8Array) => btoa(String.fromCharCode(...bytes))
const fromBase64 = (text: string) => Uint8Array.from(atob(text), (ch) => ch.charCodeAt(0))

/** A new key pair, or the public key of the private one given; both in base64. */
export function wireguardKeys(privateKey = ''): { privateKey: string; publicKey: string } {
  let secret: Uint8Array
  if (privateKey) secret = fromBase64(privateKey)
  else {
    secret = crypto.getRandomValues(new Uint8Array(32))
    clamp(secret)
  }
  return { privateKey: privateKey || toBase64(secret), publicKey: toBase64(publicKey(secret)) }
}

/** The bytes of a base64 string, e.g. WARP's client id as the reserved field. */
export const base64Bytes = (text: string) => [...fromBase64(text)]
