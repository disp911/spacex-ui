package service

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xtls/xray-core/infra/conf"

	"github.com/disp911/spacex-ui/v2/config"
	"github.com/disp911/spacex-ui/v2/database"
	"github.com/disp911/spacex-ui/v2/database/model"
)

func TestSiteInboundIsWhatXrayTakes(t *testing.T) {
	ib, client, err := newSiteInbound("docs.example.com", "/run/x-ui/site.sock")
	if err != nil {
		t.Fatal(err)
	}
	var stream conf.StreamConfig
	if err := json.Unmarshal([]byte(ib.StreamSettings), &stream); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Build(); err != nil {
		t.Fatalf("Xray refuses the stream settings: %v", err)
	}
	r := stream.REALITYSettings
	if r.Type != "unix" || r.ServerNames[0] != "docs.example.com" || len(r.ShortIds) != 8 {
		t.Fatalf("target type %q, names %v, short ids %v", r.Type, r.ServerNames, r.ShortIds)
	}
	if r.LimitFallbackDownload.AfterBytes == 0 || r.LimitFallbackUpload.BytesPerSec == 0 {
		t.Fatalf("fallback limits are not read: %+v %+v", r.LimitFallbackUpload, r.LimitFallbackDownload)
	}

	var settings conf.VLessInboundConfig
	if err := json.Unmarshal([]byte(ib.Settings), &settings); err != nil {
		t.Fatal(err)
	}
	if _, err := settings.Build(); err != nil {
		t.Fatalf("Xray refuses the settings: %v", err)
	}
	if client.Flow != "xtls-rprx-vision" || len(client.Email) != 8 || len(client.SubID) != 16 {
		t.Fatalf("client %+v", client)
	}
	if ib.Port != 443 || ib.Tag != "inbound-443" || !ib.Enable || ib.Protocol != model.VLESS {
		t.Fatalf("inbound %+v", ib)
	}
}

func TestSiteTarget(t *testing.T) {
	sock := config.GetSiteSocket()
	ib, _, err := newSiteInbound("docs.example.com", sock)
	if err != nil {
		t.Fatal(err)
	}
	if !siteTarget(ib) {
		t.Fatal("the site inbound must be told apart by its target")
	}
	foreign, _, err := newSiteInbound("docs.example.com", "www.example.org:443")
	if err != nil {
		t.Fatal(err)
	}
	if siteTarget(foreign) {
		t.Fatal("a REALITY inbound with another target is not the site's")
	}
	other := *ib
	other.Protocol = model.Trojan
	if siteTarget(&other) {
		t.Fatal("only VLESS carries the website")
	}
}

func siteTestDB(t *testing.T) {
	t.Helper()
	if err := database.InitDB(filepath.Join(t.TempDir(), "x-ui.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.CloseDB() })
	s := SettingService{}
	for k, v := range map[string]string{"siteDomain": "docs.example.com", "webBasePath": "/secret/", "webPort": "2053", "subEnable": "true"} {
		if err := s.saveSetting(k, v); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSwitchModeKeepsTheSiteInbound(t *testing.T) {
	siteTestDB(t)
	sites := SiteService{}
	inbounds := InboundService{}

	ib, err := sites.SwitchMode(SiteReality)
	if err != nil {
		t.Fatal(err)
	}
	if ib == nil || !ib.Enable || ib.Port != 443 || !siteTarget(ib) {
		t.Fatalf("site inbound %+v", ib)
	}
	link, sub, err := sites.FirstClientLinks(ib)
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(link)
	if err != nil || u.Scheme != "vless" || u.Host != "docs.example.com:443" || u.Query().Get("sni") != "docs.example.com" ||
		u.Query().Get("flow") != "xtls-rprx-vision" || u.Query().Get("security") != "reality" || u.Query().Get("pbk") == "" {
		t.Fatalf("link %s", link)
	}
	if !strings.HasPrefix(sub, "https://docs.example.com/sub/") {
		t.Fatalf("subscription %s", sub)
	}

	// The site inbound stays on, on 443 and without limits.
	if _, err := inbounds.DelInbound(ib.Id); err == nil {
		t.Fatal("the site inbound must not be deleted behind REALITY")
	}
	if _, err := inbounds.SetInboundEnable(ib.Id, false); err == nil {
		t.Fatal("the site inbound must not be switched off behind REALITY")
	}
	moved := *ib
	moved.Port = 8443
	if err := inbounds.guardSite(ib, &moved); err == nil {
		t.Fatal("the site inbound must stay on 443")
	}
	limited := *ib
	limited.ExpiryTime = 1
	if err := inbounds.guardSite(ib, &limited); err == nil {
		t.Fatal("the site inbound must not expire")
	}
	renamed := *ib
	renamed.Remark = "mine"
	if err := inbounds.guardSite(ib, &renamed); err != nil {
		t.Fatalf("a new remark is fine: %v", err)
	}
	if err := inbounds.guardSite(nil, &model.Inbound{Port: 443, Enable: true}); err == nil {
		t.Fatal("no other inbound may take 443")
	}

	// Without REALITY the inbound is kept, switched off, and 443 stays the website's.
	if got, err := sites.SwitchMode(SiteDirect); err != nil || got != nil {
		t.Fatalf("direct: %v %v", got, err)
	}
	kept, err := sites.SiteInbound()
	if err != nil || kept == nil || kept.Id != ib.Id || kept.Enable {
		t.Fatalf("kept %+v %v", kept, err)
	}
	if _, err := inbounds.SetInboundEnable(ib.Id, true); err == nil {
		t.Fatal("an inbound on 443 must not come on while the panel holds it")
	}
	if _, err := inbounds.DelInbound(ib.Id); err != nil {
		t.Fatalf("switched off, the inbound may go: %v", err)
	}

	// Back behind REALITY a new inbound is made.
	again, err := sites.SwitchMode(SiteReality)
	if err != nil || again == nil || again.Id == ib.Id {
		t.Fatalf("again %+v %v", again, err)
	}

	// Off: nothing guards 443 any more.
	if _, err := sites.SwitchMode(SiteOff); err != nil {
		t.Fatal(err)
	}
	if err := inbounds.guardSite(nil, &model.Inbound{Port: 443, Enable: true}); err != nil {
		t.Fatalf("without the website 443 is free: %v", err)
	}
}

func TestSwitchModeRefusesATakenPort(t *testing.T) {
	siteTestDB(t)
	other := &model.Inbound{Tag: "inbound-443", Port: 443, Protocol: model.Trojan, Remark: "old", Enable: true}
	if err := database.GetDB().Create(other).Error; err != nil {
		t.Fatal(err)
	}
	sites := SiteService{}
	if _, err := sites.SwitchMode(SiteReality); err == nil || !strings.Contains(err.Error(), `"old"`) {
		t.Fatalf("REALITY must not start beside another inbound on 443: %v", err)
	}
	if _, err := sites.SwitchMode(SiteDirect); err == nil {
		t.Fatal("the panel must not take 443 from an inbound that is on")
	}
	if mode, _ := (&SettingService{}).GetSiteMode(); mode != SiteOff {
		t.Fatalf("a refused switch must leave the mode, got %q", mode)
	}
}
