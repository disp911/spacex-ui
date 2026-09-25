# Handoff: 3x-ui panel redesign ("SpaceX panel")

## Overview
A full visual redesign of the 3x-ui (Xray-core) web admin panel: login, dashboard, Xray connections log, and all dashboard modals (Xray updates / geofiles / custom geo sources, panel update, panel log, backup, config.json viewer). Every screen exists in **dark and light themes** and in **desktop and mobile (390×844)** variants. UI copy is Russian (existing i18n keys should be reused; the Russian strings below are the reference wording).

## About the Design Files
The `.dc.html` files in this bundle are **design references built in HTML** — prototypes that show intended look and behavior, not production code to copy. Recreate them in the target codebase's existing environment (for 3x-ui: its Vue + Ant Design Vue frontend served by the Go backend) using its established patterns, components, and i18n. If you're building from scratch, choose an appropriate framework.

Open any file directly in a browser (keep `support.js` next to it). Each file is a canvas of artboards; every artboard has an id badge (e.g. `1a`, `5f`) and a caption. Interactive bits (dropdowns, accordions, filters, pagination, confirm dialogs, copy buttons) are wired with mock data.

## Fidelity
**High-fidelity.** Final colors, typography, spacing, radii, states and motion. Recreate pixel-accurately. All values below come straight from the files; when in doubt, inspect the inline styles in the HTML.

---

## Screens / Views

### 1. Login — `xpanel Login.dc.html`
Artboards: `1a` desktop dark 1280×840 · `1b` desktop light · `1e` mobile dark + light.

- **Background**: `#0B0E13` (light `#F3F5F7`) + radial glow `radial-gradient(120% 90% at 50% 118%, rgba(47,212,160,.16), transparent 62%)`. Bottom 440px: three stacked, horizontally-looping SVG sine waves (200% wide, `translateX(0 → -50%)`, linear infinite): A `#2FD4A0` op .30 26s; B `#4AA8E0` op .20 17s (reverse); C `#0F2E2A` op .42 11s. `pointer-events:none`.
- **Top bar** (padding 34px 44px): logo (32px rounded square `#2FD4A0`, rx 9, dark orbit mark `#06201A`) + "SpaceX **panel**" (600 15.5px, "panel" 450 `#8A91A0`). Right: language button (34px tall, globe icon + code + caret) and theme toggle (34×34, moon/sun cross-fade).
- **Language dropdown**: 232px wide, glass panel `rgba(19,23,31,.86)` + `backdrop-filter:blur(10px)`, border `rgba(255,255,255,.08)`, radius 12, shadow `0 18px 44px rgba(0,0,0,.45)`, header "ЯЗЫК ИНТЕРФЕЙСА" (mono 9.5px, tracking .09em). Rows: label + ISO code + check on selected. Opening it adds a full-screen scrim `rgba(8,11,16,.14)` + `blur(3px) saturate(.92)`. Animation `sheetUp .14s`.
- **Card**: 396px, centered, padding 30/32, radius 16, glass (same as dropdown), shadow `0 24px 60px rgba(0,0,0,.45)`.
  - Title "Вход" 600 19px, tracking -.4px. Sub "Сервер `nl-01.example.net`" (host in mono `#C7CCD6`).
  - Fields: labels mono 500 10.5px, tracking .07em, `#6E7787` ("ЛОГИН", "ПАРОЛЬ", "КОД 2FA"). Password has eye toggle.
  - 2FA (optional): 6 separate OTP cells (gap 7px); states filled / active (accent border) / empty.
  - Error banner: `rgba(242,85,90,.11)` bg, border `rgba(242,85,90,.32)`, radius 10, text `#FF9498` "Неверный логин или пароль", dismiss ✕.
  - Primary button: bg `#2FD4A0`, text `#06201A`; hover bg `#45E0AF` + `0 0 0 4px rgba(47,212,160,.16), 0 0 22px rgba(47,212,160,.4)`. States: idle → loading (18px spinner, `spin .7s linear`) → success (check drawn with `stroke-dashoffset` .34s) ; press = `btnPop` scale .975.
  - Footer (border-top `rgba(255,255,255,.07)`): pulsing dot + "Пинг · 24 мс".

