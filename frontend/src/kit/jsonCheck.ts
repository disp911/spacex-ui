// Finds the first mistake in JSON text and says what was expected there.
// Browsers word JSON.parse errors each in their own way, often without the
// place; this walks the text itself, so the line and the words are the same
// everywhere.

export type JsonIssueKind = 'commaOrBrace' | 'commaOrBracket' | 'key' | 'colon' | 'value' | 'end' | 'trailing' | 'string' | 'number'

export interface JsonIssue {
  kind: JsonIssueKind
  /** Offset in the text and its line, from 1. */
  pos: number
  line: number
}

class Stop {
  constructor(readonly kind: JsonIssueKind, readonly pos: number) {}
}

const NUMBER = /-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?/y
const ESCAPE = /["\\/bfnrt]|u[0-9a-fA-F]{4}/y

export function checkJson(text: string): JsonIssue | null {
  let i = 0
  // Where the last whole value ended: a missing comma belongs there, not on the next line.
  let lastEnd = 0

  const fail = (kind: JsonIssueKind, at = i): never => {
    throw new Stop(at >= text.length && kind !== 'trailing' ? 'end' : kind, Math.min(at, text.length))
  }
  const ws = () => {
    while (i < text.length && ' \t\n\r'.includes(text[i]!)) i++
  }

  function string() {
    i++
    for (;;) {
      if (i >= text.length) fail('string')
      const c = text[i]!
      if (c === '"') {
        i++
        return
      }
      if (c < ' ') fail('string')
      if (c === '\\') {
        ESCAPE.lastIndex = i + 1
        if (!ESCAPE.test(text)) fail('string', i)
        i = ESCAPE.lastIndex
      } else i++
    }
  }

  function value() {
    ws()
    const c = text[i]
    if (c === '{') object()
    else if (c === '[') array()
    else if (c === '"') string()
    else if (c === '-' || (c !== undefined && c >= '0' && c <= '9')) {
      NUMBER.lastIndex = i
      if (!NUMBER.test(text)) fail('number')
      i = NUMBER.lastIndex
    } else {
      const word = ['true', 'false', 'null'].find((w) => text.startsWith(w, i))
      if (!word) fail('value')
      i += word!.length
    }
    lastEnd = i
  }

  function object() {
    i++
    ws()
    if (text[i] === '}') {
      i++
      return
    }
    for (;;) {
      ws()
      if (text[i] !== '"') fail('key')
      string()
      ws()
      if (text[i] !== ':') fail('colon')
      i++
      value()
      ws()
      if (text[i] === ',') i++
      else if (text[i] === '}') {
        i++
        return
      } else fail('commaOrBrace', i >= text.length ? i : lastEnd)
    }
  }

  function array() {
    i++
    ws()
    if (text[i] === ']') {
      i++
      return
    }
    for (;;) {
      value()
      ws()
      if (text[i] === ',') i++
      else if (text[i] === ']') {
        i++
        return
      } else fail('commaOrBracket', i >= text.length ? i : lastEnd)
    }
  }

  try {
    value()
    ws()
    if (i < text.length) fail('trailing')
    return null
  } catch (e) {
    if (!(e instanceof Stop)) throw e
    let line = 1
    for (let k = 0; k < e.pos; k++) if (text.charCodeAt(k) === 10) line++
    return { kind: e.kind, pos: e.pos, line }
  }
}
