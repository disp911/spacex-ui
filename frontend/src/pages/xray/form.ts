// What the Xray settings page shares with its tabs: the template being edited
// and what is changed in it.

import { inject, type InjectionKey } from 'vue'
import type { FieldKey, XraySetting } from './model'

export interface XrayForm {
  /** The template as edited; the tabs write into it. */
  readonly draft: XraySetting
  /** The template as the server has it. */
  readonly saved: XraySetting
  isDirty(key: FieldKey): boolean
  /** Changed fields among these, for a section's badge. */
  count(keys: FieldKey[]): number
  /** Puts the panel's default template in place of the edited one; saving keeps it. */
  resetToDefault(): Promise<void>
}

export const formKey: InjectionKey<XrayForm> = Symbol('xray form')

export function useXrayForm() {
  const form = inject(formKey)
  if (!form) throw new Error('An Xray settings tab needs the Xray settings page around it')
  return form
}
