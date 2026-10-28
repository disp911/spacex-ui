// Package docsite makes the website the panel shows when the administrator
// has no site of their own: the documentation of a small open-source project
// that does not exist, with its own name, look, pages, pictures and release
// history. Everything follows from a seed, so the same seed always makes the
// same site, and every installation gets a site of its own.
package docsite

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"testing/fstest"
	"time"
)

// FeedbackPath is where the pages send «Was this page helpful?». The site
// answers it with 204 and keeps nothing.
const FeedbackPath = "/api/feedback"

//go:embed app.js
var appJS []byte

// Make makes the site of a seed for a domain. born is when the seed was
// made: the project's history starts a while before it and goes on up to
// now, so the site gets new releases as time passes.
func Make(seed uint64, domain string, born, now time.Time) fstest.MapFS {
	p := newProject(seed, domain, born, now)
	th := newTheme(newRnd(seed, "theme"))
	s := &siteData{
		Domain:   domain,
		Name:     p.Name,
		Since:    p.Releases[len(p.Releases)-1].Date.Year(),
		Year:     now.Year(),
		Latest:   p.Releases[0],
		Feedback: FeedbackPath,
		TOC:      th.TOC,
		Color:    th.Accent(),
		theme:    th,
	}
	m := &maker{p: p, s: s, r: newRnd(seed, "pages"), files: fstest.MapFS{}, now: now}
	m.build()
	return m.files
}

type siteData struct {
	Domain, Name   string
	Since, Year    int
	Latest         release
	ReleaseNote    string
	Feedback       string
	TOC            bool
	Color          string
	Logo           template.HTML
	CSS, JS, Index string
	Nav            []navGroup
	Tagline, Lead  string
	Install        template.HTML
	Art            template.HTML
	FeaturesTitle  string
	FeaturesSub    string
	Features       []featureData
	ExampleTitle   string
	ExampleText    string
	Example        template.HTML
	theme          *theme
}

// Years is the span of the copyright line.
func (s *siteData) Years() string {
	if s.Since == s.Year {
		return fmt.Sprint(s.Year)
	}
	return fmt.Sprintf("%d–%d", s.Since, s.Year)
}

type navGroup struct {
	Title string
	Pages []*page
}

type featureData struct {
	Icon, Text template.HTML
	Title      string
}

type page struct {
	Site             *siteData
	Kind             string // home, doc, plain or 404
	Path, Title, Nav string
	Section, Group   string
	Description      string
	Body             template.HTML
	Headings         []heading
	Updated          string
	Prev, Next       *page
	text             string
	date             time.Time
}

func (p *page) FullTitle() string {
	if p.Kind == "home" {
		return p.Site.Name + p.Site.theme.TitleSep + p.Title
	}
	return p.Title + p.Site.theme.TitleSep + p.Site.Name
}

type maker struct {
	p     *project
	s     *siteData
	r     *rnd
	files fstest.MapFS
	now   time.Time
	pages []*page
}

func (m *maker) t(s string) string { return m.p.text(m.r, s) }

func (m *maker) build() {
	p, s, th := m.p, m.s, m.s.theme
	modTime := s.Latest.Date

	light, second := th.Acc(44), th.hsl(th.Sat, 62)
	logo := mark(th.Mark, p.Name[:1], light, second)
	s.Logo = template.HTML(logo)
	m.add("favicon.svg", []byte(logo), modTime)
	m.add("favicon.ico", favicon(th.Hue, th.Sat), modTime)
	s.CSS = m.asset("assets/style", ".css", []byte(th.css()), modTime)
	s.JS = m.asset("assets/app", ".js", appJS, modTime)

	m.docs()
	m.home()
	m.changelog()
	m.about()

	var index []map[string]string
	for _, pg := range m.pages {
		text := pg.text
		if len(text) > 1500 {
			text = text[:1500]
		}
		section := pg.Group
		if section == "" {
			section = pg.Title
		}
		index = append(index, map[string]string{"t": pg.Title, "u": pg.Path, "s": section, "x": text})
	}
	idx, _ := json.Marshal(index)
	s.Index = m.asset("assets/search", ".json", idx, modTime)

	for _, pg := range m.pages {
		name := strings.TrimPrefix(pg.Path, "/") + "index.html"
		m.add(name, m.render(pg), pg.date)
	}
	m.add("404.html", m.render(&page{Site: s, Kind: "404", Path: "/404.html", Title: "Page not found", Description: "Page not found"}), modTime)
	m.add("robots.txt", []byte("User-agent: *\nAllow: /\n\nSitemap: https://"+p.Domain+"/sitemap.xml\n"), modTime)
	m.add("sitemap.xml", m.sitemap(), modTime)
	m.add("install.sh", []byte(m.t(installScript)), modTime)
}

