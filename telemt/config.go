// Package telemt runs the bundled telemt MTProto proxy for Telegram. Every
// panel inbound with the mtproto protocol gets its own telemt process; the
// panel renders the inbound into telemt's TOML config and reads per-client
// traffic back from telemt's loopback metrics endpoint.
package telemt

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Settings are the per-inbound options.
type Settings struct {
	// Listen is the IP to bind; empty binds all interfaces.
	Listen string
	// Port is the public TCP port Telegram clients connect to.
	Port int
	// TLSDomain is the real website the Fake-TLS handshake imitates. It is
	// encoded into every ee-link, so changing it invalidates issued links.
	TLSDomain string
	// MetricsPort is the loopback port of telemt's metrics endpoint.
	MetricsPort int
}

// User is one proxy account.
type User struct {
	// Name is the telemt user name; see UserName.
	Name   string
	Secret string
	// MaxUniqueIPs caps concurrent source IPs; 0 means unlimited.
	MaxUniqueIPs int
}

var (
	userNamePattern  = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,64}$`)
	secretPattern    = regexp.MustCompile(`^[0-9a-f]{32}$`)
	tlsDomainPattern = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?(\.[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?)+$`)
)

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
	if s.Listen != "" && net.ParseIP(s.Listen) == nil {
		return fmt.Errorf("invalid listen address %q", s.Listen)
	}
	if s.MetricsPort < 1 || s.MetricsPort > 65535 || s.MetricsPort == s.Port {
		return fmt.Errorf("invalid metrics port %d", s.MetricsPort)
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
}

type fileModes struct {
	Classic bool `toml:"classic"`
	Secure  bool `toml:"secure"`
	TLS     bool `toml:"tls"`
}

type fileServer struct {
	Port             int            `toml:"port"`
	MetricsListen    string         `toml:"metrics_listen"`
	MetricsWhitelist []string       `toml:"metrics_whitelist"`
	API              fileAPI        `toml:"api"`
	Listeners        []fileListener `toml:"listeners,omitempty"`
}

type fileAPI struct {
	Enabled bool `toml:"enabled"`
}

type fileListener struct {
	IP string `toml:"ip"`
}

type fileCensorship struct {
	TLSDomain    string `toml:"tls_domain"`
	Mask         bool   `toml:"mask"`
	TLSEmulation bool   `toml:"tls_emulation"`
	TLSFrontDir  string `toml:"tls_front_dir"`
}

type fileAccess struct {
	Users            map[string]string `toml:"users"`
	UserMaxUniqueIPs map[string]int    `toml:"user_max_unique_ips,omitempty"`
}

// BuildConfig renders the telemt TOML config for one inbound. Only Fake-TLS
// (ee) mode is enabled: it is the only mode that looks like ordinary HTTPS to
// a censor. The control API stays off; the panel only needs the metrics
// endpoint, bound to loopback. telemt refuses to start without users, so at
// least one is required.
func BuildConfig(s Settings, users []User) ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("telemt needs at least one user")
	}

	access := fileAccess{
		Users:            make(map[string]string, len(users)),
		UserMaxUniqueIPs: map[string]int{},
	}
	for _, u := range users {
		if !userNamePattern.MatchString(u.Name) {
			return nil, fmt.Errorf("invalid user name %q", u.Name)
		}
		if !ValidSecret(u.Secret) {
			return nil, fmt.Errorf("invalid secret for user %q", u.Name)
		}
		if _, dup := access.Users[u.Name]; dup {
			return nil, fmt.Errorf("duplicate user name %q", u.Name)
		}
		access.Users[u.Name] = u.Secret
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
		},
		Server: fileServer{
			Port:             s.Port,
			MetricsListen:    fmt.Sprintf("127.0.0.1:%d", s.MetricsPort),
			MetricsWhitelist: []string{"127.0.0.1/32"},
		},
		Censorship: fileCensorship{
			TLSDomain:    s.TLSDomain,
			Mask:         true,
			TLSEmulation: true,
			TLSFrontDir:  "tlsfront",
		},
		Access: access,
	}
	if s.Listen != "" {
		cfg.Server.Listeners = []fileListener{{IP: s.Listen}}
	}
	return toml.Marshal(cfg)
}

// restartKey captures the settings telemt cannot hot-reload: a change in any
// of them needs a process restart, while user changes only need a reload.
func (s Settings) restartKey() string {
	return fmt.Sprintf("%s|%d|%s|%d", s.Listen, s.Port, s.TLSDomain, s.MetricsPort)
}

// UserName maps a panel client email to a telemt user name. telemt only
// accepts [A-Za-z0-9_.-]{1,64}; other emails get a stable name derived from
// the email, so traffic can still be attributed back to it.
func UserName(email string) string {
	if userNamePattern.MatchString(email) {
		return email
	}
	var b strings.Builder
	for _, r := range email {
		if r < 128 && (r == '_' || r == '.' || r == '-' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	name := b.String()
	sum := fnv32(email)
	if len(name) > 55 {
		name = name[:55]
	}
	return fmt.Sprintf("%s-%08x", name, sum)
}

func fnv32(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}
