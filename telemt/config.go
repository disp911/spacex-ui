// Package telemt runs the bundled telemt MTProto proxy for Telegram and talks
// to its control API. The panel owns the proxy's configuration: users and
// settings live in the panel database and are rendered into telemt's TOML
// config file on every change.
package telemt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Settings are the proxy-wide options the panel exposes.
type Settings struct {
	// Port is the public TCP port Telegram clients connect to.
	Port int
	// TLSDomain is the real website the Fake-TLS handshake imitates. It is
	// encoded into every ee-link, so changing it invalidates issued links.
	TLSDomain string
	// PublicHost overrides the server address put into tg:// links. Empty
	// lets telemt use the external IP it detects at startup.
	PublicHost string
	// APIPort is the loopback port of telemt's control API.
	APIPort int
	// APIToken authorizes the panel against the control API.
	APIToken string
}

// User is one proxy account.
type User struct {
	Name    string
	Secret  string
	Enabled bool
	// ExpiresAt disables the account after this moment; zero means never.
	ExpiresAt time.Time
	// QuotaBytes caps the account's traffic; 0 means unlimited.
	QuotaBytes int64
	// MaxUniqueIPs caps concurrent source IPs; 0 means unlimited.
	MaxUniqueIPs int
}

var (
	userNamePattern  = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,64}$`)
	secretPattern    = regexp.MustCompile(`^[0-9a-f]{32}$`)
	tlsDomainPattern = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?(\.[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?)+$`)
)

// ValidUserName reports whether name is accepted by telemt as a user name.
func ValidUserName(name string) bool {
	return userNamePattern.MatchString(name)
}

// ValidSecret reports whether secret is a 32-character lowercase hex secret.
func ValidSecret(secret string) bool {
	return secretPattern.MatchString(secret)
}

// ValidTLSDomain reports whether domain is a plain host name usable as the
// Fake-TLS domain.
func ValidTLSDomain(domain string) bool {
	return len(domain) <= 253 && tlsDomainPattern.MatchString(domain)
}

// Validate checks the settings before they are written to a config file.
func (s Settings) Validate() error {
	if s.Port < 1 || s.Port > 65535 {
		return fmt.Errorf("invalid proxy port %d", s.Port)
	}
	if !ValidTLSDomain(s.TLSDomain) {
		return fmt.Errorf("invalid TLS domain %q", s.TLSDomain)
	}
	if strings.ContainsAny(s.PublicHost, " /?#&=\"") {
		return fmt.Errorf("invalid public host %q", s.PublicHost)
	}
	if s.APIPort < 1 || s.APIPort > 65535 || s.APIPort == s.Port {
		return fmt.Errorf("invalid API port %d", s.APIPort)
	}
	return nil
}

type fileConfig struct {
	General    fileGeneral    `toml:"general"`
	Server     fileServer     `toml:"server"`
	Censorship fileCensorship `toml:"censorship"`
	Access     fileAccess     `toml:"access"`
}

type fileGeneral struct {
	LogLevel       string    `toml:"log_level"`
	DisableColors  bool      `toml:"disable_colors"`
	QuotaStatePath string    `toml:"quota_state_path"`
	Modes          fileModes `toml:"modes"`
	Links          fileLinks `toml:"links"`
}

type fileModes struct {
	Classic bool `toml:"classic"`
	Secure  bool `toml:"secure"`
	TLS     bool `toml:"tls"`
}

type fileLinks struct {
	PublicHost string `toml:"public_host,omitempty"`
}

type fileServer struct {
	Port int     `toml:"port"`
	API  fileAPI `toml:"api"`
}

type fileAPI struct {
	Enabled    bool     `toml:"enabled"`
	Listen     string   `toml:"listen"`
	Whitelist  []string `toml:"whitelist"`
	AuthHeader string   `toml:"auth_header"`
}

