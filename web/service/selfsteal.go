package service

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/disp911/spacex-ui/v2/config"
	"github.com/disp911/spacex-ui/v2/database"
	"github.com/disp911/spacex-ui/v2/database/model"
	"github.com/disp911/spacex-ui/v2/util/common"
	"github.com/disp911/spacex-ui/v2/util/random"
)

// SiteService switches where the website on port 443 answers. Behind
// REALITY the panel keeps an inbound of its own for it, the site inbound:
// VLESS REALITY on 443 whose target is the website's socket, so whoever is
// not a client gets the website itself ("selfsteal"). The inbound is kept
// switched off while the website runs any other way, so its clients keep
// their keys for when it comes back.
type SiteService struct {
	settingService SettingService
	inboundService InboundService
}

// SwitchMode puts the website in the mode, with the site inbound made,
// fitted or switched off to suit. It changes only the database: the panel
// takes the mode on its next start. The site inbound is returned while the
// mode is SiteReality.
func (s *SiteService) SwitchMode(mode string) (*model.Inbound, error) {
	if mode != SiteOff && mode != SiteDirect && mode != SiteReality {
		return nil, common.NewErrorf("unknown website mode %q", mode)
	}
	domain, err := s.settingService.GetSiteDomain()
	if err != nil {
		return nil, err
	}
	if mode != SiteOff {
		if domain == "" {
			return nil, common.NewError("no domain is set for the website")
		}
		if p, _ := s.settingService.GetBasePath(); p == "/" {
			return nil, common.NewError("the panel needs a secret path to share port 443 with the website")
		}
		if p, _ := s.settingService.GetPort(); p == 443 {
			return nil, common.NewError("the panel's own port is 443, which the website takes: set another one first")
		}
	}

	ib, err := s.SiteInbound()
	if err != nil {
		return nil, err
	}
	switch mode {
	case SiteReality:
		if err := s.free443(ib, false); err != nil {
			return nil, err
		}
		if ib == nil {
			ib, err = s.makeInbound(domain)
		} else {
			err = s.fitInbound(ib, domain)
		}
		if err != nil {
			return nil, err
		}
	case SiteDirect:
		if err := s.free443(ib, true); err != nil {
			return nil, err
		}
		fallthrough
	default:
		if ib != nil && ib.Enable && ib.Port == 443 {
			if err := database.GetDB().Model(model.Inbound{}).Where("id = ?", ib.Id).Update("enable", false).Error; err != nil {
				return nil, err
			}
		}
		ib = nil
	}
	if err := s.settingService.SetSiteMode(mode); err != nil {
		return nil, err
	}
	return ib, nil
}

// SiteInbound is the inbound the panel made for the website, or nil when
// there is none. One made into an inbound of another kind while the website
// ran without it is the administrator's now and no longer counts.
func (s *SiteService) SiteInbound() (*model.Inbound, error) {
	id, err := s.settingService.getInt("siteInbound")
	if err != nil || id == 0 {
		return nil, err
	}
	var list []*model.Inbound
	if err := database.GetDB().Model(model.Inbound{}).Where("id = ?", id).Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 || !siteTarget(list[0]) {
		return nil, nil
	}
	return list[0], nil
}

// free443 makes sure no other inbound holds port 443: a switched-on one
// would keep the website from it, and with REALITY in front a switched-off
// one could not be edited or switched on without breaking it.
func (s *SiteService) free443(own *model.Inbound, enabledOnly bool) error {
	q := database.GetDB().Model(model.Inbound{}).Where("port = ?", 443)
	if own != nil {
		q = q.Where("id != ?", own.Id)
	}
	if enabledOnly {
		q = q.Where("enable = ?", true)
	}
	var others []*model.Inbound
	if err := q.Find(&others).Error; err != nil {
		return err
	}
	if len(others) == 0 {
		return nil
	}
	names := make([]string, 0, len(others))
	for _, o := range others {
		names = append(names, fmt.Sprintf("%q (id %d)", o.Remark, o.Id))
	}
	return common.NewErrorf("port 443 is taken by inbound %s: move it to another port first", strings.Join(names, ", "))
}

func (s *SiteService) makeInbound(domain string) (*model.Inbound, error) {
	db := database.GetDB()
	var userId int
	if err := db.Model(model.User{}).Select("id").Order("id").Limit(1).Scan(&userId).Error; err != nil {
		return nil, err
	}
	ib, client, err := newSiteInbound(domain, config.GetSiteSocket())
	if err != nil {
		return nil, err
	}
	ib.UserId = userId
	for {
		taken, err := s.inboundService.checkEmailsExistForClients([]model.Client{*client})
		if err != nil {
			return nil, err
		}
		if taken == "" {
			break
		}
		client.Email = lowerNum(8)
		if err := setClients(ib, client); err != nil {
			return nil, err
		}
	}

	tx := db.Begin()
	if err := tx.Save(ib).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := s.inboundService.AddClientStat(tx, ib.Id, client); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return ib, s.settingService.setInt("siteInbound", ib.Id)
}

