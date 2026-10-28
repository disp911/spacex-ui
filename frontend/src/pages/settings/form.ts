// What the settings page shares with its tabs: the settings being edited and
// what is changed or wrong in them.

import { inject, type InjectionKey } from 'vue'
import { ApiError } from '../../api'
import type { AllSetting, SettingKey } from './model'

export interface SettingsForm {
  /** The settings as edited; the tabs write into it. */
  readonly draft: AllSetting
  /** The settings as the server has them. */
  readonly saved: AllSetting
  isDirty(key: SettingKey): boolean
  /** The message of a setting that fails its check. */
  error(key: SettingKey): string | undefined
  /** Changed and failing settings among these, for a section's badges. */
  count(keys: SettingKey[]): { dirty: number; errors: number }
  /**
   * Saves a few settings at once, apart from the other changes on the page,
   * e.g. two-factor authentication. Throws the server's refusal.
   */
  applyNow(patch: Partial<AllSetting>): Promise<void>
  /**
   * Restarts the panel as the page's own button does, without asking, and
   * goes on at the address given, e.g. where the website's settings put it.
   */
  restart(address?: { panel: string; page: string }): void
}

export const formKey: InjectionKey<SettingsForm> = Symbol('settings form')

export function useSettingsForm() {
  const form = inject(formKey)
  if (!form) throw new Error('A settings tab needs the settings page around it')
  return form
}

/** The server wraps the cause in parentheses after its own words; the cause is what helps. */
export function cause(e: unknown) {
  const msg = e instanceof ApiError ? e.message : String(e)
  return /\((.+)\)$/s.exec(msg)?.[1] ?? msg
}
