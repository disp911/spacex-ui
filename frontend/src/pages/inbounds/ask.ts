import { shallowRef } from 'vue'

/** A question window 19 of the mockup puts before an action: what it is about, what it asks and the action itself. */
export interface Ask {
  title: string
  text: string
  cta: string
  /** Cannot be undone, e.g. deleting or resetting traffic: the window shows it in red. */
  danger?: boolean
  run: () => Promise<unknown>
}

/** The question the Inbounds page shows now; the page holds the one window that asks it. */
export const asked = shallowRef<Ask | null>(null)

export function ask(question: Ask) {
  asked.value = question
}