### 2. Dashboard — `xpanel Dashboard.dc.html`
Artboards: `1a` desktop dark 1440×900 · `2a` desktop light · `3a`/`3b` mobile dark/light 390×844.

**Desktop layout**: flex row. Sidebar 236px (`#0E1218`, border-right `rgba(255,255,255,.06)`, padding 20/14/16) + content (flex 1, padding 22/26/24, scrolls).
- **Sidebar**: logo 26px + "SpaceX panel" 600 14.5px. Nav vertically centered: items 38px tall, radius 10, gap 11, icon 16px, 450 13px `#9AA3B2`; hover `rgba(255,255,255,.05)` / `#E6E9EF`; active bg `#172A26`, text `#2FD4A0`, weight 500. Items: Дашборд, Подключения, Настройки, Настройки Xray. Bottom: Xray status pill (`#111520`, pulsing 7px dot in state color, "XRAY" mono + state label) and "Выход" (hover red `rgba(242,85,90,.1)` / `#FF9498`).
- **Header** (45px): "Дашборд" 600 21px + "Сервер nl-01.example.net". Right side: HTTP warning banner, width `calc(50% - 7px)` so it aligns with the right column — amber `rgba(232,161,58,.1)`, border `.3`, radius 11, title `#F0B75C` "Соединение не защищено (HTTP)", body `#C9A473`, dismiss ✕.
- **Grid**: 2 equal columns, gap 12. Cards: bg `rgba(19,23,31,.86)`, border `rgba(255,255,255,.08)`, radius 16, padding 14/20. Card header row: mono caps label (10.5px, `#6E7787`, tracking .07em) + right chip (mono 13px, `#161B24` bg, radius 8, 24px tall), separated by 1px `rgba(255,255,255,.07)` rule.
  - **Left column**
    - *Xray* card: version chip `v26.4.25`; state row = 9px pulsing dot + state label 600 18px in state color; metrics row (КЛИЕНТЫ ОНЛАЙН 37 · ВРЕМЯ РАБОТЫ 3d 4h · ПАМЯТЬ 48.20 MB — label mono 9.5px, value mono 500 16px) with icon buttons Restart (neutral) and Stop (red-tinted) 34×34 radius 10. Error state shows a mono red log box (max-height 76px). Link rows (34px): "Журнал Xray ›", "Выбор версии ›".
    - *Telegram-прокси* card (optional, prop `telegramProxy`): same structure, version 3.5.7, own state & error text.
    - *Управление* card: rows "Журнал панели", "Конфигурация", "Резервная копия" (each opens a modal) + "Версия панели v2.10.0" with update CTA "Обновить панель: v2.10.1" when an update is available.
    - *IP-адрес сервера*: IPv4 / IPv6 values, blurred (`filter:blur(5px)`) until the eye toggle is clicked.
  - **Right column**
    - *Ресурсы сервера*: ЦП / ОЗУ / ПОДКАЧКА / ДИСК rows with detail text + percentage + bar.
    - *Загрузка ЦП*: area chart (28 points, SVG path, y-ticks every 25%, auto top = ceil(max·1.15/25)·25, min 25), "пик N%" label, range dropdown "за 2 мин" with options 2 мин / 30 мин / 1 час / 2 часа / 3 часа / 5 часов (selected row `#172A26`/`#2FD4A0` + check). Dropdown open → scrim blur like login.
    - *Сеть*: Всего отправлено / получено, Отправка / Загрузка (per second), TCP / UDP counts.
    - *Система*: ОС uptime, load average 1/5/15 with helper text.
- **Status colors** (dark / light): running "Запущен" `#2FD4A0`/`#0A7A59`; stopped "Остановлен" `#E8A13A`/`#B4761A`; error "Ошибка" `#F2555A`/`#C42A31`; unknown "Неизвестно" `#6E7787`/`#686F7D`.

