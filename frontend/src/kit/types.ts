export type OptionValue = string | number | boolean

export interface Option<T extends OptionValue = string> {
  value: T
  label: string
  disabled?: boolean
  /** What a cell of a grid list shows in place of the label, e.g. ␣ for a space. */
  glyph?: string
}

export type Tone = 'ok' | 'warn' | 'err' | 'unl' | 'on' | 'off' | 'neutral'
