package locale

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

// useTestBundle loads a two-language bundle for the test.
func useTestBundle(t *testing.T) {
	t.Helper()
	bundle := newBundle()
	for path, data := range map[string]string{
		"translation/translate.en_US.toml": "\"hello\" = \"Hello\"\n\"bye\" = \"Bye {{ .name }}\"\n",
		"translation/translate.ru_RU.toml": "\"hello\" = \"Привет\"\n\"bye\" = \"Пока, {{ .name }}\"\n",
	} {
		if _, err := bundle.ParseMessageFileBytes([]byte(data), path); err != nil {
			t.Fatal(err)
		}
	}
	bundleMu.Lock()
	old := i18nBundle
	i18nBundle = bundle
	bundleMu.Unlock()
	t.Cleanup(func() {
		bundleMu.Lock()
		i18nBundle = old
		bundleMu.Unlock()
	})
}

func testEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.SetFuncMap(template.FuncMap{"i18n": func(key string, params ...string) string { return I18n(Web, key, params...) }})
	engine.SetHTMLTemplate(template.Must(template.New("").Funcs(engine.FuncMap).Parse(
		`{{ define "greet" }}{{ i18n "hello" }}{{ end }}{{ define "page.html" }}{{ template "greet" }} {{ i18n "bye" "name==Ann" }}{{ end }}`)))
	engine.HTMLRender = HTMLRender(engine.HTMLRender)
	engine.Use(LocalizerMiddleware())
	engine.GET("/page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "page.html", PageData(c, gin.H{}))
	})
	engine.GET("/msg", func(c *gin.Context) {
		f, _ := c.Get("I18n")
		c.String(http.StatusOK, f.(func(I18nType, string, ...string) string)(Web, "hello"))
	})
	return engine
}

func get(engine *gin.Engine, path, lang string) string {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.AddCookie(&http.Cookie{Name: "lang", Value: lang})
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec.Body.String()
}

func TestRequestsKeepTheirOwnLanguage(t *testing.T) {
	useTestBundle(t)
	engine := testEngine(t)
	want := map[string][2]string{
		"ru-RU": {"Привет Пока, Ann", "Привет"},
		"en-US": {"Hello Bye Ann", "Hello"},
	}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		for lang, w := range want {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if got := get(engine, "/page", lang); got != w[0] {
					t.Errorf("%s page = %q, want %q", lang, got, w[0])
				}
				if got := get(engine, "/msg", lang); got != w[1] {
					t.Errorf("%s message = %q, want %q", lang, got, w[1])
				}
			}()
		}
	}
	wg.Wait()
}

func TestPagesWithoutALocalizerAreEnglish(t *testing.T) {
	useTestBundle(t)
	engine := testEngine(t)
	engine.GET("/plain", func(c *gin.Context) { c.HTML(http.StatusOK, "page.html", gin.H{}) })
	for range 3 {
		if got := get(engine, "/plain", "ru-RU"); got != "Hello Bye Ann" {
			t.Fatalf("got %q", got)
		}
	}
}