func (m *maker) add(name string, data []byte, mod time.Time) {
	m.files[name] = &fstest.MapFile{Data: data, Mode: 0o444, ModTime: mod.Truncate(time.Second)}
}

// asset adds a file under a name with a hash of its contents, as build
// tools name them, and returns its path.
func (m *maker) asset(base, ext string, data []byte, mod time.Time) string {
	sum := sha256.Sum256(data)
	name := base + "." + hex.EncodeToString(sum[:4]) + ext
	m.add(name, data, mod)
	return "/" + name
}

func (m *maker) render(pg *page) []byte {
	var b bytes.Buffer
	if err := layout.Execute(&b, pg); err != nil {
		panic(err)
	}
	return b.Bytes()
}

// page makes a page from markup and keeps it for the search and the
// sitemap.
func (m *maker) page(kind, path, title, nav, section, group, src string) *page {
	d := render(src)
	body := d.HTML
	if strings.Contains(body, "<p>%%diagram%%</p>") {
		body = strings.Replace(body, "<p>%%diagram%%</p>", `<figure class="diagram">`+diagram(m.p)+`</figure>`, 1)
		d.Text = strings.Replace(d.Text, "%%diagram%%", "", 1)
	}
	// Each page was last touched at one of the recent releases.
	rel := m.p.Releases[m.r.intn(min(len(m.p.Releases), 4))]
	pg := &page{
		Site: m.s, Kind: kind, Path: path, Title: title, Nav: nav, Section: section, Group: group,
		Body: template.HTML(body), Headings: d.Headings, text: d.Text, date: rel.Date,
		Updated:     rel.Date.Format("January 2, 2006"),
		Description: describe(d.Text),
	}
	m.pages = append(m.pages, pg)
	return pg
}

// describe is the first sentence or two of a text, for the meta description.
func describe(text string) string {
	if len(text) <= 160 {
		return text
	}
	cut := strings.LastIndex(text[:160], ". ")
	if cut < 60 {
		cut = strings.LastIndex(text[:157], " ")
		return text[:cut] + "..."
	}
	return text[:cut+1]
}

func (m *maker) docs() {
	k := m.p.kind
	start := m.t("{Getting started|Introduction|Basics|Start here}")
	guides := m.t("{Guides|Usage|Concepts|Using {name}}")
	ref := m.t("{Reference|Reference|API & reference}")
	help := m.t("{Help|Support|More}")

	intro := m.page("doc", "/docs/", m.t("{Introduction|Overview|What is {name}?}"), m.t("{Introduction|Overview}"), "docs", start, m.introText())
	install := m.page("doc", "/docs/installation/", "Installation", "Installation", "docs", start, m.installText())
	quick := m.page("doc", "/docs/quick-start/", m.t("{Quick start|Getting started|Your first {steps|run}}"), "Quick start", "docs", start, m.t(k.quick))
	groups := []navGroup{{start, []*page{intro, install, quick}}}

	var gp []*page
	for _, g := range k.guides {
		nav := g.nav
		if nav == "" {
			nav = g.title
		}
		gp = append(gp, m.page("doc", "/docs/"+g.slug+"/", g.title, nav, "docs", guides, m.t(g.body)))
	}
	gp = append(gp, m.page("doc", "/docs/deployment/", m.t("{Deployment|Running in production|Running as a service}"), "Deployment", "docs", guides, m.deployText()))
	groups = append(groups, navGroup{guides, gp})

	conf := m.page("doc", "/docs/configuration/", "Configuration", "Configuration", "docs", ref, m.configText())
	cli := m.page("doc", "/docs/cli/", m.t("{Command line|Command-line reference|CLI reference}"), m.t("{Command line|CLI}"), "docs", ref, m.cliText())
	groups = append(groups, navGroup{ref, []*page{conf, cli}})

	faq := m.page("doc", "/docs/faq/", m.t("{FAQ|Frequently asked questions}"), "FAQ", "docs", help, m.faqText())
	groups = append(groups, navGroup{help, []*page{faq}})
	m.s.Nav = groups

	var order []*page
	for _, g := range groups {
		order = append(order, g.Pages...)
	}
	for i, pg := range order {
		if i > 0 {
			pg.Prev = order[i-1]
		}
		if i < len(order)-1 {
			pg.Next = order[i+1]
		}
	}
}