**Mobile**: header 72px (`#0E1218`): burger button 36×36 radius 11 + title 600 17px + host mono 11px. Content scrolls (hidden scrollbar), padding 13/14/16, gap 12, single column in the same card order; Restart/Stop move to the state row (36×36). Burger opens a left drawer 282px (`drawerIn .24s cubic-bezier(.22,.7,.3,1)`) over a scrim `rgba(6,9,13,.58)` + blur(2px); nav items 46px tall, 14px text.

### 3. Connections log — `xpanel Connections Log.dc.html`
Artboards: `1a` desktop dark · `2a` desktop light (wide modal ~80% of viewport) · `3a` mobile dark · `4a` mobile light.

- **Desktop modal**: title "Журнал соединений Xray", subtitle with counts, close ✕. Filter bar: ДАТА dropdown (days with record counts, last 7 days), КЛИЕНТ dropdown (all / per email), search input (IP, домен, email…; clear ✕), connection-type toggle chips (direct / blocked / proxy), "Скачать" button.
- **Table** columns: ДАТА · ОТКУДА · КУДА · ПОДКЛЮЧЕНИЕ (inbound) · ИСХОДЯЩИЙ (outbound) · EMAIL. Mono values, tabular numbers. Optional row tint by outbound type (prop `rowTint`), optional dense rows (`denseRows`). Only the table body scrolls.
- Empty state: "Ничего не найдено" / "Смягчите фильтры или выберите другой день".
- Footer: range label ("1–50 из 6 200") + pagination (‹ pages ›).
- **Mobile**: full-screen sheet; search field (36px) + three equal filter buttons (Day / Client / Types, 34px). Each filter opens a **bottom sheet** (white/`rgba(19,23,31,.96)`, radius 22 22 0 0, grab handle 38×4, `sheetUp .16s`). Rows render as cards instead of a table.

### 4. Modals — `xpanel Modals.dc.html`
Artboards: `1*` desktop dark · `4*` desktop light · `5*` mobile dark · `6*` mobile light.

Desktop dialog shell: bg `rgba(19,23,31,.96)` (light `#FFFFFF`), border `rgba(255,255,255,.08)`, radius 16, shadow `0 24px 60px rgba(0,0,0,.45)`. Widths: Updates 596, Add/Edit source 464, Panel update 436, Panel log 800, Backup 436, confirm 400. Mobile: full-screen sheets; confirmations become bottom sheets.

- **Обновления Xray** (`1a`): subtitle "Ядро, geo-базы и пользовательские источники". Three accordion sections (one open at a time; caret ▸ rotates):
  1. *Xray* — amber "Важно" notice ("Старые версии могут не поддерживать текущие настройки. После смены версии ядро будет перезапущено."), scrollable version list (max-height 290) with "установлена" chip on current and release date. Picking a version → confirm dialog.
  2. *Geofiles* (6 файлов) — info note, "Обновить все", table ФАЙЛ / ОБНОВЛЕНО / ОБНОВИТЬ with per-row refresh icon → confirm.
  3. *Пользовательские GeoSite / GeoIP* — hint "ext:файл.dat:тег", "Обновить все", "+ Добавить", table ТИП (geosite/geoip badge) / ПСЕВДОНИМ / МАРШРУТИЗАЦИЯ / ОБНОВЛЕНО / ДЕЙСТВИЯ (edit, refresh, delete-red). Empty state: dashed icon + "Пользовательских источников нет".
  - Confirm dialog: title/text/CTA vary by action; CTA shows busy label + spinner; in/out animations `cDlgIn`/`cDlgOut` (desktop) and `cSheetIn`/`cSheetOut` (mobile).
