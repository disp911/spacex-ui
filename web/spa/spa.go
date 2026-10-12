// Package spa serves the panel's single-page frontend built by Vite.
package spa

import (
	"bytes"
	"encoding/json"
	"html"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Config is handed to the frontend in index.html.
type Config struct {
	// Base is the panel's base path, e.g. "/" or "/secret/".
	Base string `json:"base"`
	// Mount is where the frontend is served, ending with a slash.
	Mount string `json:"mount"`
	// LoggedIn tells the frontend whether the request carried a session.
	LoggedIn bool   `json:"loggedIn"`
	Version  string `json:"version"`
}

// placeholder is the comment in frontend/index.html that the server replaces
// with the <base> element and the config.
const placeholder = "<!--app-config-->"

// Handler serves fsys: files under assets/ with a year-long cache, any other
// path with index.html, so that client-side routes survive a reload. The
// route must be registered with a *path wildcard.
func Handler(fsys fs.FS, config func(c *gin.Context) Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.TrimPrefix(path.Clean("/"+c.Param("path")), "/")
		if name != "" && !strings.HasPrefix(path.Base(name), ".") && path.Ext(name) != "" {
			if data, err := fs.ReadFile(fsys, name); err == nil {
				if strings.HasPrefix(name, "assets/") {
					// Vite puts a content hash in every asset name.
					c.Header("Cache-Control", "public, max-age=31536000, immutable")
				}
				c.Data(http.StatusOK, contentType(name), data)
				return
			}
			if strings.HasPrefix(name, "assets/") {
				c.Status(http.StatusNotFound)
				return
			}
		}

		index, err := fs.ReadFile(fsys, "index.html")
		if err != nil {
			c.String(http.StatusServiceUnavailable, "The panel frontend is not built: run npm --prefix frontend run build")
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", render(index, config(c)))
	}
}

// render puts the <base> element and the config into index.html.
func render(index []byte, cfg Config) []byte {
	data, _ := json.Marshal(cfg)
	head := `<base href="` + html.EscapeString(cfg.Mount) + `">` +
		`<meta name="app-config" content="` + html.EscapeString(string(data)) + `">`
	return bytes.Replace(index, []byte(placeholder), []byte(head), 1)
}

func contentType(name string) string {
	switch path.Ext(name) {
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	case ".json", ".webmanifest":
		return "application/json"
	case ".png":
		return "image/png"
	case ".ico":
		return "image/x-icon"
	case ".html":
		return "text/html; charset=utf-8"
	}
	return "application/octet-stream"
}
