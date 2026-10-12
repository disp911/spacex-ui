export type OptionValue = string | number | boolean

export interface Option<T extends OptionValue = string> {
  value: T
  label: string
  disabled?: boolean
}

export type Tone = 'ok' | 'warn' | 'err' | 'unl' | 'on' | 'off' | 'neutral'
