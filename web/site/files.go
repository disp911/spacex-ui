// Package site is the website the panel shows on port 443 to everyone who is
// not one of its clients: the files of the site, the certificates it is
// served with and the check that it is up. It answers the way an ordinary web
// server does, so nothing about it tells a visitor or a scanner what else
// runs behind the domain.
package site

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"testing/fstest"
	"time"
)

// Files serves a website from fsys: a file at its path, a folder by its
// index.html, a folder asked for without its slash by a redirect to it, and
// anything missing by 404.html with status 404. Only GET and HEAD are
// answered; other methods get 405, as a static server answers them. Text is
// compressed for clients that take gzip, and every file carries an ETag and
// its modification time, so revisits get 304.
func Files(fsys fs.FS) http.Handler {
	return &files{fsys: fsys}
}

type files struct {
	fsys fs.FS
	// gz holds compressed copies by name and ETag, made on first request.
	gz sync.Map
}

func (h *files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		h.status(w, r, http.StatusMethodNotAllowed)
		return
	}
	name, ok := fileName(r.URL.Path)
	if !ok {
		h.status(w, r, http.StatusNotFound)
		return
	}
	info, err := fs.Stat(h.fsys, name)
	if err == nil && info.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			to := r.URL.Path + "/"
			if r.URL.RawQuery != "" {
				to += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, to, http.StatusMovedPermanently)
			return
		}
		name = path.Join(name, "index.html")
		info, err = fs.Stat(h.fsys, name)
	}
	if err != nil || !info.Mode().IsRegular() {
		h.status(w, r, http.StatusNotFound)
		return
	}
	h.serve(w, r, name, info)
}

// fileName turns a request path into a name in the site's files: "." for
// the root, nothing for a path that reaches for a hidden file.
func fileName(urlPath string) (string, bool) {
	clean := path.Clean("/" + urlPath)
	for _, part := range strings.Split(clean, "/") {
		if strings.HasPrefix(part, ".") {
			return "", false
		}
	}
	if clean == "/" {
		return ".", true
	}
	return strings.TrimPrefix(clean, "/"), true
}

func (h *files) serve(w http.ResponseWriter, r *http.Request, name string, info fs.FileInfo) {
	body, err := fs.ReadFile(h.fsys, name)
	if err != nil {
		h.status(w, r, http.StatusNotFound)
		return
	}
	header := w.Header()
	header.Set("Content-Type", contentType(name))
	// The ETag a static server makes: modification time and size in hex.
	etag := fmt.Sprintf(`"%x-%x"`, info.ModTime().Unix(), info.Size())
	if compressible(name) && len(body) > 512 && r.Header.Get("Range") == "" && acceptsGzip(r) {
		header.Add("Vary", "Accept-Encoding")
		body = h.gzipped(name, etag, body)
		header.Set("Content-Encoding", "gzip")
		etag = "W/" + etag
	}
	header.Set("ETag", etag)
	http.ServeContent(w, r, "", info.ModTime(), bytes.NewReader(body))
}

func (h *files) gzipped(name, etag string, body []byte) []byte {
	key := name + " " + etag
	if gz, ok := h.gz.Load(key); ok {
		return gz.([]byte)
	}
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
	zw.Write(body)
	zw.Close()
	h.gz.Store(key, buf.Bytes())
	return buf.Bytes()
}

// status answers with an error page: the site's own 404.html for 404 if it
// has one, a plain page otherwise.
func (h *files) status(w http.ResponseWriter, r *http.Request, code int) {
	body := []byte(fmt.Sprintf("<html>\r\n<head><title>%d %s</title></head>\r\n<body>\r\n<center><h1>%d %s</h1></center>\r\n</body>\r\n</html>\r\n",
		code, http.StatusText(code), code, http.StatusText(code)))
	if code == http.StatusNotFound {
		if page, err := fs.ReadFile(h.fsys, "404.html"); err == nil {
			body = page
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	w.WriteHeader(code)
	if r.Method != http.MethodHead {
		w.Write(body)
	}
}

func contentType(name string) string {
	switch path.Ext(name) {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".txt":
		return "text/plain; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".webmanifest":
		return "application/manifest+json"
	}
	if t := mime.TypeByExtension(path.Ext(name)); t != "" {
		return t
	}
	return "application/octet-stream"
}

func compressible(name string) bool {
	switch path.Ext(name) {
	case ".html", ".htm", ".css", ".js", ".mjs", ".json", ".xml", ".txt", ".svg", ".webmanifest":
		return true
	}
	return false
}

func acceptsGzip(r *http.Request) bool {
	for _, part := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		coding, q, _ := strings.Cut(strings.TrimSpace(part), ";")
		if strings.EqualFold(strings.TrimSpace(coding), "gzip") && strings.ReplaceAll(strings.TrimSpace(q), " ", "") != "q=0" {
			return true
		}
	}
	return false
}

// Source returns the site to show: the administrator's own in dir when it
// has an index.html, made otherwise.
func Source(dir string, made fs.FS) fs.FS {
	own := os.DirFS(dir)
	if info, err := fs.Stat(own, "index.html"); err == nil && info.Mode().IsRegular() {
		return own
	}
	return made
}

// Placeholder is a one-page site shown until the panel makes a real one.
func Placeholder(page []byte, modTime time.Time) fs.FS {
	return fstest.MapFS{"index.html": {Data: page, Mode: 0o444, ModTime: modTime}}
}
