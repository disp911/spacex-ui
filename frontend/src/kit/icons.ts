// Panel icons, drawn on a 16×16 grid with a currentColor stroke.
// Filled shapes set fill="currentColor" and stroke="none" themselves.

const circle = (cx: number, cy: number, r: number) => `<circle cx="${cx}" cy="${cy}" r="${r}"/>`
const dot = (cx: number, cy: number, r = 1.35) => `<circle cx="${cx}" cy="${cy}" r="${r}" fill="currentColor" stroke="none"/>`
const path = (d: string) => `<path d="${d}"/>`

export const icons = {
  plus: path('M8 3v10M3 8h10'),
  close: path('M4 4l8 8M12 4l-8 8'),
  check: path('M3.4 8.4 6.5 11.5 12.6 4.9'),
  chevron: path('M6 3.5 10.5 8 6 12.5'),
  chevronLeft: path('M10 3.5 5.5 8l4.5 4.5'),
  caret: path('M4 6l4 4 4-4'),
  menu: path('M2.4 4h11.2M2.4 8h11.2M2.4 12h7'),
  refresh: path('M13.4 8a5.4 5.4 0 1 1-1.6-3.8M13.6 2.2v3.2h-3.2'),
  search: circle(7, 7, 4.6) + path('M10.5 10.5l3.1 3.1'),
  filter: path('M2.4 3.2h11.2L9.3 8.4v4.2l-2.6 1.2V8.4z'),
  dots: dot(8, 3.4) + dot(8, 8) + dot(8, 12.6),
  grip: dot(5.5, 3.5, 1.25) + dot(10.5, 3.5, 1.25) + dot(5.5, 8, 1.25) + dot(10.5, 8, 1.25) + dot(5.5, 12.5, 1.25) + dot(10.5, 12.5, 1.25),
  edit: path('M10.6 2.9l2.5 2.5-7.7 7.7H2.9v-2.5z'),
  delete: path('M3.2 4.4h9.6M6.4 4.4V3.1h3.2v1.3M4.6 4.4l.6 8.5h5.6l.6-8.5'),
  info: circle(8, 8, 6.3) + path('M8 7.3v3.9M8 4.9v.1'),
  warning: path('M8 1.8 14.6 13.6H1.4z') + path('M8 6.4v3M8 11.6v.1'),
  user: circle(8, 5.2, 2.6) + path('M3 13.4c.6-2.5 2.6-3.9 5-3.9s4.4 1.4 5 3.9'),
  infinity: path('M8 8C6.6 6.1 5.4 5.3 4.3 5.3a2.7 2.7 0 0 0 0 5.4C5.4 10.7 6.6 9.9 8 8s2.6-2.7 3.7-2.7a2.7 2.7 0 0 1 0 5.4C10.6 10.7 9.4 9.9 8 8z'),
  pause: path('M5.5 3.5v9M10.5 3.5v9'),
  eye: path('M1.6 8S4 3.6 8 3.6 14.4 8 14.4 8 12 12.4 8 12.4 1.6 8 1.6 8z') + circle(8, 8, 2),
  eyeOff: path('M1.6 8S4 3.6 8 3.6 14.4 8 14.4 8 12 12.4 8 12.4 1.6 8 1.6 8z') + circle(8, 8, 2) + path('M3 13 13 3'),
  globe: circle(8, 8, 6.1) + '<ellipse cx="8" cy="8" rx="2.6" ry="6.1"/>' + path('M2 8h12'),
  bolt: path('M9.2 1.8 3.8 9h4l-1 5.2L12.2 7h-4z'),
  arrowRight: path('M2.5 8h9M8.5 4.5 12 8l-3.5 3.5'),
  balancer: path('M2 8h3.6M5.6 8l3-4h3.4M5.6 8l3 4h3.4') + dot(13.2, 4, 1.2) + dot(13.2, 12, 1.2),
  logout: path('M6.2 2.4H3.4A1.4 1.4 0 0 0 2 3.8v8.4a1.4 1.4 0 0 0 1.4 1.4h2.8M10.4 11.2 13.6 8l-3.2-3.2M13.2 8H6.4'),
  lock: '<rect x="3.4" y="7" width="9.2" height="7" rx="1.6"/>' + path('M5.4 7V5.2a2.6 2.6 0 0 1 5.2 0V7'),
  shield: path('M8 1.6 13.6 3.7v4.1c0 3.2-2.3 5.6-5.6 6.6-3.3-1-5.6-3.4-5.6-6.6V3.7z'),
  sun: circle(8, 8, 3) + path('M8 1.5v1.6M8 12.9v1.6M1.5 8h1.6M12.9 8h1.6M3.4 3.4l1.1 1.1M11.5 11.5l1.1 1.1M3.4 12.6l1.1-1.1M11.5 4.5l1.1-1.1'),
  moon: path('M13.2 9.8A5.6 5.6 0 0 1 6.2 2.8a5.6 5.6 0 1 0 7 7z'),
  dashboard: [[1.8, 1.8], [9.1, 1.8], [1.8, 9.1], [9.1, 9.1]].map(([x, y]) => `<rect x="${x}" y="${y}" width="5.1" height="5.1" rx="1.4"/>`).join(''),
  xray: path('M8 1.8l5.4 3.1v6.2L8 14.2l-5.4-3.1V4.9z') + circle(8, 8, 2),
  sliders: path('M2 4.6h12M2 11.4h12') + '<circle cx="10.2" cy="4.6" r="2.1" style="fill:var(--acc-soft)"/><circle cx="5.8" cy="11.4" r="2.1" style="fill:var(--acc-soft)"/>',
  doc: path('M3.4 2.4h9.2v11.2H3.4z') + path('M5.8 5.6h4.4M5.8 8h4.4M5.8 10.4h2.6'),
  download: path('M8 2.4v7.2M5.2 7l2.8 2.8L10.8 7M2.8 12.4h10.4'),
  upload: path('M8 9.6V2.4M5.2 5.2 8 2.4l2.8 2.8M2.8 12.4h10.4'),
  code: path('M6 3.4 2.8 8 6 12.6M10 3.4 13.2 8 10 12.6'),
  database: '<ellipse cx="8" cy="4" rx="5" ry="2"/>' + path('M3 4v8c0 1.1 2.2 2 5 2s5-.9 5-2V4M3 8c0 1.1 2.2 2 5 2s5-.9 5-2'),
  stop: '<rect x="3.2" y="3.2" width="9.6" height="9.6" rx="2"/>',
  arrowUp: path('M8 13.4V3M8 2.6l4.2 4.2M8 2.6 3.8 6.8'),
  arrowDown: path('M8 2.6V13M8 13.4l4.2-4.2M8 13.4 3.8 9.2'),
  alert: circle(8, 8, 6.9) + path('M8 4.6v4.1') + dot(8, 11.2, 0.95),
  swap: path('M2.6 5.4h10.4M10.4 2.8 13 5.4 10.4 8M13.4 10.6H3M5.6 8 3 10.6l2.6 2.6'),
} as const

export type IconName = keyof typeof icons
