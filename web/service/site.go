package service

import (
	"sync"

	"github.com/disp911/spacex-ui/v2/web/site"
)

// siteStatus is what the running website reports of itself: its
// certificate and the last look at it from outside.
var siteStatus struct {
	sync.Mutex
	certs   *site.Certs
	checker *site.Checker
}

// SetSiteStatus records the running website's certificates and checks; nil
// when the website stops.
func SetSiteStatus(certs *site.Certs, checker *site.Checker) {
	siteStatus.Lock()
	siteStatus.certs, siteStatus.checker = certs, checker
	siteStatus.Unlock()
}

// GetSiteStatus describes the website's certificate and the last check of
// it; running is false while there is no website.
func GetSiteStatus() (cert site.CertStatus, check site.CheckResult, running bool) {
	siteStatus.Lock()
	certs, checker := siteStatus.certs, siteStatus.checker
	siteStatus.Unlock()
	if certs == nil {
		return cert, check, false
	}
	return certs.Status(), checker.Last(), true
}
