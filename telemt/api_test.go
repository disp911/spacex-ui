package telemt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIClientParsesEnvelopesAndSendsToken(t *testing.T) {
	var gotAuth, gotResetPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/users":
			w.Write([]byte(`{"ok":true,"data":[{"username":"alice","enabled":true,"current_connections":2,
				"active_unique_ips":1,"active_unique_ips_list":["5.6.7.8"],"total_octets":12345,
				"links":{"classic":[],"secure":[],"tls":["tg://proxy?server=1.2.3.4&port=8443&secret=eeabc"],"tls_domains":[]}}],"revision":"x"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/system/info":
			w.Write([]byte(`{"ok":true,"data":{"version":"3.5.7","uptime_seconds":12.5}}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/reset-quota"):
			gotResetPath = r.URL.EscapedPath()
			w.Write([]byte(`{"ok":true,"data":{"username":"a b","used_bytes":0}}`))
		default:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"ok":false,"error":{"code":"unauthorized","message":"bad token"},"request_id":1}`))
		}
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "secret-token")
	ctx := context.Background()

	users, err := c.Users(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer secret-token" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if len(users) != 1 || users[0].Username != "alice" || users[0].CurrentConnections != 2 ||
		users[0].TotalOctets != 12345 || len(users[0].ActiveIPs) != 1 || len(users[0].Links.TLS) != 1 {
		t.Fatalf("unexpected users %+v", users)
	}

	info, err := c.SystemInfo(ctx)
	if err != nil || info.Version != "3.5.7" {
		t.Fatalf("SystemInfo = %+v, %v", info, err)
	}

	if err := c.ResetQuota(ctx, "a b"); err != nil {
		t.Fatal(err)
	}
	if gotResetPath != "/v1/users/a%20b/reset-quota" {
		t.Fatalf("reset path = %q", gotResetPath)
	}

	err = c.do(ctx, http.MethodGet, "/v1/unknown", nil)
	if err == nil || !strings.Contains(err.Error(), "unauthorized") || !strings.Contains(err.Error(), "bad token") {
		t.Fatalf("API error should carry code and message, got %v", err)
	}
}
