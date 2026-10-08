package web

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebProxyBase(t *testing.T) {
	for path, want := range map[string]string{
		"/c0ffee/":          "c0ffee",
		"/c0ffee/api/v1/up": "c0ffee",
		"/c0ffee":           "",
		"/":                 "",
		"":                  "",
		"//x":               "",
	} {
		got, ok := webProxyBase(path)
		if got != want || ok != (want != "") {
			t.Errorf("webProxyBase(%q) = %q, %v; want %q", path, got, ok, want)
		}
	}
}

// panelStub stands in for the panel and says that it answered.
var panelStub = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "panel %s %s", r.Host, r.URL.Path)
})

func routeTo(port int) func(string) (int, bool) {
	return func(base string) (int, bool) { return port, base == "c0ffee" }
}

func backendPort(t *testing.T, srv *httptest.Server) int {
	t.Helper()
	return srv.Listener.Addr().(*net.TCPAddr).Port
}

func TestWebProxyFrontSendsTheWebSubtreeToTelemt(t *testing.T) {
	var seen *http.Request
	var body string
	telemt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		w.Header().Set("X-Carrier-Mode", "https")
		fmt.Fprint(w, "from telemt")
	}))
	defer telemt.Close()
	front := webProxyFront{panel: panelStub, route: routeTo(backendPort(t, telemt))}

	req := httptest.NewRequest(http.MethodPost, "https://proxy.example.com/c0ffee/api/v1/up?x=1", strings.NewReader("frame"))
	req.RemoteAddr = "198.51.100.7:5555"
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	front.ServeHTTP(rec, req)

	if rec.Body.String() != "from telemt" || rec.Header().Get("X-Carrier-Mode") != "https" {
		t.Fatalf("response = %q %v", rec.Body.String(), rec.Header())
	}
	if seen == nil {
		t.Fatal("telemt got nothing")
	}
	if seen.Host != "proxy.example.com" {
		t.Errorf("telemt must see the public host, got %q", seen.Host)
	}
	if got := seen.Header.Values("X-Forwarded-For"); len(got) != 1 || got[0] != "198.51.100.7" {
		t.Errorf("X-Forwarded-For = %v, want only the client address", got)
	}
	if seen.URL.Path != "/c0ffee/api/v1/up" || seen.URL.RawQuery != "x=1" || body != "frame" {
		t.Errorf("telemt got %s?%s %q", seen.URL.Path, seen.URL.RawQuery, body)
	}
	if rec.Header().Get("Content-Encoding") != "" {
		t.Error("WEB responses must not be compressed")
	}
}

func TestWebProxyFrontLeavesEverythingElseToThePanel(t *testing.T) {
	telemt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("telemt must not get %s", r.URL.Path)
	}))
	defer telemt.Close()
	front := webProxyFront{panel: panelStub, route: routeTo(backendPort(t, telemt))}
	for _, path := range []string{"/", "/c0ffee", "/other/api/v1/up", "/mPY0LuFc1d78pw4lTm/panel/"} {
		rec := httptest.NewRecorder()
		front.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "https://proxy.example.com"+path, nil))
		if want := "panel proxy.example.com " + path; rec.Body.String() != want {
			t.Errorf("%s: got %q, want %q", path, rec.Body.String(), want)
		}
	}
}

func TestWebProxyFrontShowsThePanelWhenTelemtIsDown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	front := webProxyFront{panel: panelStub, route: routeTo(port)}
	rec := httptest.NewRecorder()
	front.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "https://proxy.example.com/c0ffee/", nil))
	if !strings.HasPrefix(rec.Body.String(), "panel ") {
		t.Fatalf("got %d %q, want the panel's page", rec.Code, rec.Body.String())
	}
}

func TestWebProxyFrontPassesWebSocketUpgrades(t *testing.T) {
	telemt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "websocket" || r.Header.Get("Sec-WebSocket-Protocol") != "tproxy-v1.token" {
			http.Error(w, "not an upgrade", http.StatusBadRequest)
			return
		}
		conn, rw, err := w.(http.Hijacker).Hijack()
		if err != nil {
			return
		}
		defer conn.Close()
		rw.WriteString("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n")
		rw.Flush()
		line, _ := rw.ReadString('\n')
		rw.WriteString("echo " + line)
		rw.Flush()
	}))
	defer telemt.Close()
	front := httptest.NewServer(webProxyFront{panel: panelStub, route: routeTo(backendPort(t, telemt))})
	defer front.Close()

	conn, err := net.Dial("tcp", front.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	fmt.Fprint(conn, "GET /c0ffee/api/v1/ws HTTP/1.1\r\nHost: proxy.example.com\r\nConnection: Upgrade\r\nUpgrade: websocket\r\n"+
		"Sec-WebSocket-Version: 13\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Protocol: tproxy-v1.token\r\n\r\n")
	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("upgrade answered %s", resp.Status)
	}
	fmt.Fprint(conn, "ping\n")
	line, err := r.ReadString('\n')
	if err != nil || line != "echo ping\n" {
		t.Fatalf("after the upgrade got %q, %v", line, err)
	}
}

func TestWebDecoyHandlerServesOnlyWebPathsAsThePanelDomain(t *testing.T) {
	h := webDecoyHandler{panel: panelStub, domain: "proxy.example.com", route: routeTo(1)}
	for path, want := range map[string]string{
		"/c0ffee/":               "panel proxy.example.com /c0ffee/",
		"/c0ffee/api/v1/session": "panel proxy.example.com /c0ffee/api/v1/session",
		"/c0ffee":                "panel proxy.example.com /c0ffee",
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "http://127.0.0.1:1234"+path, nil))
		if rec.Body.String() != want {
			t.Errorf("%s: got %q, want %q", path, rec.Body.String(), want)
		}
	}
	for _, path := range []string{"/", "/mPY0LuFc1d78pw4lTm/panel/", "/other/x"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "http://127.0.0.1:1234"+path, nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s must not reach the panel over the loopback decoy, got %d %q", path, rec.Code, rec.Body.String())
		}
	}
}