type fileCensorship struct {
	TLSDomain    string `toml:"tls_domain"`
	Mask         bool   `toml:"mask"`
	TLSEmulation bool   `toml:"tls_emulation"`
	TLSFrontDir  string `toml:"tls_front_dir"`
}

type fileAccess struct {
	Users            map[string]string `toml:"users"`
	UserEnabled      map[string]bool   `toml:"user_enabled,omitempty"`
	UserExpirations  map[string]string `toml:"user_expirations,omitempty"`
	UserDataQuota    map[string]int64  `toml:"user_data_quota,omitempty"`
	UserMaxUniqueIPs map[string]int    `toml:"user_max_unique_ips,omitempty"`
}

// BuildConfig renders the telemt TOML config for settings and users. Only
// Fake-TLS (ee) mode is enabled: it is the only mode that looks like ordinary
// HTTPS to a censor. telemt refuses to start without users, so at least one
// is required.
func BuildConfig(s Settings, users []User) ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("telemt needs at least one user")
	}

	access := fileAccess{
		Users:            make(map[string]string, len(users)),
		UserEnabled:      map[string]bool{},
		UserExpirations:  map[string]string{},
		UserDataQuota:    map[string]int64{},
		UserMaxUniqueIPs: map[string]int{},
	}
	for _, u := range users {
		if !ValidUserName(u.Name) {
			return nil, fmt.Errorf("invalid user name %q", u.Name)
		}
		if !ValidSecret(u.Secret) {
			return nil, fmt.Errorf("invalid secret for user %q", u.Name)
		}
		if _, dup := access.Users[u.Name]; dup {
			return nil, fmt.Errorf("duplicate user name %q", u.Name)
		}
		access.Users[u.Name] = u.Secret
		if !u.Enabled {
			access.UserEnabled[u.Name] = false
		}
		if !u.ExpiresAt.IsZero() {
			access.UserExpirations[u.Name] = u.ExpiresAt.UTC().Format(time.RFC3339)
		}
		if u.QuotaBytes > 0 {
			access.UserDataQuota[u.Name] = u.QuotaBytes
		}
		if u.MaxUniqueIPs > 0 {
			access.UserMaxUniqueIPs[u.Name] = u.MaxUniqueIPs
		}
	}

	cfg := fileConfig{
		General: fileGeneral{
			LogLevel:       "normal",
			DisableColors:  true,
			QuotaStatePath: "telemt.limit.json",
			Modes:          fileModes{TLS: true},
			Links:          fileLinks{PublicHost: s.PublicHost},
		},
		Server: fileServer{
			Port: s.Port,
			API: fileAPI{
				Enabled:    true,
				Listen:     fmt.Sprintf("127.0.0.1:%d", s.APIPort),
				Whitelist:  []string{"127.0.0.1/32"},
				AuthHeader: authHeader(s.APIToken),
			},
		},
		Censorship: fileCensorship{
			TLSDomain:    s.TLSDomain,
			Mask:         true,
			TLSEmulation: true,
			TLSFrontDir:  "tlsfront",
		},
		Access: access,
	}
	return toml.Marshal(cfg)
}

// restartKey captures the settings telemt cannot hot-reload: a change in any
// of them needs a process restart, while user changes only need a reload.
func (s Settings) restartKey() string {
	return fmt.Sprintf("%d|%s|%s|%d|%s", s.Port, s.TLSDomain, s.PublicHost, s.APIPort, s.APIToken)
}

// Link builds the Fake-TLS tg://proxy link for secret, the same way telemt
// does, for when the running proxy cannot be asked for it.
func Link(host string, port int, secret, tlsDomain string) string {
	return fmt.Sprintf("tg://proxy?server=%s&port=%d&secret=ee%s%s", host, port, secret, hex.EncodeToString([]byte(tlsDomain)))
}

func authHeader(token string) string {
	if token == "" {
		return ""
	}
	return "Bearer " + token
}