// fitInbound switches the site inbound back on, on port 443 and for the
// domain set now.
func (s *SiteService) fitInbound(ib *model.Inbound, domain string) error {
	var stream map[string]any
	if err := json.Unmarshal([]byte(ib.StreamSettings), &stream); err != nil {
		return err
	}
	reality, _ := stream["realitySettings"].(map[string]any)
	names, _ := reality["serverNames"].([]any)
	if !slices.Contains(names, any(domain)) {
		old := ""
		if len(names) > 0 {
			old, _ = names[0].(string)
		}
		reality["serverNames"] = []string{domain}
		proxies, _ := stream["externalProxy"].([]any)
		for _, p := range proxies {
			if p, ok := p.(map[string]any); ok && p["dest"] == old {
				p["dest"] = domain
			}
		}
		out, err := json.MarshalIndent(stream, "", "  ")
		if err != nil {
			return err
		}
		ib.StreamSettings = string(out)
	}
	ib.Enable = true
	ib.Port = 443
	ib.Tag = "inbound-443"
	if ib.Listen != "" && ib.Listen != "0.0.0.0" && ib.Listen != "::" && ib.Listen != "::0" {
		ib.Tag = fmt.Sprintf("inbound-%v:443", ib.Listen)
	}
	return database.GetDB().Save(ib).Error
}

// newSiteInbound makes the site inbound with its first client: VLESS with
// Vision over REALITY on port 443 with fresh keys, which sends whoever is
// not a client to the website's socket.
func newSiteInbound(domain, socket string) (*model.Inbound, *model.Client, error) {
	key, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().UnixMilli()
	client := &model.Client{
		ID:        uuid.NewString(),
		Flow:      "xtls-rprx-vision",
		Email:     lowerNum(8),
		Enable:    true,
		SubID:     lowerNum(16),
		CreatedAt: now,
		UpdatedAt: now,
	}
	stream := map[string]any{
		"network":  "tcp",
		"security": "reality",
		// Links name the domain, wherever the panel is opened from.
		"externalProxy": []any{map[string]any{"forceTls": "same", "dest": domain, "port": 443, "remark": ""}},
		"realitySettings": map[string]any{
			"show":         false,
			"xver":         0,
			"target":       socket,
			"serverNames":  []string{domain},
			"privateKey":   base64.RawURLEncoding.EncodeToString(key.Bytes()),
			"minClientVer": "",
			"maxClientVer": "",
			"maxTimediff":  0,
			"shortIds":     shortIds(),
			"mldsa65Seed":  "",
			"settings": map[string]any{
				"publicKey":     base64.RawURLEncoding.EncodeToString(key.PublicKey().Bytes()),
				"fingerprint":   "chrome",
				"serverName":    "",
				"spiderX":       "/",
				"mldsa65Verify": "",
			},
			// Whoever is passed on to the website gets it at full speed for
			// as much as the site has to give, then slowly: the server
			// cannot be measured or used through it with large transfers.
			"limitFallbackUpload":   map[string]any{"afterBytes": 1 << 20, "bytesPerSec": 128 << 10, "burstBytesPerSec": 256 << 10},
			"limitFallbackDownload": map[string]any{"afterBytes": 16 << 20, "bytesPerSec": 1 << 20, "burstBytesPerSec": 2 << 20},
		},
		"tcpSettings": map[string]any{"acceptProxyProtocol": false, "header": map[string]any{"type": "none"}},
	}
	streamJSON, err := json.MarshalIndent(stream, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	sniffing, err := json.MarshalIndent(map[string]any{
		"enabled": false, "destOverride": []string{"http", "tls", "quic", "fakedns"}, "metadataOnly": false, "routeOnly": false,
	}, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	ib := &model.Inbound{
		Remark:         domain,
		Enable:         true,
		TrafficReset:   "never",
		Port:           443,
		Protocol:       model.VLESS,
		StreamSettings: string(streamJSON),
		Tag:            "inbound-443",
		Sniffing:       string(sniffing),
	}
	if err := setClients(ib, client); err != nil {
		return nil, nil, err
	}
	return ib, client, nil
}

func setClients(ib *model.Inbound, client *model.Client) error {
	settings, err := json.MarshalIndent(map[string]any{
		"clients":    []*model.Client{client},
		"decryption": "none",
		"encryption": "none",
		"fallbacks":  []any{},
		"testseed":   []int{900, 500, 900, 256},
	}, "", "  ")
	if err != nil {
		return err
	}
	ib.Settings = string(settings)
	return nil
}

// shortIds are of every even length from 2 to 16 hex digits in random
// order, as the inbound window makes them.
func shortIds() []string {
	ids := make([]string, 0, 8)
	for n := 2; n <= 16; n += 2 {
		b := make([]byte, n/2)
		rand.Read(b)
		ids = append(ids, hex.EncodeToString(b))
	}
	for i := len(ids) - 1; i > 0; i-- {
		j := random.Num(i + 1)
		ids[i], ids[j] = ids[j], ids[i]
	}
	return ids
}

func lowerNum(n int) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[random.Num(len(chars))]
	}
	return string(b)
}

