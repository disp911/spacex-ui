package sub

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestPageFilesAreServedUnderTheLinksPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := fstest.MapFS{"assets/sub-1.js": {Data: []byte("x")}, "sub.html": {Data: []byte("<html>")}}
	engine := gin.New()
	engine.Use(pageFiles(app, "/s/"))
	engine.GET("/s/:subid", func(c *gin.Context) { c.String(http.StatusOK, "sub "+c.Param("subid")) })

	for target, want := range map[string]int{
		"/s/assets/sub-1.js":        http.StatusOK,
		"/s/assets/missing.js":      http.StatusNotFound,
		"/s/assets/../sub.html":     http.StatusNotFound,
		"/s/assets/%2e%2e/sub.html": http.StatusNotFound,
		"/s/abc":                    http.StatusOK,
		"/assets/sub-1.js":          http.StatusNotFound,
	} {
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		if rec.Code != want {
			t.Fatalf("%s: HTTP %d, want %d (%q)", target, rec.Code, want, rec.Body.String())
		}
	}
}