func (m *maker) introText() string {
	k := m.p.kind
	var b strings.Builder
	b.WriteString(m.t("{name} is " + pick(m.r, k.what) + ".\n\n"))
	b.WriteString(m.t("It is written in {lang} and ships as a single binary for Linux, macOS and FreeBSD{, on amd64 and arm64|}. {There is nothing else to install: no runtime, no database, no agent.|Releases come out every few weeks; [the changelog](/changelog/) lists what changed.}\n\n"))
	b.WriteString("%%diagram%%\n\n")
	b.WriteString(m.t("## {Features|What it does|Highlights}\n\n"))
	for _, f := range k.features {
		b.WriteString("- **" + f.title + ".** " + m.t(f.text) + "\n")
	}
	b.WriteString(m.t("\n## {Where to go next|Next steps|Find your way}\n\n"))
	b.WriteString(m.t("- [Install {name}](/docs/installation/) on your machine or server.\n"))
	b.WriteString(m.t("- Follow the [quick start](/docs/quick-start/) {to see it working in a few minutes|for a first run}.\n"))
	b.WriteString(m.t("- Look up every option in the [configuration reference](/docs/configuration/).\n"))
	b.WriteString(m.t("- {Questions?|Stuck?} The [FAQ](/docs/faq/) {may already have the answer|answers the common ones}.\n"))
	return b.String()
}

func (m *maker) installText() string {
	p := m.p
	src := `
{name} {is distributed|ships} as a single binary {with no dependencies|with nothing else to install}. Builds are available for Linux, macOS and FreeBSD on amd64 and arm64.

## {Install script|Quick install}

` + "```shell" + `
$ curl -fsSL https://{domain}/install.sh | sh
` + "```" + `

The script {works out|detects} your operating system and architecture, downloads the archive of the latest release, checks its SHA-256 sum and puts ` + "`{bin}`" + ` into ` + "`/usr/local/bin`" + `. {Read it first if you like:|You can read it before running it:} [install.sh](/install.sh).

## From an archive

Every release has archives for each platform, named after the version:

` + "```shell" + `
$ curl -LO https://{domain}/releases/{ver}/{bin}-{ver}-linux-amd64.tar.gz
$ tar -xzf {bin}-{ver}-linux-amd64.tar.gz
$ sudo install -m 755 {bin} /usr/local/bin/
` + "```" + `

## From source

`
	if p.Lang == "Go" {
		src += `Building needs Go 1.22 or newer:

` + "```shell" + `
$ curl -LO https://{domain}/releases/{ver}/{bin}-{ver}.tar.gz
$ tar -xzf {bin}-{ver}.tar.gz && cd {bin}-{ver}
$ go build -o {bin} ./cmd/{bin}
` + "```" + `
`
	} else {
		src += `Building needs a Rust toolchain, 1.74 or newer:

` + "```shell" + `
$ curl -LO https://{domain}/releases/{ver}/{bin}-{ver}.tar.gz
$ tar -xzf {bin}-{ver}.tar.gz && cd {bin}-{ver}
$ cargo build --release
$ sudo install -m 755 target/release/{bin} /usr/local/bin/
` + "```" + `
`
	}
	src += `
## {Check the installation|Verify}

` + "```shell" + `
$ {bin} --version
{bin} {ver} (linux/amd64)
` + "```" + `

## Shell completion

` + "```shell" + `
$ {bin} completion bash | sudo tee /etc/bash_completion.d/{bin} > /dev/null
` + "```" + `

Completion is also available for ` + "`zsh`" + ` and ` + "`fish`" + `. Run ` + "`{bin} completion --help`" + ` for where to put the files.
`
	return m.t(src)
}

