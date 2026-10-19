import { shallowRef } from 'vue'
import type { IconName } from '../../kit'

/** A question window 19 of the mockup puts before an action: what it is about, what it asks and the action itself. */
export interface Ask {
  title: string
  text: string
  cta: string
  /** Cannot be undone, e.g. deleting or resetting traffic: the window shows it in red. */
  danger?: boolean
  /** What the action is, in place of the warning sign, e.g. info for a clone. */
  icon?: IconName
  run: () => Promise<unknown>
}

/** The question the Inbounds page shows now; the page holds the one window that asks it. */
export const asked = shallowRef<Ask | null>(null)

export function ask(question: Ask) {
  asked.value = question
}
