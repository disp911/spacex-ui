package web

import (
	"embed"
	"encoding/json"
	"html"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/disp911/spacex-ui/v2/config"
	"github.com/disp911/spacex-ui/v2/logger"

	"github.com/gin-gonic/gin"
)

// The redesigned Vue 3 panel is built from web/frontend into web/ui
// (`npm run build`) and served under <basePath>ui/. It runs alongside the
// legacy template pages while screens are migrated.
//
//go:embed all:ui
var uiFS embed.FS

// spaBootMarker is replaced in index.html with the runtime configuration,
// so one build works under any base path.
const spaBootMarker = "<!--spx:boot-->"

// registerSPA serves the built single-page app: static files as-is, any
// other path under /ui/ as index.html so client-side routes survive reloads.
func registerSPA(g *gin.RouterGroup, basePath string) {
	var files fs.FS
	if config.IsDebug() {
		files = os.DirFS("web/ui")
	} else {
		sub, err := fs.Sub(uiFS, "ui")
		if err != nil {
			logger.Warning("SPA: embedded ui not found:", err)
			return
		}
		files = sub
	}

	g.GET("/ui/*path", func(c *gin.Context) {
		name := strings.TrimPrefix(path.Clean("/"+c.Param("path")), "/")

		if name != "" && name != "index.html" {
			if data, err := fs.ReadFile(files, name); err == nil {
				if strings.HasPrefix(name, "assets/") {
					// Vite fingerprints asset names, so they never change.
					c.Header("Cache-Control", "public, max-age=31536000, immutable")
				}
				ctype := mime.TypeByExtension(path.Ext(name))
				if ctype == "" {
					ctype = http.DetectContentType(data)
				}
				c.Data(http.StatusOK, ctype, data)
				return
			}
			if strings.HasPrefix(name, "assets/") {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		index, err := fs.ReadFile(files, "index.html")
		if err != nil {
			c.String(http.StatusNotFound, "UI is not built. Run `npm run build` in web/frontend.")
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", injectSPABoot(index, basePath, requestHost(c)))
	})
}

// injectSPABoot writes the <base> tag and window.__SPX__ into index.html.
func injectSPABoot(index []byte, basePath string, host string) []byte {
	boot, _ := json.Marshal(map[string]string{
		"basePath": basePath,
		"host":     host,
		"version":  config.GetVersion(),
	})
	// json.Marshal escapes <, > and & so the payload cannot close the script.
	snippet := `<base href="` + html.EscapeString(basePath+"ui/") + `">` +
		`<script>window.__SPX__=` + string(boot) + `</script>`
	return []byte(strings.Replace(string(index), spaBootMarker, snippet, 1))
}

// requestHost mirrors the host shown by the legacy templates.
func requestHost(c *gin.Context) string {
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.GetHeader("X-Real-IP")
	}
	if host == "" {
		if h, _, err := net.SplitHostPort(c.Request.Host); err == nil {
			host = h
		} else {
			host = c.Request.Host
		}
	}
	return host
}
