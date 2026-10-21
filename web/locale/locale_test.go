package locale

import (
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
	engine.Use(LocalizerMiddleware())
	engine.GET("/msg", func(c *gin.Context) {
		f, _ := c.Get("I18n")
		i18n := f.(func(I18nType, string, ...string) string)
		c.String(http.StatusOK, i18n(Web, "hello")+" "+i18n(Web, "bye", "name==Ann"))
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
	want := map[string]string{"ru-RU": "Привет Пока, Ann", "en-US": "Hello Bye Ann"}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		for lang, w := range want {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if got := get(engine, "/msg", lang); got != w {
					t.Errorf("%s message = %q, want %q", lang, got, w)
				}
			}()
		}
	}
	wg.Wait()
}
