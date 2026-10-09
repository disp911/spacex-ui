package xray

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/xtls/xray-core/infra/conf"
)

// The configs below have the shapes the panel's forms write. They are built
// with the xray-core loader the panel bundles, so a core bump that drops or
// renames one of them fails here instead of keeping Xray from starting.

func buildInbound(t *testing.T, raw string) error {
	t.Helper()
	c := new(conf.InboundDetourConfig)
	if err := json.Unmarshal([]byte(raw), c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	_, err := c.Build()
	return err
}

func buildOutbound(t *testing.T, raw string) error {
	t.Helper()
	c := new(conf.OutboundDetourConfig)
	if err := json.Unmarshal([]byte(raw), c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	_, err := c.Build()
	return err
}

const vlessSettings = `{"clients":[{"id":"8e3a3c8e-0f57-4a57-9a6e-4b3c7f0f5e11","email":"a"}],"decryption":"none"}`

func TestKCPInboundWithMkcpLegacyMasksBuilds(t *testing.T) {
	for name, settings := range map[string]string{
		"original": `{"header":"","value":""}`,
		"aes":      `{"header":"","value":"secret"}`,
		"dns":      `{"header":"dns","value":"www.example.com"}`,
		"wechat":   `{"header":"wechat","value":""}`,
	} {
		raw := `{"protocol":"vless","port":10001,"settings":` + vlessSettings + `,
			"streamSettings":{"network":"kcp","kcpSettings":{"mtu":1350},
				"finalmask":{"udp":[{"type":"mkcp-legacy","settings":` + settings + `}]}}}`
		if err := buildInbound(t, raw); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestRemovedMkcpMaskTypesAreRefused(t *testing.T) {
	raw := `{"protocol":"vless","port":10001,"settings":` + vlessSettings + `,
		"streamSettings":{"network":"kcp","finalmask":{"udp":[{"type":"mkcp-aes128gcm","settings":{"password":"x"}}]}}}`
	if err := buildInbound(t, raw); err == nil {
		t.Fatal("xray-core accepts mkcp-aes128gcm again; the panel offers mkcp-legacy only")
	}
}

func TestXHTTPInboundWithSessionIDKeysBuilds(t *testing.T) {
	raw := `{"protocol":"vless","port":10002,"settings":` + vlessSettings + `,
		"streamSettings":{"network":"xhttp","xhttpSettings":{"path":"/x","mode":"packet-up",
			"sessionIDPlacement":"header","sessionIDKey":"X-Session"}}}`
	if err := buildInbound(t, raw); err != nil {
		t.Fatal(err)
	}
	// The core checks the placement under the new name only: a bad one
	// there is refused, the same under the old name goes unread.
	bad := strings.Replace(raw, `"sessionIDPlacement":"header"`, `"sessionIDPlacement":"bogus"`, 1)
	if err := buildInbound(t, bad); err == nil {
		t.Fatal("xray-core no longer reads sessionIDPlacement")
	}
	old := strings.Replace(raw, `"sessionIDPlacement":"header"`, `"sessionPlacement":"bogus"`, 1)
	if err := buildInbound(t, old); err != nil {
		t.Fatalf("xray-core reads the old sessionPlacement again: %v", err)
	}
}

func TestXicmpOutboundMaskBuilds(t *testing.T) {
	raw := `{"protocol":"vless","settings":{"vnext":[{"address":"127.0.0.1","port":443,
		"users":[{"id":"8e3a3c8e-0f57-4a57-9a6e-4b3c7f0f5e11","encryption":"none"}]}]},
		"streamSettings":{"network":"kcp","finalmask":{"udp":[{"type":"xicmp","settings":{"dgram":true,"ips":["1.2.3.4"]}}]}}}`
	if err := buildOutbound(t, raw); err != nil {
		t.Fatal(err)
	}
}

func TestDNSOutboundRulesBuild(t *testing.T) {
	raw := `{"protocol":"dns","settings":{"rules":[
		{"action":"return","qType":"1,28"},
		{"action":"drop","qType":65},
		{"action":"hijack","domain":["domain:example.com"]},
		{"action":"direct"}]}}`
	if err := buildOutbound(t, raw); err != nil {
		t.Fatal(err)
	}
	raw = `{"protocol":"dns","settings":{"rules":[{"action":"reject"}]}}`
	if err := buildOutbound(t, raw); err == nil {
		t.Fatal(`the DNS rule action "reject" is back; the panel writes "return"`)
	}
}
