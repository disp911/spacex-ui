package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Embedded fonts must be served whatever the host's MIME table knows:
// net/http sniffs their type by reading and seeking back, and a file it
// cannot seek comes back as 500.
func TestEmbeddedAssetsServeFonts(t *testing.T) {
	handler := http.FileServer(http.FS(&wrapAssetsFS{FS: assetsFS}))
	for _, path := range []string{
		"/fonts/ibm-plex/plex-sans-400-700-latin.woff2",
		"/fonts/ibm-plex/plex-mono-500-cyrillic.woff2",
		"/Vazirmatn-UI-NL-Regular.woff2",
	} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK || rec.Body.Len() == 0 {
			t.Errorf("GET %s = %d with %d bytes, want 200 with the font", path, rec.Code, rec.Body.Len())
		}
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/fonts/ibm-plex/plex-sans-400-700-latin.woff2", nil)
	req.Header.Set("Range", "bytes=0-3")
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusPartialContent || rec.Body.String() != "wOF2" {
		t.Errorf("ranged GET = %d %q, want 206 %q", rec.Code, rec.Body.String(), "wOF2")
	}
}
