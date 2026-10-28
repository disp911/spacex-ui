package web

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/disp911/spacex-ui/v2/config"
	"github.com/disp911/spacex-ui/v2/logger"
	"github.com/disp911/spacex-ui/v2/web/network"
	"github.com/disp911/spacex-ui/v2/web/service"
	"github.com/disp911/spacex-ui/v2/web/site"
)

// siteFront answers on the site's listeners: the panel's secret path is the
// panel's, the subscription paths are the subscription server's, and every
// other path is the website's. A path is matched as a whole first segment
// or more, so "/sub" does not take "/subscribe".
type siteFront struct {
	basePath string
	panel    http.Handler
	site     http.Handler
}

func (f siteFront) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if underPath(p, f.basePath) {
		f.panel.ServeHTTP(w, r)
		return
	}
	if sub := subFront.Load(); sub != nil {
		for _, prefix := range sub.paths {
			if underPath(p, prefix) {
				sub.handler.ServeHTTP(w, r)
				return
			}
		}
	}
	f.site.ServeHTTP(w, r)
}

// underPath reports whether p is prefix, a folder path ending in "/", or
// lies under it; the folder asked for without its slash counts too.
func underPath(p, prefix string) bool {
	if prefix == "" || prefix == "/" {
		return false
	}
	return strings.HasPrefix(p, prefix) || p == strings.TrimSuffix(prefix, "/")
}

type subRoute struct {
	handler http.Handler
	paths   []string
}

var subFront atomic.Pointer[subRoute]

// SetSubFront hands the website the subscription server's handler and the
// paths it answers, while the site serves the subscriptions on port 443 in
// place of the server's own port; nil takes them back.
func SetSubFront(handler http.Handler, paths []string) {
	if handler == nil {
		subFront.Store(nil)
		return
	}
	subFront.Store(&subRoute{handler: handler, paths: paths})
}

var siteRunning atomic.Bool

// SiteRunning reports whether the website serves port 443, and with it the
// subscriptions. The web server starts first, so the subscription server
// asks this once the website is up or has failed to come up.
func SiteRunning() bool {
	return siteRunning.Load()
}

// siteRedirect sends a plain HTTP visitor to the same address over HTTPS.
func siteRedirect(domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		if host == "" {
			host = domain
		}
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusMovedPermanently)
	})
}

// quietLog drops what net/http would print about every visitor whose TLS
// handshake fails, which a public port sees all day.
var quietLog = log.New(io.Discard, "", 0)

// startSite opens the website on port 443, directly or on the socket a
// REALITY inbound passes visitors to, with the panel at its secret path,
// plus port 80 for the authorities' checks and a redirect to HTTPS. The
// panel's own port stays shut unless asked for behind REALITY. It reports
// whether the panel's own port is still to be opened.
func (s *Server) startSite(mode string, panel http.Handler) (ownPort bool, err error) {
	domain, err := s.settingService.GetSiteDomain()
	if err != nil {
		return false, err
	}
	if domain == "" {
		return false, errSite("no domain is set for the website")
	}
	basePath, err := s.settingService.GetBasePath()
	if err != nil {
		return false, err
	}
	if basePath == "/" {
		return false, errSite("the panel needs a secret path to share port 443 with the website")
	}
	email, _ := s.settingService.GetSiteEmail()
	listen, _ := s.settingService.GetListen()

	certs := site.NewCerts(config.GetCertFolder(), domain, email)
	if err := certs.Start(s.ctx); err != nil {
		return false, err
	}

	var ln net.Listener
	switch mode {
	case service.SiteDirect:
		ln, err = net.Listen("tcp", net.JoinHostPort(listen, "443"))
	case service.SiteReality:
		sock := config.GetSiteSocket()
		os.MkdirAll(filepath.Dir(sock), 0o755)
		os.Remove(sock)
		ln, err = net.Listen("unix", sock)
		if err == nil {
			os.Chmod(sock, 0o600)
		}
	default:
		err = errSite("unknown site mode " + mode)
	}
	if err != nil {
		certs.Stop()
		return false, err
	}

	made := site.Placeholder(decoyPage, buildTime())
	front := webProxyFront{
		panel: siteFront{basePath: basePath, panel: panel, site: site.Files(site.Source(config.GetSiteFolder(), made))},
		route: service.TelemtWebRoute,
	}
	s.siteServer = &http.Server{
		Handler:           front,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       75 * time.Second,
		ErrorLog:          quietLog,
	}
	tlsLn := tls.NewListener(ln, certs.TLSConfig())
	go s.siteServer.Serve(tlsLn)
	logger.Info("Website running on", ln.Addr(), "for", domain)

	if l80, err := net.Listen("tcp", net.JoinHostPort(listen, "80")); err != nil {
		logger.Warning("Website: port 80 is not available, certificates are checked over 443 only:", err)
	} else {
		s.redirectServer = &http.Server{
			Handler:           certs.HTTPChallenge(siteRedirect(domain)),
			ReadHeaderTimeout: 60 * time.Second,
			ErrorLog:          quietLog,
		}
		go s.redirectServer.Serve(l80)
	}

	s.siteCerts = certs
	s.siteTLS = certs.TLSConfig()
	s.siteHandler = front
	s.siteCheck = &site.Checker{}
	service.SetSiteStatus(certs, s.siteCheck)
	siteRunning.Store(true)

	if mode == service.SiteReality {
		ownPort, _ = s.settingService.GetSitePanelPort()
	}
	return ownPort, nil
}

// listenOwnPort opens the panel's own port while the website runs: the same
// answers as on 443, with the website's certificate.
func (s *Server) listenOwnPort() (net.Listener, error) {
	listen, err := s.settingService.GetListen()
	if err != nil {
		return nil, err
	}
	port, err := s.settingService.GetPort()
	if err != nil {
		return nil, err
	}
	ln, err := net.Listen("tcp", net.JoinHostPort(listen, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	return tls.NewListener(network.NewAutoHttpsListener(ln), s.siteTLS), nil
}

// checkSite looks at the website from outside and logs when it goes down or
// comes back.
func (s *Server) checkSite() {
	domain, err := s.settingService.GetSiteDomain()
	if err != nil || domain == "" || s.siteCheck == nil || s.ctx.Err() != nil {
		return
	}
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()
	res := site.Check(ctx, domain, net.JoinHostPort(domain, "443"), nil)
	if !s.siteCheck.Record(res) {
		return
	}
	if res.OK {
		logger.Info("Website check: the site is up again")
	} else {
		logger.Warning("Website check failed:", res.Error)
	}
}

func (s *Server) stopSite() {
	siteRunning.Store(false)
	if s.siteServer != nil {
		s.siteServer.Close()
	}
	if s.redirectServer != nil {
		s.redirectServer.Close()
	}
	if s.siteCerts != nil {
		s.siteCerts.Stop()
		service.SetSiteStatus(nil, nil)
	}
}

// buildTime is when the panel's binary was written, the date the placeholder
// site claims for its page.
func buildTime() time.Time {
	if exe, err := os.Executable(); err == nil {
		if info, err := os.Stat(exe); err == nil {
			return info.ModTime()
		}
	}
	return time.Now()
}

type errSite string

func (e errSite) Error() string { return string(e) }
