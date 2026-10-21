package spa

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func testEngine(fsys fstest.MapFS) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/base/v3/*path", Handler(fsys, func(c *gin.Context) Config {
		return Config{Base: "/base/", Mount: "/base/v3/", LoggedIn: c.GetHeader("X-Test-Login") == "1", Version: "3.0.0"}
	}))
	return engine
}

func get(engine *gin.Engine, target string, header ...string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

var built = fstest.MapFS{
	"index.html":        {Data: []byte(`<html><head><!--app-config--><script src="./assets/app-1.js"></script></head></html>`)},
	"assets/app-1.js":   {Data: []byte("console.log(1)")},
	"assets/app-1.css":  {Data: []byte("body{}")},
	"favicon.svg":       {Data: []byte("<svg/>")},
	".gitkeep":          {Data: nil},
	"assets/font.woff2": {Data: []byte("woff")},
}

func TestClientRoutesGetIndexWithConfig(t *testing.T) {
	engine := testEngine(built)
	for _, target := range []string{"/base/v3/", "/base/v3/inbounds", "/base/v3/settings/panel", "/base/v3/.gitkeep", "/base/v3/client.v2"} {
		rec := get(engine, target, "X-Test-Login", "1")
		body := rec.Body.String()
		if rec.Code != http.StatusOK || !strings.Contains(body, `<base href="/base/v3/">`) {
			t.Fatalf("%s: HTTP %d %q", target, rec.Code, body)
		}
		want := `content="{&#34;base&#34;:&#34;/base/&#34;,&#34;mount&#34;:&#34;/base/v3/&#34;,&#34;loggedIn&#34;:true,&#34;version&#34;:&#34;3.0.0&#34;}"`
		if !strings.Contains(body, want) || strings.Contains(body, "<!--app-config-->") {
			t.Fatalf("%s: config not rendered: %q", target, body)
		}
		if rec.Header().Get("Cache-Control") != "no-cache" {
			t.Fatalf("%s: index must not be cached, got %q", target, rec.Header().Get("Cache-Control"))
		}
	}
	if body := get(engine, "/base/v3/").Body.String(); !strings.Contains(body, "&#34;loggedIn&#34;:false") {
		t.Fatalf("anonymous request: %q", body)
	}
}

func TestAssetsAreServedAndCached(t *testing.T) {
	engine := testEngine(built)
	for target, typ := range map[string]string{
		"/base/v3/assets/app-1.js":   "text/javascript; charset=utf-8",
		"/base/v3/assets/app-1.css":  "text/css; charset=utf-8",
		"/base/v3/assets/font.woff2": "font/woff2",
	} {
		rec := get(engine, target)
		if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != typ || !strings.Contains(rec.Header().Get("Cache-Control"), "immutable") {
			t.Fatalf("%s: HTTP %d, %q, %q", target, rec.Code, rec.Header().Get("Content-Type"), rec.Header().Get("Cache-Control"))
		}
	}
	if rec := get(engine, "/base/v3/favicon.svg"); rec.Code != http.StatusOK || rec.Header().Get("Cache-Control") != "" {
		t.Fatalf("favicon: HTTP %d, cache %q", rec.Code, rec.Header().Get("Cache-Control"))
	}
	if rec := get(engine, "/base/v3/assets/missing-2.js"); rec.Code != http.StatusNotFound {
		t.Fatalf("missing asset: HTTP %d", rec.Code)
	}
	if rec := get(engine, "/base/v3/assets/../../../etc/passwd"); rec.Code == http.StatusOK && !strings.Contains(rec.Body.String(), "<base") {
		t.Fatalf("path escaped the frontend: %q", rec.Body.String())
	}
}

func TestUnbuiltFrontendSaysSo(t *testing.T) {
	engine := testEngine(fstest.MapFS{".gitkeep": {}})
	rec := get(engine, "/base/v3/")
	if rec.Code != http.StatusServiceUnavailable || !strings.Contains(rec.Body.String(), "npm --prefix frontend run build") {
		t.Fatalf("HTTP %d %q", rec.Code, rec.Body.String())
	}
}

func TestConfigIsEscaped(t *testing.T) {
	out := string(render([]byte("<head><!--app-config--></head>"), `/"><b>/`, Config{Base: `/"><script>x</script>/`}))
	if strings.Contains(out, "<script>") || strings.Contains(out, "<b>") {
		t.Fatalf("unescaped config: %s", out)
	}
}

func TestOtherPageGetsItsMountAndConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	fsys := fstest.MapFS{"sub.html": {Data: []byte(`<head><!--app-config--></head>`)}}
	engine.GET("/sub/:id", func(c *gin.Context) {
		ServePage(c, fsys, "sub.html", "/sub/", map[string]string{"subId": c.Param("id")})
	})
	body := get(engine, "/sub/abc").Body.String()
	if !strings.Contains(body, `<base href="/sub/">`) || !strings.Contains(body, `{&#34;subId&#34;:&#34;abc&#34;}`) {
		t.Fatalf("page: %q", body)
	}
}
