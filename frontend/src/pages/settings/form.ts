// What the settings page shares with its tabs: the settings being edited and
// what is changed or wrong in them.

import { inject, type InjectionKey } from 'vue'
import type { AllSetting, SettingKey } from './model'

export interface SettingsForm {
  /** The settings as edited; the tabs write into it. */
  readonly draft: AllSetting
  isDirty(key: SettingKey): boolean
  /** The message of a setting that fails its check. */
  error(key: SettingKey): string | undefined
  /** Changed and failing settings among these, for a section's badges. */
  count(keys: SettingKey[]): { dirty: number; errors: number }
}

export const formKey: InjectionKey<SettingsForm> = Symbol('settings form')

export function useSettingsForm() {
  const form = inject(formKey)
  if (!form) throw new Error('A settings tab needs the settings page around it')
  return form
}
