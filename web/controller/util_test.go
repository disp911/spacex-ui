package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRemoteIpBelievesOnlyALocalProxy(t *testing.T) {
	for _, tc := range []struct {
		remote, realIP, forwarded, want string
	}{
		{"203.0.113.7:51000", "", "", "203.0.113.7"},
		{"203.0.113.7:51000", "198.51.100.9", "198.51.100.10", "203.0.113.7"},
		{"127.0.0.1:40000", "198.51.100.9", "", "198.51.100.9"},
		{"[::1]:40000", "", "198.51.100.10, 127.0.0.1", "198.51.100.10"},
		{"127.0.0.1:40000", "", "", "127.0.0.1"},
		{"@", "198.51.100.9", "", "unknown"},
	} {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = tc.remote
		if tc.realIP != "" {
			c.Request.Header.Set("X-Real-IP", tc.realIP)
		}
		if tc.forwarded != "" {
			c.Request.Header.Set("X-Forwarded-For", tc.forwarded)
		}
		if got := getRemoteIp(c); got != tc.want {
			t.Errorf("from %s with X-Real-IP %q, X-Forwarded-For %q: got %s, want %s", tc.remote, tc.realIP, tc.forwarded, got, tc.want)
		}
	}
}
