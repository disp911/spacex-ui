package service

import (
	"strings"
	"testing"
)

func TestApplySiteChecksWhatItIsGiven(t *testing.T) {
	siteTestDB(t)
	sites := SiteService{}
	for _, set := range []SiteSettings{
		{Domain: "not a domain"},
		{Domain: "docs.example.com", CertFile: "/nowhere/cert.pem", KeyFile: "/nowhere/key.pem"},
	} {
		if err := sites.ApplySite(set); err == nil {
			t.Errorf("%+v was taken", set)
		}
	}
}

func TestApplySiteKeepsTheModeAndNothingRefused(t *testing.T) {
	siteTestDB(t)
	sites := SiteService{}
	settings := SettingService{}

	// On its own port any name or address will do, and the mode stays.
	if err := sites.ApplySite(SiteSettings{Domain: "203.0.113.7"}); err != nil {
		t.Fatal(err)
	}
	if mode, _ := settings.GetSiteMode(); mode != SitePort {
		t.Fatalf("mode %q", mode)
	}
	if err := sites.ApplySite(SiteSettings{Domain: " Docs.Example.com "}); err != nil {
		t.Fatal(err)
	}

	// Behind REALITY the site inbound takes a new domain; an address is
	// refused and the domain stays as it was.
	if _, err := sites.SwitchMode(SiteReality); err != nil {
		t.Fatal(err)
	}
	if err := sites.ApplySite(SiteSettings{Domain: "vpn.example.com"}); err != nil {
		t.Fatal(err)
	}
	ib, _ := sites.SiteInbound()
	if ib == nil || !ib.Enable || !strings.Contains(ib.StreamSettings, `"vpn.example.com"`) {
		t.Fatalf("site inbound %+v", ib)
	}
	err := sites.ApplySite(SiteSettings{Domain: "203.0.113.7"})
	if err == nil || !strings.Contains(err.Error(), "domain name") {
		t.Fatalf("err %v", err)
	}
	set, err := sites.siteSettings()
	if err != nil {
		t.Fatal(err)
	}
	if mode, _ := settings.GetSiteMode(); set.Domain != "vpn.example.com" || mode != SiteReality {
		t.Fatalf("kept %+v in %q", set, mode)
	}
}

func TestSiteWithoutAWebsite(t *testing.T) {
	siteTestDB(t)
	info, err := (&SiteService{}).Site("panel.example.com:2053")
	if err != nil {
		t.Fatal(err)
	}
	if info.URL != "" || info.Cert != nil || info.Check != nil || info.Socket != "" || info.Port != 2053 || info.Domain != "docs.example.com" {
		t.Fatalf("info %+v", info)
	}
}
