import type { IconName } from './XIcon.vue'

export interface SelectOption<T> {
  value: T
  label: string
  meta?: string
  /** CSS color of a status dot shown before the label. */
  dot?: string
}

export interface MenuItem {
  key: string
  label: string
  icon?: IconName
  danger?: boolean
  disabled?: boolean
  href?: string
  /** Draw a divider above this item. */
  divided?: boolean
}
