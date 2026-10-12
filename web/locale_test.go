package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/disp911/spacex-ui/v2/web/locale"

	"github.com/gin-gonic/gin"
)

// TestPagesRenderInTheRequestLanguage renders the panel's own templates for
// requests in both languages, one after another, as the panel serves them.
func TestPagesRenderInTheRequestLanguage(t *testing.T) {
	if err := locale.InitLocalizer(i18nFS); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	funcMap := template.FuncMap{"i18n": func(key string, params ...string) string { return locale.I18n(locale.Web, key, params...) }}
	engine.SetFuncMap(funcMap)
	tmpl, err := (&Server{}).getHtmlTemplate(funcMap)
	if err != nil {
		t.Fatal(err)
	}
	engine.SetHTMLTemplate(tmpl)
	engine.HTMLRender = locale.HTMLRender(engine.HTMLRender)
	engine.Use(locale.LocalizerMiddleware())
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", locale.PageData(c, gin.H{"base_path": "/", "title": "pages.login.title"}))
	})

	for _, c := range []struct{ lang, want, not string }{
		{"ru-RU", `<b class="is-visible">Добро пожаловать!</b>`, "placeholder='Username'"},
		{"en-US", `<b class="is-visible">Welcome</b>`, "Добро пожаловать!"},
		{"ru-RU", `<b class="is-visible">Добро пожаловать!</b>`, "placeholder='Username'"},
	} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "lang", Value: c.lang})
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		body := rec.Body.String()
		if rec.Code != http.StatusOK || !strings.Contains(body, c.want) || strings.Contains(body, c.not) {
			t.Fatalf("%s: HTTP %d, want %q without %q in the page", c.lang, rec.Code, c.want, c.not)
		}
	}
}
