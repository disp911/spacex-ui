/**
 * Chrome shows a textarea's text cursor over its styled scrollbars too. While
 * the pointer is over a bar the textarea is marked x-on-scrollbar, and kit.css
 * gives it the arrow there.
 */
export function scrollbarCursor(e: PointerEvent) {
  const el = e.currentTarget as HTMLElement
  const box = el.getBoundingClientRect()
  const x = e.clientX - box.left - el.clientLeft
  const y = e.clientY - box.top - el.clientTop
  const onBar = e.type === 'pointermove' && (x >= el.clientWidth || y >= el.clientHeight)
  el.classList.toggle('x-on-scrollbar', onBar)
}