// FirstClientLinks are what the site inbound's first client connects with:
// its VLESS link and, while subscriptions are on, its subscription link,
// both through the website's domain on port 443.
func (s *SiteService) FirstClientLinks(ib *model.Inbound) (link, sub string, err error) {
	clients, err := s.inboundService.GetClients(ib)
	if err != nil || len(clients) == 0 {
		return "", "", err
	}
	c := clients[0]
	var stream struct {
		Reality struct {
			ServerNames []string `json:"serverNames"`
			ShortIds    []string `json:"shortIds"`
			Settings    struct {
				PublicKey   string `json:"publicKey"`
				Fingerprint string `json:"fingerprint"`
				SpiderX     string `json:"spiderX"`
			} `json:"settings"`
		} `json:"realitySettings"`
	}
	if err := json.Unmarshal([]byte(ib.StreamSettings), &stream); err != nil {
		return "", "", err
	}
	r := stream.Reality
	domain, err := s.settingService.GetSiteDomain()
	if err != nil {
		return "", "", err
	}
	q := url.Values{}
	q.Set("type", "tcp")
	q.Set("encryption", "none")
	q.Set("security", "reality")
	q.Set("pbk", r.Settings.PublicKey)
	q.Set("fp", r.Settings.Fingerprint)
	if len(r.ServerNames) > 0 {
		q.Set("sni", r.ServerNames[0])
	}
	if len(r.ShortIds) > 0 {
		q.Set("sid", r.ShortIds[0])
	}
	q.Set("spx", r.Settings.SpiderX)
	if c.Flow != "" {
		q.Set("flow", c.Flow)
	}
	link = (&url.URL{
		Scheme:   "vless",
		User:     url.User(c.ID),
		Host:     net.JoinHostPort(domain, "443"),
		RawQuery: q.Encode(),
		Fragment: ib.Remark + "-" + c.Email,
	}).String()

	if on, _ := s.settingService.GetSubEnable(); on && c.SubID != "" {
		sub, _ = s.settingService.GetSubURI()
		if sub == "" {
			path, _ := s.settingService.GetSubPath()
			sub = "https://" + domain + path
		}
		if !strings.HasSuffix(sub, "/") {
			sub += "/"
		}
		sub += c.SubID
	}
	return link, sub, nil
}

// siteTarget reports whether an inbound sends whoever is not a client to the
// website's socket: VLESS over REALITY with the socket as its target.
func siteTarget(ib *model.Inbound) bool {
	if ib.Protocol != model.VLESS {
		return false
	}
	var stream struct {
		Security string `json:"security"`
		Reality  struct {
			Target string `json:"target"`
			Dest   string `json:"dest"`
		} `json:"realitySettings"`
	}
	if json.Unmarshal([]byte(ib.StreamSettings), &stream) != nil || stream.Security != "reality" {
		return false
	}
	target := stream.Reality.Target
	if target == "" {
		target = stream.Reality.Dest
	}
	return target == config.GetSiteSocket()
}

// guardSite refuses a change that would take port 443 from the website: no
// other inbound may hold it while the website runs, and behind REALITY the
// site inbound carries the website and the panel itself, so it stays on,
// on 443, and without limits that would one day switch it off. old is nil
// for a new inbound and next is nil for a deleted one.
func (s *InboundService) guardSite(old, next *model.Inbound) error {
	settings := SettingService{}
	mode, err := settings.GetSiteMode()
	if err != nil || mode == SiteOff {
		return err
	}
	siteId := 0
	if mode == SiteReality {
		siteId, _ = settings.getInt("siteInbound")
	}
	if old != nil && siteId != 0 && old.Id == siteId {
		const how = "Switch the website to another mode first: x-ui site -mode direct"
		switch {
		case next == nil:
			return common.NewError("This inbound carries the website and the panel on port 443 and cannot be deleted. " + how)
		case !next.Enable:
			return common.NewError("This inbound carries the website and the panel on port 443 and cannot be switched off. " + how)
		case next.Port != 443 || next.Listen != old.Listen || !siteTarget(next):
			return common.NewError("This inbound carries the website and the panel on port 443: its port, address, protocol, security and target stay as they are. " + how)
		case next.Total != 0 || next.ExpiryTime != 0:
			return common.NewError("This inbound carries the website and the panel on port 443 and takes no traffic limit or expiry date: they would switch it off.")
		}
		return nil
	}
	if next != nil && next.Enable && next.Port == 443 {
		return common.NewError("Port 443 belongs to the panel's website. Switch the website off first: x-ui site -mode off")
	}
	return nil
}
