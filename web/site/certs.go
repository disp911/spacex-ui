package site

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/mholt/acmez/v3/acme"
	"go.uber.org/zap"
)

// Certs obtains and renews the site's certificate on its own, with no
// script beside the panel: from Let's Encrypt, and from ZeroSSL should that
// fail and an email be known. A renewed certificate is served from the next
// handshake on, with nothing restarted. The authorities check the domain
// either over port 80 (HTTPChallenge) or over port 443 by TLS (TLSConfig);
// behind REALITY the latter reaches the site like any other visitor.
type Certs struct {
	domain string
	cache  *certmagic.Cache
	cfg    *certmagic.Config
	le     *certmagic.ACMEIssuer

	mu     sync.Mutex
	failed string
	at     time.Time
}

// CertStatus is what the panel shows of the certificate.
type CertStatus struct {
	Domain string `json:"domain"`
	// Issuer and NotAfter describe the certificate served, zero while
	// there is none.
	Issuer   string    `json:"issuer"`
	NotAfter time.Time `json:"notAfter"`
	// Error is the last failure to obtain or renew it, cleared by success.
	Error   string    `json:"error"`
	ErrorAt time.Time `json:"errorAt"`
}

// NewCerts manages the certificate of domain, keeping it and the accounts in
// folder. Without email there is no second authority.
func NewCerts(folder, domain, email string) *Certs {
	c := &Certs{domain: domain}
	var cfg *certmagic.Config
	c.cache = certmagic.NewCache(certmagic.CacheOptions{
		GetConfigForCert: func(certmagic.Certificate) (*certmagic.Config, error) { return cfg, nil },
		Logger:           zap.NewNop(),
	})
	cfg = certmagic.New(c.cache, certmagic.Config{
		Storage: &certmagic.FileStorage{Path: folder},
		Logger:  zap.NewNop(),
		OnEvent: c.onEvent,
	})
	c.le = certmagic.NewACMEIssuer(cfg, certmagic.ACMEIssuer{
		CA:     certmagic.LetsEncryptProductionCA,
		TestCA: certmagic.LetsEncryptStagingCA,
		Email:  email,
		Agreed: true,
		Logger: zap.NewNop(),
	})
	cfg.Issuers = []certmagic.Issuer{c.le}
	if email != "" {
		cfg.Issuers = append(cfg.Issuers, &zeroSSL{cfg: cfg, email: email, eabURL: zeroSSLEABURL})
	}
	c.cfg = cfg
	return c
}

// Start loads the certificate from storage, or obtains one in the
// background, and keeps it renewed until Stop.
func (c *Certs) Start(ctx context.Context) error {
	return c.cfg.ManageAsync(ctx, []string{c.domain})
}

// Stop ends the renewals.
func (c *Certs) Stop() {
	c.cache.Stop()
}

// Renew obtains a new certificate now, even if the one served is fresh.
func (c *Certs) Renew(ctx context.Context) error {
	return c.cfg.RenewCertAsync(ctx, c.domain, true)
}

// TLSConfig serves the certificate over HTTP/2 and HTTP/1.1 and answers the
// authorities' TLS check.
func (c *Certs) TLSConfig() *tls.Config {
	t := c.cfg.TLSConfig()
	t.MinVersion = tls.VersionTLS12
	t.NextProtos = append([]string{"h2", "http/1.1"}, t.NextProtos...)
	return t
}

// HTTPChallenge answers the authorities' check over port 80 and passes
// every other request to next.
func (c *Certs) HTTPChallenge(next http.Handler) http.Handler {
	return c.le.HTTPChallengeHandler(next)
}

// Status describes the certificate served and the last failure.
func (c *Certs) Status() CertStatus {
	s := CertStatus{Domain: c.domain}
	for _, cert := range c.cache.AllMatchingCertificates(c.domain) {
		if cert.Leaf != nil && cert.Leaf.NotAfter.After(s.NotAfter) {
			s.NotAfter = cert.Leaf.NotAfter
			s.Issuer = issuerName(cert.Leaf)
		}
	}
	c.mu.Lock()
	s.Error, s.ErrorAt = c.failed, c.at
	c.mu.Unlock()
	return s
}

func issuerName(leaf *x509.Certificate) string {
	if len(leaf.Issuer.Organization) > 0 {
		return leaf.Issuer.Organization[0]
	}
	return leaf.Issuer.CommonName
}

func (c *Certs) onEvent(_ context.Context, event string, data map[string]any) error {
	switch event {
	case "cert_failed":
		c.mu.Lock()
		c.failed, c.at = fmt.Sprint(data["error"]), time.Now()
		c.mu.Unlock()
	case "cert_obtained":
		c.mu.Lock()
		c.failed, c.at = "", time.Time{}
		c.mu.Unlock()
	}
	return nil
}

// zeroSSLEABURL hands out the account binding ZeroSSL asks of an ACME client,
// in exchange for an email.
const zeroSSLEABURL = "https://api.zerossl.com/acme/eab-credentials-email"

// zeroSSL issues from ZeroSSL's ACME endpoint. Its account binding is
// fetched on the first issue, so a panel that never falls back never asks.
type zeroSSL struct {
	cfg    *certmagic.Config
	email  string
	eabURL string

	mu     sync.Mutex
	issuer *certmagic.ACMEIssuer
}

func (z *zeroSSL) IssuerKey() string {
	return (&certmagic.ACMEIssuer{CA: certmagic.ZeroSSLProductionCA}).IssuerKey()
}

func (z *zeroSSL) Issue(ctx context.Context, csr *x509.CertificateRequest) (*certmagic.IssuedCertificate, error) {
	issuer, err := z.acme(ctx)
	if err != nil {
		return nil, err
	}
	return issuer.Issue(ctx, csr)
}

func (z *zeroSSL) acme(ctx context.Context) (*certmagic.ACMEIssuer, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	if z.issuer != nil {
		return z.issuer, nil
	}
	eab, err := fetchEAB(ctx, z.eabURL, z.email)
	if err != nil {
		return nil, err
	}
	z.issuer = certmagic.NewACMEIssuer(z.cfg, certmagic.ACMEIssuer{
		CA:              certmagic.ZeroSSLProductionCA,
		Email:           z.email,
		Agreed:          true,
		ExternalAccount: eab,
		Logger:          zap.NewNop(),
	})
	return z.issuer, nil
}

func fetchEAB(ctx context.Context, endpoint, email string) (*acme.EAB, error) {
	form := url.Values{"email": {email}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ZeroSSL account binding: %w", err)
	}
	defer resp.Body.Close()
	var answer struct {
		Success bool   `json:"success"`
		KeyID   string `json:"eab_kid"`
		MACKey  string `json:"eab_hmac_key"`
		Error   struct {
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		return nil, fmt.Errorf("ZeroSSL account binding: %s: %w", resp.Status, err)
	}
	if !answer.Success || answer.KeyID == "" || answer.MACKey == "" {
		if answer.Error.Type != "" {
			return nil, errors.New("ZeroSSL account binding: " + answer.Error.Type)
		}
		return nil, errors.New("ZeroSSL account binding: " + resp.Status)
	}
	return &acme.EAB{KeyID: answer.KeyID, MACKey: answer.MACKey}, nil
}
