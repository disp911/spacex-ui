package service

import (
	"strings"
	"testing"
)

func TestApplySiteChecksWhatItIsGiven(t *testing.T) {
	sites := SiteService{}
	for _, set := range []SiteSettings{
		{Mode: "direct", Domain: "docs.example.com"},
		{Mode: SitePort, Domain: "not a domain"},
		{Mode: SitePort, Domain: "docs.example.com", CertFile: "/nowhere/cert.pem", KeyFile: "/nowhere/key.pem"},
	} {
		if err := sites.ApplySite(set); err == nil {
			t.Errorf("%+v was taken", set)
		}
	}
}

func TestApplySiteKeepsNothingWhenRefused(t *testing.T) {
	siteTestDB(t)
	sites := SiteService{}
	settings := SettingService{}

	// REALITY wants a domain name: the address is refused, and the domain
	// stays as it was.
	err := sites.ApplySite(SiteSettings{Mode: SiteReality, Domain: "203.0.113.7"})
	if err == nil || !strings.Contains(err.Error(), "domain name") {
		t.Fatalf("err %v", err)
	}
	set, err := sites.siteSettings()
	if err != nil {
		t.Fatal(err)
	}
	if set.Domain != "docs.example.com" || set.Mode != SitePort {
		t.Fatalf("kept %+v", set)
	}

	if err := sites.ApplySite(SiteSettings{Mode: SiteReality, Domain: " Docs.Example.com "}); err != nil {
		t.Fatal(err)
	}
	set, _ = sites.siteSettings()
	if set.Mode != SiteReality || set.Domain != "docs.example.com" {
		t.Fatalf("applied %+v", set)
	}
	if ib, _ := sites.SiteInbound(); ib == nil || !ib.Enable {
		t.Fatalf("site inbound %+v", ib)
	}

	// Back on the panel's own port, by IP: the inbound rests.
	if err := sites.ApplySite(SiteSettings{Mode: SitePort, Domain: "203.0.113.7"}); err != nil {
		t.Fatal(err)
	}
	if mode, _ := settings.GetSiteMode(); mode != SitePort {
		t.Fatalf("mode %q", mode)
	}
	if ib, _ := sites.SiteInbound(); ib == nil || ib.Enable {
		t.Fatalf("site inbound %+v", ib)
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
