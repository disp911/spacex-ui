package web

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// webProxyTransport carries WEB proxy requests to telemt on loopback. It
// ignores proxy settings from the environment and never compresses, so the
// bytes telemt sends reach Telegram unchanged.
var webProxyTransport = &http.Transport{
	Proxy:               nil,
	DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
	MaxIdleConns:        256,
	MaxIdleConnsPerHost: 256,
	IdleConnTimeout:     90 * time.Second,
	DisableCompression:  true,
}

// webProxyBase returns the first segment of a path that has more after it:
// "abc" for "/abc/" and "/abc/api/v1/up", nothing for "/abc".
func webProxyBase(path string) (string, bool) {
	rest, ok := strings.CutPrefix(path, "/")
	if !ok {
		return "", false
	}
	base, _, ok := strings.Cut(rest, "/")
	if !ok || base == "" {
		return "", false
	}
	return base, true
}

// webProxyFront is the panel's HTTPS handler. A request under the secret
// path of a running WEB proxy goes to that proxy's telemt listener; it is
// passed on before the panel's middleware, whose compression and sessions
// would break the long polls and WebSockets Telegram uses. Everything else,
// the bare secret path included, is the panel's.
type webProxyFront struct {
	panel http.Handler
	// route returns the loopback port of the WEB proxy at a base path.
	route func(basePath string) (int, bool)
}

func (f webProxyFront) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base, ok := webProxyBase(r.URL.EscapedPath())
	if !ok {
		f.panel.ServeHTTP(w, r)
		return
	}
	port, ok := f.route(base)
	if !ok {
		f.panel.ServeHTTP(w, r)
		return
	}
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.Out.URL.Scheme = "http"
			pr.Out.URL.Host = fmt.Sprintf("127.0.0.1:%d", port)
			// telemt matches the public host name and trusts exactly one
			// client address, which replaces whatever the client sent.
			pr.Out.Host = pr.In.Host
			if ip, _, err := net.SplitHostPort(pr.In.RemoteAddr); err == nil {
				pr.Out.Header.Set("X-Forwarded-For", ip)
			}
		},
		Transport:     webProxyTransport,
		FlushInterval: -1,
		// Should telemt be down, the visitor sees the site, not an error.
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, _ error) {
			f.panel.ServeHTTP(w, r)
		},
	}
	proxy.ServeHTTP(w, r)
}

// webDecoyHandler is the site telemt shows a WEB request that carries no
// valid credentials. It is the panel's own handler, so such a request gets
// the same decoy page as any other unknown path. telemt replaces the host
// name with this listener's address; the panel's domain is put back, or the
// panel's domain check would answer differently. Only the WEB paths are
// served: this listener is plain HTTP and must not be a second way into the
// panel.
type webDecoyHandler struct {
	panel  http.Handler
	domain string
	route  func(basePath string) (int, bool)
}

func (h webDecoyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()
	base, ok := webProxyBase(path)
	if !ok {
		base = strings.TrimPrefix(path, "/")
	}
	if _, ok := h.route(base); !ok {
		http.NotFound(w, r)
		return
	}
	if h.domain != "" {
		r.Host = h.domain
	}
	h.panel.ServeHTTP(w, r)
}
