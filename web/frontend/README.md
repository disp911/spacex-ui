# SpaceX panel — frontend (Vue 3)

The redesigned panel UI: Vue 3 + TypeScript + Vite. The Go server embeds the
build output from `web/ui` and serves it at `<webBasePath>ui/` (see
`web/spa.go`). The legacy template pages keep working at their usual paths
while screens are migrated.

Design rules, tokens and components: [DESIGN_SYSTEM.md](DESIGN_SYSTEM.md).

## Commands

```sh
npm ci
npm run dev        # dev server on :5173, API proxied to PANEL_URL (default http://127.0.0.1:2053)
npm run build      # typecheck + production build into ../ui
```

`web/ui` is committed, so `go build` works without Node. **Run
`npm run build` and commit `web/ui` whenever you change the frontend.**

With `XUI_DEBUG=true` the panel serves `web/ui` from disk instead of the
embedded copy, so a rebuild shows up without recompiling Go.

## Layout

```
src/
  styles/tokens.css     design tokens (both themes) — the only place with raw values
  components/ui/        design-system components (X*)
  components/shell/     app frame: sidebar, drawer, theme & language switches
  features/dashboard/   dashboard cards
  features/modals/      dialogs opened from the dashboard
  views/                routed pages (login, dashboard)
  api/                  typed wrappers over the panel HTTP API and WebSocket
  i18n/                 ru / en dictionaries
```