func (m *maker) deployText() string {
	k := m.p.kind
	var src string
	if k.daemon == "" {
		src = `
{name} is a command you run, not a service that keeps running. The usual way to run it regularly is a systemd timer.

## {A systemd timer|Running on a schedule}

` + "```ini" + `
# /etc/systemd/system/{bin}.service
[Unit]
Description={name} backup
After=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/{bin} backup /etc /home /srv
ExecStartPost=/usr/local/bin/{bin} prune
Nice=10
IOSchedulingClass=idle
` + "```" + `

` + "```ini" + `
# /etc/systemd/system/{bin}.timer
[Unit]
Description=Daily {name} backup

[Timer]
OnCalendar=*-*-* 02:00
RandomizedDelaySec=30m
Persistent=true

[Install]
WantedBy=timers.target
` + "```" + `

` + "```shell" + `
$ sudo systemctl daemon-reload
$ sudo systemctl enable --now {bin}.timer
$ systemctl list-timers {bin}
` + "```" + `

` + "`Persistent=true`" + ` makes up for a run missed while the machine was off.

## Cron

If you prefer cron, one line does the same:

` + "```shell" + `
0 2 * * * /usr/local/bin/{bin} backup /etc /home /srv && /usr/local/bin/{bin} prune
` + "```" + `
`
	} else {
		src = `
{name} runs as a long-lived service. On most Linux systems that means a systemd unit.

## systemd

Create a user for it and the unit:

` + "```shell" + `
$ sudo useradd --system --home /var/lib/{bin} --create-home --shell /usr/sbin/nologin {bin}
` + "```" + `

` + "```ini" + `
# /etc/systemd/system/{bin}.service
[Unit]
Description={name}
After=network-online.target
Wants=network-online.target

[Service]
User={bin}
ExecStart=/usr/local/bin/{bin} ` + k.daemon + `
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
ProtectSystem=strict
ReadWritePaths=/var/lib/{bin}

[Install]
WantedBy=multi-user.target
` + "```" + `

` + "```shell" + `
$ sudo systemctl daemon-reload
$ sudo systemctl enable --now {bin}
$ journalctl -u {bin} -f
` + "```" + `

{name} reloads its configuration on SIGHUP, so ` + "`systemctl reload {bin}`" + ` applies changes without a restart.

## Containers

Release archives include a ` + "`Dockerfile`" + `. Build the image once and run it with the configuration mounted:

` + "```shell" + `
$ docker build -t {bin}:{ver} .
$ docker run -d --name {bin} --restart unless-stopped -v /etc/{bin}:/etc/{bin}:ro {bin}:{ver} ` + k.daemon + `
` + "```" + `
`
	}
	src += `
## Upgrading

Replace the binary and restart. Within a major version the configuration stays compatible; [the changelog](/changelog/) says when anything needs your attention.
`
	return m.t(src)
}

func (m *maker) configText() string {
	p, k := m.p, m.p.kind
	format := "YAML"
	if p.Ext == "toml" {
		format = "TOML"
	}
	var b strings.Builder
	b.WriteString(m.t("{name} reads its configuration from `{conf}`. Pass `--config` to {use another file|point it elsewhere}. The file is " + format + "{; every option has a default, so it only needs what you change|, and every option in it is optional}.\n\n"))
	b.WriteString(m.t("## {Example|A complete example}\n\n```{ext}\n{example}\n```\n\n"))
	example := k.options[0]
	for _, o := range k.options {
		if strings.Contains(o.key, ".") {
			example = o
			break
		}
	}
	b.WriteString("## Environment variables\n\n")
	b.WriteString(m.t("Every option can also be set in the environment: write the key in upper case, replace dots with underscores and put `{env}` in front, so `" + example.key + "` becomes `" + envName(p.Env, example.key) + "`. The environment {wins over|takes precedence over} the file.\n\n"))
	b.WriteString(m.t("## {Options|All options}\n\n"))
	b.WriteString("| Option | Type | Default | Description |\n| --- | --- | --- | --- |\n")
	for _, o := range k.options {
		def := "–"
		if o.def != "" {
			def = "`" + m.t(o.def) + "`"
		}
		b.WriteString("| `" + o.key + "` | " + o.typ + " | " + def + " | " + m.t(o.desc) + " |\n")
	}
	return b.String()
}

