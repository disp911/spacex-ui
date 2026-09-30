export interface SelectOption<T> {
  value: T
  label: string
  meta?: string
  /** CSS color of a status dot shown before the label. */
  dot?: string
}
