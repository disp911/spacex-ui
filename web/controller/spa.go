package controller

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"sync"

	"github.com/disp911/spacex-ui/v2/config"
	"github.com/disp911/spacex-ui/v2/logger"
	"github.com/disp911/spacex-ui/v2/web/service"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

// spaHeadMarker is the placeholder in the built index.html that the server
// replaces with the page's <base href> and its boot config.
const spaHeadMarker = "<!--spx-head-->"

// spaLanguages maps the codes the panel stores in its "lang" cookie to their
// translation files. The first entry is the fallback for missing keys.
var spaLanguages = []struct{ code, file string }{
	{"en-US", "translation/translate.en_US.toml"},
	{"ru-RU", "translation/translate.ru_RU.toml"},
}

var spa struct {
	dist fs.FS
	i18n fs.FS

	once     sync.Once
	messages map[string]map[string]string
}

// InitSPA hands the controllers the built frontend (web/dist) and the
// translation files. Pages fall back to the old templates while the frontend
// has not been built.
func InitSPA(dist fs.FS, i18n fs.FS) {
	spa.dist = dist
	spa.i18n = i18n
}

// spaIndex returns the built index.html, or nil when there is no build.
func spaIndex() []byte {
	if spa.dist == nil {
		return nil
	}
	b, err := fs.ReadFile(spa.dist, "index.html")
	if err != nil || !bytes.Contains(b, []byte(spaHeadMarker)) {
		return nil
	}
	return b
}

// spaConfig is the boot config the frontend reads from window.__SPX__.
type spaConfig struct {
	Page      string            `json:"page"`
	BasePath  string            `json:"basePath"`
	Version   string            `json:"version"`
	Host      string            `json:"host"`
	Lang      string            `json:"lang"`
	Languages []spaLanguage     `json:"languages"`
	TwoFactor bool              `json:"twoFactor"`
	Messages  map[string]string `json:"messages"`
}

type spaLanguage struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// renderSPA serves the new frontend for the given page. It reports false when
// the frontend is not built, so the caller can render the old template.
func renderSPA(c *gin.Context, page string, title string) bool {
	index := spaIndex()
	if index == nil {
		return false
	}

	lang := spaLang(c)
	cfg := spaConfig{
		Page:     page,
		BasePath: c.GetString("base_path"),
		Version:  config.GetVersion(),
		Host:     requestHost(c),
		Lang:     lang,
		Languages: []spaLanguage{
			{Code: "ru-RU", Name: "Русский"},
			{Code: "en-US", Name: "English"},
		},
		Messages: spaMessages(lang),
	}
	if page == "login" {
		var settingService service.SettingService
		cfg.TwoFactor, _ = settingService.GetTwoFactorEnable()
	}

	// json.Marshal escapes <, > and & so the config cannot close the script.
	data, err := json.Marshal(cfg)
	if err != nil {
		logger.Warning("Unable to encode frontend config:", err)
		return false
	}

	var head bytes.Buffer
	head.WriteString(`<base href="`)
	head.WriteString(escapeAttr(cfg.BasePath))
	head.WriteString(`"><title>`)
	head.WriteString(escapeAttr(cfg.Host + " – " + cfg.Messages[title]))
	head.WriteString(`</title><script>window.__SPX__=`)
	head.Write(data)
	head.WriteString(`</script>`)

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/html; charset=utf-8", bytes.Replace(index, []byte(spaHeadMarker), head.Bytes(), 1))
	return true
}

// spaLang picks the interface language the same way the old pages do: the
// "lang" cookie first, then the browser's Accept-Language header.
func spaLang(c *gin.Context) string {
	want, _ := c.Cookie("lang")
	if want == "" {
		want = c.GetHeader("Accept-Language")
	}
	want = strings.ToLower(want)
	for _, l := range spaLanguages {
		if strings.HasPrefix(want, strings.ToLower(l.code[:2])) {
			return l.code
		}
	}
	return spaLanguages[0].code
}

// spaMessages returns every web translation for lang as a flat
// "section.key" map, with English filling in keys the language lacks. The
// Telegram bot's strings are left out.
func spaMessages(lang string) map[string]string {
	spa.once.Do(func() {
		spa.messages = map[string]map[string]string{}
		var fallback map[string]string
		for i, l := range spaLanguages {
			msgs := map[string]string{}
			if spa.i18n != nil {
				if b, err := fs.ReadFile(spa.i18n, l.file); err == nil {
					var tree map[string]any
					if err := toml.Unmarshal(b, &tree); err != nil {
						logger.Warning("Unable to parse", l.file, err)
					}
					flattenMessages("", tree, msgs)
				}
			}
			if i == 0 {
				fallback = msgs
			} else {
				for k, v := range fallback {
					if _, ok := msgs[k]; !ok {
						msgs[k] = v
					}
				}
			}
			spa.messages[l.code] = msgs
		}
	})
	return spa.messages[lang]
}

func flattenMessages(prefix string, tree map[string]any, out map[string]string) {
	for k, v := range tree {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			out[key] = val
		case map[string]any:
			if key == "tgbot" {
				continue
			}
			flattenMessages(key, val, out)
		}
	}
}

func escapeAttr(s string) string {
	return strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;", `"`, "&quot;").Replace(s)
}