func envName(prefix, key string) string {
	return prefix + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

func (m *maker) cliText() string {
	k := m.p.kind
	var b strings.Builder
	b.WriteString(m.t("All commands take the global flags below. `{bin} help COMMAND` {shows the same information in the terminal|prints the help of a command}.\n\n"))
	b.WriteString("## Global flags\n\n| Flag | Description |\n| --- | --- |\n")
	b.WriteString(m.t("| `--config FILE` | Configuration file. Default: `{conf}`. |\n"))
	b.WriteString("| `--log-level LEVEL` | `debug`, `info`, `warn` or `error`. |\n")
	b.WriteString("| `-h, --help` | Show help for a command. |\n")
	b.WriteString("| `-v, --version` | Print the version and exit. |\n\n")
	cmds := append(append([]command(nil), k.commands...), command{"completion", "SHELL", "Print a shell completion script.", "`SHELL` is `bash`, `zsh` or `fish`.", nil})
	for _, c := range cmds {
		b.WriteString(m.t("## `{bin} " + c.name + "`\n\n"))
		b.WriteString(m.t(c.summary))
		if c.text != "" {
			b.WriteString(" " + m.t(c.text))
		}
		b.WriteString("\n\n```shell\n")
		usage := "{bin} " + c.name
		if len(c.flags) > 0 {
			usage += " [flags]"
		}
		if c.args != "" {
			usage += " " + c.args
		}
		b.WriteString(m.t(usage) + "\n```\n\n")
		if len(c.flags) > 0 {
			b.WriteString("| Flag | Default | Description |\n| --- | --- | --- |\n")
			for _, f := range c.flags {
				name := f.name
				if f.value != "" {
					name += " " + f.value
				}
				def := "–"
				if f.def != "" {
					def = "`" + m.t(f.def) + "`"
				}
				b.WriteString("| `" + name + "` | " + def + " | " + m.t(f.desc) + " |\n")
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}

var genericFAQ = []qa{
	{"Does {name} collect telemetry?", "No. It never contacts anything you did not configure, and it does not check for updates on its own."},
	{"Which platforms are supported?", "Linux, macOS and FreeBSD on amd64 and arm64. Linux on armv7 and riscv64 is built for every release too, but tested less."},
	{"How do I upgrade?", "Replace the binary and restart. Within a major version the configuration stays compatible; [the changelog](/changelog/) points out anything that needs your attention."},
	{"Is {name} free to use at work?", "Yes. It is released under the [MIT License](/about/#license), which allows commercial use."},
}

func (m *maker) faqText() string {
	k := m.p.kind
	var b strings.Builder
	list := append(append([]qa(nil), k.faq...), shuffled(m.r, genericFAQ)[:3]...)
	for _, f := range list {
		b.WriteString(m.t("## " + f.q + "\n\n" + f.a + "\n\n"))
	}
	return b.String()
}

func (m *maker) home() {
	p, s, k := m.p, m.s, m.p.kind
	s.Tagline = pick(m.r, k.taglines)
	what := m.t(pick(m.r, k.what))
	s.Lead = p.Name + " is " + what + "."
	s.ReleaseNote = m.t("{released|is out} " + s.Latest.Date.Format("January 2"))
	s.Install = template.HTML(render(m.t("```shell\n$ curl -fsSL https://{domain}/install.sh | sh\n```")).HTML)
	s.Art = template.HTML(illustration(m.r, s.theme.Illustration, p))
	s.FeaturesTitle = m.t("{Why {name}|Features|What you get|Built for the job}")
	s.FeaturesSub = m.t("{Everything you need, nothing you have to babysit.|Small, predictable and well documented.|The parts that matter, done properly.}")
	for _, f := range shuffled(m.r, k.features) {
		s.Features = append(s.Features, featureData{Icon: template.HTML(icon(f.icon)), Title: f.title, Text: template.HTML(inline(m.t(f.text)))})
	}
	s.ExampleTitle = m.t("{One file to configure it|Configured in one file|Readable configuration}")
	s.ExampleText = m.t("{name} {is configured|is set up} with a single {format} file. Every option has a sensible default and can be overridden from the environment, which makes it easy to run in containers{ and CI|}.")
	s.Example = template.HTML(render("```" + p.Ext + "\n" + p.vars["example"] + "\n```").HTML)

	home := &page{Site: s, Kind: "home", Path: "/", Title: s.Tagline, Section: "home", Description: describe(s.Lead), date: s.Latest.Date}
	home.text = s.Lead
	m.pages = append([]*page{home}, m.pages...)
}

func (m *maker) changelog() {
	p := m.p
	var b strings.Builder
	b.WriteString(m.t("All notable changes to {name} are listed here. {name} follows [semantic versioning](https://semver.org/){; patch releases only fix bugs|}.\n\n"))
	shown := p.Releases[:min(len(p.Releases), 16)]
	dateFormat := pick(m.r, []string{"2006-01-02", "January 2, 2006", "2 January 2006"})
	for _, rel := range shown {
		b.WriteString("## " + rel.Version + " — " + rel.Date.Format(dateFormat) + "\n\n")
		for _, part := range []struct {
			title string
			lines []string
		}{{"Added", rel.Added}, {"Changed", rel.Changed}, {"Fixed", rel.Fixed}} {
			if len(part.lines) == 0 {
				continue
			}
			b.WriteString("### " + part.title + "\n\n")
			for _, l := range part.lines {
				b.WriteString("- " + m.t(l) + "\n")
			}
			b.WriteString("\n")
		}
	}
	if len(p.Releases) > len(shown) {
		b.WriteString(m.t("> Older releases are listed in the `CHANGELOG.md` file {of every source archive|inside the source archive}.\n"))
	}
	pg := m.page("plain", "/changelog/", "Changelog", "Changelog", "changelog", "", b.String())
	pg.date = p.Releases[0].Date
}

func (m *maker) about() {
	p := m.p
	since := p.Releases[len(p.Releases)-1].Date.Year()
	src := `
{name} is ` + pick(m.r, p.kind.what) + `. It {started|began} in ` + fmt.Sprint(since) + ` as {a weekend project to replace a pile of shell scripts|an internal tool at a small hosting company|a replacement for a tool that had grown too large for the job} and is {developed in the open by a small group of maintainers|maintained by a handful of people in their spare time}.

## {Goals|What we care about}

- **{Boring operation|Predictability}.** {name} should do the same thing every time and say clearly when it cannot.
- **{Small surface|Few moving parts}.** One binary, one configuration file, no daemons it depends on.
- **Good documentation.** If something is not written down here, that is a bug.

## Versioning

{name} follows semantic versioning. Within a major version, configuration files and command lines keep working; anything deprecated keeps working for at least two minor releases and is listed in [the changelog](/changelog/).

## License

{name} is free software, released under the MIT License:

` + "```text" + `
Copyright (c) ` + fmt.Sprint(since) + ` The {name} Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
` + "```" + `
`
	m.page("plain", "/about/", m.t("About {name}"), "About", "about", "", m.t(src))
}

func (m *maker) sitemap() []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, pg := range m.pages {
		fmt.Fprintf(&b, "  <url><loc>https://%s%s</loc><lastmod>%s</lastmod></url>\n", m.p.Domain, pg.Path, pg.date.Format("2006-01-02"))
	}
	b.WriteString("</urlset>\n")
	return []byte(b.String())
}

const installScript = `#!/bin/sh
# Installs the latest release of {name}.
#   curl -fsSL https://{domain}/install.sh | sh
set -eu

BASE="https://{domain}/releases"
VERSION="${VERSION:-{ver}}"
PREFIX="${PREFIX:-/usr/local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) echo "{bin}: unsupported architecture $(uname -m)" >&2; exit 1 ;;
esac
case "$os" in
  linux | darwin | freebsd) ;;
  *) echo "{bin}: unsupported operating system $os" >&2; exit 1 ;;
esac

name="{bin}-$VERSION-$os-$arch.tar.gz"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading {name} $VERSION for $os/$arch"
curl -fsSL -o "$tmp/$name" "$BASE/$VERSION/$name"
curl -fsSL -o "$tmp/SHA256SUMS" "$BASE/$VERSION/SHA256SUMS"
(cd "$tmp" && grep " $name\$" SHA256SUMS | sha256sum -c - >/dev/null)

tar -xzf "$tmp/$name" -C "$tmp" {bin}
if [ -w "$PREFIX" ]; then
  install -m 755 "$tmp/{bin}" "$PREFIX/{bin}"
else
  sudo install -m 755 "$tmp/{bin}" "$PREFIX/{bin}"
fi
echo "Installed $("$PREFIX/{bin}" --version)"
`