- **Добавить источник** (`1d`): ТИП segmented (geosite | geoip), ПСЕВДОНИМ (hint "Разрешены a-z 0-9 _ - · войдёт в имя файла"), URL, live preview "В ПРАВИЛАХ ext:geosite_myads.dat:тег". Отмена / Сохранить.
- **Изменить источник** (`1e`): subtitle = filename; type and alias locked ("Псевдоним нельзя изменить — он привязан к имени файла"); URL error state "Укажите полный адрес, начинающийся с https://" (red border + message); "Обновлён 18.09.2026 04:12".
- **Обновить панель** (`1f`): text, Текущая v2.10.0 → Последняя v2.10.1, CTA "Обновить панель" → "Обновление…" with spinner.
- **Журнал панели** (`1h`): subtitle "Последние N строк · journalctl -u x-ui". Filters СТРОК (dropdown), УРОВЕНЬ (dropdown), SysLog toggle, "Скачать". Mono log lines: timestamp, level (colored), source, message. Only the list scrolls. Empty state "Нет записей". Mobile: filters scroll horizontally.
- **Резервная копия** (`1k`): two action rows — "Экспорт базы данных" (download .db) and "Импорт базы данных" (upload .db) with descriptions.
- **Конфигурация** (`1l`): read-only config.json viewer with line-number gutter, status footer "JSON валиден · N стр." or the parse error (prop `configError`), "Копировать" → "Скопировано" feedback.

---

## Interactions & Behavior
- **Hover glow (signature)** on outlined buttons/inputs: border `rgba(47,212,160,.5)`, bg `#1A2B28`, shadow `0 0 0 3px rgba(47,212,160,.12), 0 0 14px rgba(47,212,160,.3)`, text `#E6E9EF`. Light: border `rgba(10,122,89,.45)`, bg `#DCF1E8`, shadow `0 0 0 3px rgba(10,122,89,.13), 0 0 14px rgba(10,122,89,.22)`. Transition `box-shadow/border-color/background .22s ease`.
- Danger buttons: `rgba(242,85,90,.08)` bg, `.3` border, `#FF9498` text; hover bg `.16`, border `.5`, ring `0 0 0 3px rgba(242,85,90,.12)`.
- Nav/list rows: `background, color .18s ease`.
- Popovers/dropdowns: full-area scrim with `backdrop-filter: blur(3px) saturate(.92)`; click scrim to close; enter `sheetUp` (translateY 24px→0, fade) .14–.16s ease-out.
- Status dots: `pulseDot` 2.4s infinite (opacity 1 → .35).
- Spinner: `spin .7s linear infinite`.
- Keyframes (copy verbatim from `<helmet><style>` in each file): `pulseDot, spin, fadeIn, drawerIn, sheetUp, scrimIn, scrimOut, checkDraw, btnPop, cSheetIn/Out, cDlgIn/Out, waveDriftA/B/C`.
- Scrollbars: 9px, thumb `#2A3141` (hover `#3A4352`), light `#C9CED7`/`#AAB1BC`, rounded, no arrows. Mobile scroll areas hide scrollbars.
- Numbers: `font-variant-numeric: tabular-nums` on all app frames.

## State Management
- Theme (dark/light) + language — persisted per user.
- Login: `username, password, otp[6], showPassword, status: idle|loading|success|error`, `twoFactorEnabled` (from server).
- Dashboard: `xrayState, telegramState: running|stopped|error|unknown`, error texts, `telegramProxyEnabled`, `httpWarningDismissed`, `updateAvailable`, `ipHidden`, `cpuRange: 2m|30m|1h|2h|3h|5h`, `mobileMenuOpen`. Server status polled (existing 3x-ui `/server/status`).
- Connections log: `day, client, query, kinds{direct,blocked,proxy}, page, openMenu, loading`.
- Modals: `openSection: xray|geo|custom`, `selectedVersion`, `confirm: {kind, target} | null`, `busy`, `panelLog{lines, level, syslog}`, `configText, copied`, source form `{type, alias, url, errors}`.

## Design Tokens
Font: **IBM Plex Sans** (400/450/500/600/700) for UI, **IBM Plex Mono** (400/500/600) for labels, values, hosts, logs, code (Google Fonts).

Typography scale: page title 600 21px/-.5px (mobile 17px) · card state 600 18px · login title 600 19px · dialog title ~16.5px 600 · body 12–13.5px · metric value mono 500 16px · caps label mono 500 9.5–10.5px, tracking .07em · small mono 10.5–11.5px.

