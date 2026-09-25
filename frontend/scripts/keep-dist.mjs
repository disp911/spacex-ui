// vite build empties web/dist; put back the placeholder that keeps the
// directory in git, so the Go embed compiles without a frontend build.
import { writeFileSync } from 'node:fs'

writeFileSync(new URL('../../web/dist/.gitkeep', import.meta.url), '')
