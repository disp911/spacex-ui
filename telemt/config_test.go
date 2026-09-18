package telemt

import (
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func testSettings() Settings {
	return Settings{Port: 8443, TLSDomain: "www.example.com", MetricsPort: 19090}
}

func lookup(m map[string]any, path string) any {
	var cur any = m
	for _, key := range strings.Split(path, "/") {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = obj[key]
	}
	return cur
}

func TestBuildConfigRendersInbound(t *testing.T) {
	users := []User{
		{Name: "alice", Secret: "00112233445566778899aabbccddeeff", MaxUniqueIPs: 3},
		{Name: "bob.pc", Secret: "0123456789abcdef0123456789abcdef"},
	}
	s := testSettings()
	s.Listen = "10.0.0.5"
	data, err := BuildConfig(s, users)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := toml.Unmarshal(data, &got); err != nil {
		t.Fatalf("generated config is not valid TOML: %v\n%s", err, data)
	}
	want := map[string]any{
		"general/modes/tls":                    true,
		"general/modes/classic":                false,
		"general/modes/secure":                 false,
		"general/disable_colors":               true,
		"general/quota_state_path":             "telemt.limit.json",
		"general/use_middle_proxy":             false,
		"server/port":                          int64(8443),
		"server/metrics_listen":                "127.0.0.1:19090",
		"server/api/enabled":                   false,
		"censorship/tls_domain":                "www.example.com",
		"censorship/mask":                      true,
		"censorship/tls_emulation":             true,
		"access/users/alice":                   "00112233445566778899aabbccddeeff",
		"access/users/bob.pc":                  "0123456789abcdef0123456789abcdef",
		"access/user_max_unique_ips/alice":     int64(3),
		"access/user_max_unique_ips/bob.pc":    nil,
		"general/links":                        nil,
		"access/user_enabled":                  nil,
		"access/users/does-not-exist-sanity":   nil,
		"server/metrics_whitelist":             []any{"127.0.0.1/32"},
		"server/listeners":                     []any{map[string]any{"ip": "10.0.0.5"}},
		"censorship/tls_front_dir":             "tlsfront",
		"general/log_level":                    "normal",
		"access/user_max_unique_ips/unlimited": nil,
	}
	for path, value := range want {
		if g := lookup(got, path); !equalValue(g, value) {
			t.Errorf("%s = %#v, want %#v", path, g, value)
		}
	}
}

func equalValue(a, b any) bool {
	as, aok := a.([]any)
	bs, bok := b.([]any)
	if aok || bok {
		if len(as) != len(bs) {
			return false
		}
		for i := range as {
			if !equalValue(as[i], bs[i]) {
				return false
			}
		}
		return true
	}
	am, aok := a.(map[string]any)
	bm, bok := b.(map[string]any)
	if aok || bok {
		if len(am) != len(bm) {
			return false
		}
		for k := range am {
			if !equalValue(am[k], bm[k]) {
				return false
			}
		}
		return true
	}
	return a == b
}

func TestBuildConfigWithoutListenBindsEverywhere(t *testing.T) {
	data, err := BuildConfig(testSettings(), []User{{Name: "a", Secret: strings.Repeat("a", 32)}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "listeners") {
		t.Fatalf("empty listen must not write listeners:\n%s", data)
	}
}

func TestBuildConfigRejectsInvalidInput(t *testing.T) {
	ok := User{Name: "a", Secret: strings.Repeat("a", 32)}
	cases := map[string]struct {
		settings Settings
		users    []User
	}{
		"no users":           {testSettings(), nil},
		"bad port":           {Settings{Port: 0, TLSDomain: "example.com", MetricsPort: 9090}, []User{ok}},
		"metrics port clash": {Settings{Port: 9090, TLSDomain: "example.com", MetricsPort: 9090}, []User{ok}},
		"no metrics port":    {Settings{Port: 8443, TLSDomain: "example.com"}, []User{ok}},
		"empty domain":       {Settings{Port: 8443, MetricsPort: 9090}, []User{ok}},
		"domain with path":   {Settings{Port: 8443, TLSDomain: "example.com/x", MetricsPort: 9090}, []User{ok}},
		"domain with quote":  {Settings{Port: 8443, TLSDomain: `exa"mple.com`, MetricsPort: 9090}, []User{ok}},
		"listen not an IP":   {Settings{Listen: "example.com", Port: 8443, TLSDomain: "example.com", MetricsPort: 9090}, []User{ok}},
		"bad user name":      {testSettings(), []User{{Name: "bad name", Secret: ok.Secret}}},
		"quote in user name": {testSettings(), []User{{Name: `a"b`, Secret: ok.Secret}}},
		"short secret":       {testSettings(), []User{{Name: "a", Secret: "abc"}}},
		"uppercase secret":   {testSettings(), []User{{Name: "a", Secret: strings.Repeat("A", 32)}}},
		"duplicate user":     {testSettings(), []User{ok, ok}},
	}
	for name, tc := range cases {
		if _, err := BuildConfig(tc.settings, tc.users); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestUserNameMapsEmailsToValidUniqueNames(t *testing.T) {
	if got := UserName("iphone-my"); got != "iphone-my" {
		t.Fatalf("valid emails must be kept, got %q", got)
	}
	seen := map[string]string{}
	for _, email := range []string{"user@mail.ru", "user_mail.ru", "Телефон Степана", "a b", "a_b", strings.Repeat("x", 80) + "@y", ""} {
		name := UserName(email)
		if !userNamePattern.MatchString(name) {
			t.Errorf("UserName(%q) = %q is not a valid telemt user name", email, name)
		}
		if other, dup := seen[name]; dup {
			t.Errorf("UserName(%q) and UserName(%q) collide on %q", email, other, name)
		}
		seen[name] = email
		if UserName(email) != name {
			t.Errorf("UserName(%q) is not stable", email)
		}
	}
}

func TestRestartKeyIgnoresUsers(t *testing.T) {
	a := testSettings()
	if a.restartKey() != testSettings().restartKey() {
		t.Fatal("identical settings must share a restart key")
	}
	for name, mutate := range map[string]func(*Settings){
		"listen":       func(s *Settings) { s.Listen = "10.0.0.1" },
		"port":         func(s *Settings) { s.Port++ },
		"domain":       func(s *Settings) { s.TLSDomain = "other.example.com" },
		"metrics port": func(s *Settings) { s.MetricsPort++ },
	} {
		c := testSettings()
		mutate(&c)
		if c.restartKey() == a.restartKey() {
			t.Errorf("changing %s must change the restart key", name)
		}
	}
}