Radii: 3–5 (chips/badges) · 7–9 (small buttons/rows) · 10–12 (buttons, inputs, notices, dropdowns) · 16 (cards, dialogs) · 22 (bottom sheets top corners) · 26 (mobile frame).

Shadows: card/dialog `0 24px 60px rgba(0,0,0,.45)` (light `rgba(16,24,40,.16)`); dropdown `0 18px 44px rgba(0,0,0,.45)`; drawer `26px 0 70px rgba(0,0,0,.5)`; bottom sheet (light) `0 -24px 70px rgba(16,24,40,.18)`.

Dark → light color map (used consistently across all files — implement as theme tokens):
```
Surfaces
#0B0E13 → #F3F5F7   app background
#0E1218 → #EFF1F4   sidebar / mobile header
rgba(19,23,31,.86)  card (glass)          → white
rgba(19,23,31,.96)  dialog                → #FFFFFF
#111520 → #F5F6F8   button / field fill
#161B24 → #F5F6F8   chip / inset fill
#14181F → #EFF1F4   #1A2028 → #EDEFF3 (row hover)
#1C2130, #242938 → #E3E6EB   borders (solid)
#2A3141 → #C9CED7   #3A4150 → #B4BAC4
rgba(255,255,255,.06/.07/.08/.1) → rgba(16,20,28,.08/.07/.09/.12)  hairlines
Text
#E6E9EF → #13171E   primary
#C7CCD6 → #39404D   secondary
#9AA3B2, #8A91A0 → #5C6371   tertiary
#6E7787 → #686F7D   muted / labels
#4E5563 → #8A909C   faint
Accent (green)
#2FD4A0 → #0A7A59   primary accent     hover #45E0AF → #086A4D
#06201A → #FFFFFF   text on accent
#172A26 → #E4F5EE   active nav / selected row
#1A2B28 → #DCF1E8   hover fill
Danger   #F2555A → #C42A31  · text #FF9498 → #C22D34 · #FF7B7F → #B23037
Warning  #E8A13A → #8A5A10 (dashboard state uses #B4761A) · title #F0B75C → #8A5A10 · body #C9A473 → #7A6337
Info     #4AA8E0 → #15628F (also #1E7FC0) · #6FC0F0 → #15628F
Tinted fills keep alpha and swap base: rgba(47,212,160,a) → rgba(10,122,89,a'), rgba(242,85,90,a) → rgba(196,42,49,a'), rgba(232,161,58,a) → rgba(180,118,26,a'), rgba(74,168,224,a) → rgba(30,127,192,a')
```
Common spacing: 12px grid/card gaps; card padding 14×20 (desktop) / 13×15 (mobile); row heights 34 (desktop link rows) / 38 (desktop nav) / 46 (mobile nav); icon buttons 34 (desktop) / 36 (mobile); touch targets ≥ 32px.

## Assets
- Logo: inline SVG (rounded square `#2FD4A0` + orbit ellipse + two dots in `#06201A`) — in every file's sidebar/top bar. "SpaceX panel" is the working product name.
- Icons: simple 16×16 stroked inline SVGs (stroke 1.3–1.6, round caps), `currentColor`. Map to the codebase's icon set if one exists, otherwise copy from the HTML.
- Fonts: IBM Plex Sans / Mono via Google Fonts.
- All data (hosts, IPs `203.0.113.10` / `2001:db8::1`, versions, logs, config) is mock.

## Files
- `xpanel Login.dc.html` — login (desktop/mobile, dark/light)
- `xpanel Dashboard.dc.html` — dashboard (desktop dark `1a`, desktop light `2a`, mobile `3a`/`3b`); tweak props: `xrayState`, `telegramProxy`, `telegramState`, `httpWarning`, `updateAvailable`
- `xpanel Connections Log.dc.html` — Xray connections log (desktop/mobile, dark/light); props `rowTint`, `denseRows`
- `xpanel Modals.dc.html` — all dashboard modals (desktop/mobile, dark/light); props `configError`, `logEmpty`, `sourcesEmpty`
- `support.js` — runtime needed only to open the reference files in a browser; not part of the implementation.
