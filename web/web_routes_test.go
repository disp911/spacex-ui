package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUnroutedPagesGoToTheFrontendUnderTheBasePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.POST("/secret/login", func(c *gin.Context) { c.String(http.StatusOK, "login endpoint") })
	app := func(c *gin.Context) { c.String(http.StatusOK, "app "+c.Param("path")) }
	decoy := func(c *gin.Context) { c.String(http.StatusOK, "decoy") }
	engine.NoRoute(frontendOrDecoy("/secret/", app, decoy))

	for _, tc := range []struct{ method, target, want string }{
		{http.MethodGet, "/secret/", "app /"},
		{http.MethodGet, "/secret/login", "app /login"},
		{http.MethodGet, "/secret/inbounds", "app /inbounds"},
		{http.MethodGet, "/secret/assets/index-1.js", "app /assets/index-1.js"},
		{http.MethodGet, "/secret/panel/api/missing", "decoy"},
		{http.MethodGet, "/", "decoy"},
		{http.MethodGet, "/other", "decoy"},
		{http.MethodPost, "/secret/login", "login endpoint"},
		{http.MethodPost, "/secret/missing", ""},
	} {
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.target, nil))
		if rec.Body.String() != tc.want {
			t.Fatalf("%s %s: %q, want %q", tc.method, tc.target, rec.Body.String(), tc.want)
		}
	}
}
