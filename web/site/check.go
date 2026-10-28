package site

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CheckResult is the outcome of the last look at the site from outside.
type CheckResult struct {
	OK    bool      `json:"ok"`
	At    time.Time `json:"at"`
	Error string    `json:"error"`
}

// Check opens the site as a visitor does: a TLS connection to addr naming
// domain, a domain or an IP address, a certificate a browser trusts for it, and a page at the root.
// Behind REALITY this goes the way of anyone without a client's key, so it
// fails when the inbound no longer passes visitors on as well as when the
// site or its certificate is down. Nil roots are the system's.
func Check(ctx context.Context, domain, addr string, roots *x509.CertPool) CheckResult {
	res := CheckResult{At: time.Now()}
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
				return (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, network, addr)
			},
			TLSClientConfig:   &tls.Config{ServerName: domain, RootCAs: roots},
			ForceAttemptHTTP2: true,
			DisableKeepAlives: true,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	host := domain
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+host+"/", nil)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	resp, err := client.Do(req)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		res.Error = fmt.Sprintf("the site answered %s", resp.Status)
		return res
	}
	res.OK = true
	return res
}

// Checker keeps the last result of Check and reports when it changes.
type Checker struct {
	mu   sync.Mutex
	last CheckResult
	seen bool
}

// Record stores res and reports whether the site went up or down with it;
// the first result counts as a change only if it is a failure.
func (c *Checker) Record(res CheckResult) (changed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	changed = (!c.seen && !res.OK) || (c.seen && c.last.OK != res.OK)
	c.last, c.seen = res, true
	return changed
}

// Last returns the last result, zero before the first check.
func (c *Checker) Last() CheckResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}
