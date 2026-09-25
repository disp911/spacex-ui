package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func setupSPA(t *testing.T, index string) {
	t.Helper()
	dist := fstest.MapFS{}
	if index != "" {
		dist["index.html"] = &fstest.MapFile{Data: []byte(index)}
	}
	i18n := fstest.MapFS{
		"translation/translate.en_US.toml": &fstest.MapFile{Data: []byte(
			"\"close\" = \"Close\"\n[pages.index]\n\"title\" = \"Dashboard\"\n[tgbot]\n\"hello\" = \"hi\"\n",
		)},
		"translation/translate.ru_RU.toml": &fstest.MapFile{Data: []byte(
			"[pages.index]\n\"title\" = \"Дашборд\"\n",
		)},
	}
	InitSPA(dist, i18n)
	spa.once = sync.Once{}
	spa.messages = nil
	t.Cleanup(func() { InitSPA(nil, nil) })
}

func renderPage(t *testing.T, header http.Header) (*httptest.ResponseRecorder, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/secret/panel/", nil)
	for k, v := range header {
		c.Request.Header[k] = v
	}
	c.Request.Host = "panel.example.net:2053"
	c.Set("base_path", "/secret/")
	ok := renderSPA(c, "dashboard", "pages.index.title")
	return w, ok
}

func bootConfig(t *testing.T, body string) spaConfig {
	t.Helper()
	start := strings.Index(body, "window.__SPX__=")
	end := strings.Index(body, "</script>")
	if start < 0 || end < start {
		t.Fatalf("no boot config in %q", body)
	}
	var cfg spaConfig
	if err := json.Unmarshal([]byte(body[start+len("window.__SPX__="):end]), &cfg); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestRenderSPA_InjectsBaseAndConfig(t *testing.T) {
	setupSPA(t, "<head><!--spx-head--></head><body></body>")
	w, ok := renderPage(t, http.Header{"Cookie": {"lang=ru-RU"}})
	if !ok {
		t.Fatal("renderSPA reported no build")
	}
	body := w.Body.String()
	if !strings.Contains(body, `<base href="/secret/">`) {
		t.Errorf("missing base href: %s", body)
	}
	if !strings.Contains(body, "<title>panel.example.net – Дашборд</title>") {
		t.Errorf("unexpected title: %s", body)
	}
	cfg := bootConfig(t, body)
	if cfg.Page != "dashboard" || cfg.BasePath != "/secret/" || cfg.Lang != "ru-RU" {
		t.Errorf("unexpected config %+v", cfg)
	}
	// Russian lacks "close", so English fills it in; bot strings stay out.
	if cfg.Messages["close"] != "Close" {
		t.Errorf("fallback message missing: %v", cfg.Messages)
	}
	if _, ok := cfg.Messages["tgbot.hello"]; ok {
		t.Error("telegram bot strings leaked into the web messages")
	}
}

func TestRenderSPA_ConfigCannotCloseTheScript(t *testing.T) {
	setupSPA(t, "<head><!--spx-head--></head>")
	w, _ := renderPage(t, http.Header{"X-Forwarded-Host": {"</script><script>alert(1)</script>"}})
	body := w.Body.String()
	if strings.Count(body, "</script>") != 1 {
		t.Errorf("host escaped the boot script: %s", body)
	}
}

func TestRenderSPA_FallsBackWithoutBuild(t *testing.T) {
	setupSPA(t, "")
	if _, ok := renderPage(t, nil); ok {
		t.Error("renderSPA should report false when index.html is missing")
	}
}

func TestSpaLang(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		cookie, accept, want string
	}{
		{"ru-RU", "", "ru-RU"},
		{"en-US", "ru", "en-US"},
		{"", "ru-RU,ru;q=0.9", "ru-RU"},
		{"", "de-DE", "en-US"},
		{"fa-IR", "", "en-US"},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		if tc.cookie != "" {
			c.Request.AddCookie(&http.Cookie{Name: "lang", Value: tc.cookie})
		}
		if tc.accept != "" {
			c.Request.Header.Set("Accept-Language", tc.accept)
		}
		if got := spaLang(c); got != tc.want {
			t.Errorf("spaLang(cookie=%q, accept=%q) = %q, want %q", tc.cookie, tc.accept, got, tc.want)
		}
	}
}
