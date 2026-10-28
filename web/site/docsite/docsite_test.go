package docsite

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

var (
	born = time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	now  = time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)
)

func TestMakeIsTheSameForTheSameSeed(t *testing.T) {
	a := Make(42, "example.org", born, now)
	b := Make(42, "example.org", born, now)
	if len(a) != len(b) {
		t.Fatalf("%d files, then %d", len(a), len(b))
	}
	for name, f := range a {
		if g, ok := b[name]; !ok || !bytes.Equal(f.Data, g.Data) || !f.ModTime.Equal(g.ModTime) {
			t.Errorf("%s differs", name)
		}
	}
}

func TestSeedsMakeDifferentSites(t *testing.T) {
	names := map[string]bool{}
	kindsSeen := map[string]bool{}
	for seed := range uint64(40) {
		p := newProject(seed, "example.org", born, now)
		names[p.Name] = true
		kindsSeen[p.kind.id] = true
	}
	if len(names) < 30 {
		t.Errorf("40 seeds made only %d names", len(names))
	}
	if len(kindsSeen) != len(kinds) {
		t.Errorf("40 seeds made %d of %d kinds", len(kindsSeen), len(kinds))
	}
}

var (
	refRe = regexp.MustCompile(`(?:href|src|data-index)="([^"]*)"`)
	idRe  = regexp.MustCompile(` id="([^"]+)"`)
)

// Every link on every page of many sites leads to a file of the site, and
// every anchor to a heading on its page.
func TestLinksLeadSomewhere(t *testing.T) {
	for seed := range uint64(24) {
		fsys := Make(seed, "example.org", born, now)
		for name, f := range fsys {
			if !strings.HasSuffix(name, ".html") {
				continue
			}
			for _, m := range refRe.FindAllStringSubmatch(string(f.Data), -1) {
				ref := m[1]
				if strings.HasPrefix(ref, "https://") {
					continue
				}
				target, anchor, _ := strings.Cut(ref, "#")
				file := name
				if target != "" {
					if !strings.HasPrefix(target, "/") {
						t.Errorf("seed %d %s: relative link %q", seed, name, ref)
						continue
					}
					file = fileOf(target)
				}
				page, ok := fsys[file]
				if !ok {
					t.Errorf("seed %d %s: %q leads nowhere", seed, name, ref)
					continue
				}
				if anchor != "" && !hasID(page.Data, anchor) {
					t.Errorf("seed %d %s: no #%s in %s", seed, name, anchor, file)
				}
			}
		}
	}
}

func fileOf(path string) string {
	path = strings.TrimPrefix(path, "/")
	if path == "" || strings.HasSuffix(path, "/") {
		return path + "index.html"
	}
	return path
}

func hasID(page []byte, id string) bool {
	for _, m := range idRe.FindAllSubmatch(page, -1) {
		if string(m[1]) == id {
			return true
		}
	}
	return false
}

func TestSiteHasTheUsualFiles(t *testing.T) {
	fsys := Make(7, "docs.example.org", born, now)
	for _, name := range []string{"index.html", "404.html", "robots.txt", "sitemap.xml", "favicon.ico", "favicon.svg", "install.sh", "docs/index.html", "docs/installation/index.html", "docs/configuration/index.html", "docs/cli/index.html", "docs/faq/index.html", "changelog/index.html", "about/index.html"} {
		if _, ok := fsys[name]; !ok {
			t.Errorf("no %s", name)
		}
	}
	if !strings.Contains(string(fsys["robots.txt"].Data), "Sitemap: https://docs.example.org/sitemap.xml") {
		t.Errorf("robots.txt: %s", fsys["robots.txt"].Data)
	}
	ico := fsys["favicon.ico"].Data
	if !bytes.HasPrefix(ico, []byte{0, 0, 1, 0, 1, 0, 32, 32}) || !bytes.Contains(ico, []byte("\x89PNG")) {
		t.Errorf("favicon.ico is not an icon")
	}
}

func TestSearchIndexCoversThePages(t *testing.T) {
	fsys := Make(3, "example.org", born, now)
	var index []struct{ T, U, S, X string }
	for name, f := range fsys {
		if strings.HasPrefix(name, "assets/search.") {
			if err := json.Unmarshal(f.Data, &index); err != nil {
				t.Fatal(err)
			}
		}
	}
	pages := 0
	for name := range fsys {
		if strings.HasSuffix(name, "index.html") {
			pages++
		}
	}
	if len(index) != pages {
		t.Fatalf("%d entries for %d pages", len(index), pages)
	}
	for _, e := range index {
		if e.T == "" || e.X == "" || strings.Contains(e.X, "{") && strings.Contains(e.X, "|") {
			t.Errorf("entry %+v", e)
		}
	}
}

// Nothing of the markup or the fillings may show on a page.
func TestPagesHaveNoLeftovers(t *testing.T) {
	leftover := regexp.MustCompile(`\{[a-z_]+\}|\{[^{}]*\|[^{}]*\}|%%|\*\*|\]\(`)
	for seed := range uint64(24) {
		fsys := Make(seed, "example.org", born, now)
		for name, f := range fsys {
			if !strings.HasSuffix(name, ".html") && name != "install.sh" {
				continue
			}
			text := strings.ReplaceAll(string(f.Data), "<code>**</code>", "")
			if name != "install.sh" {
				// Scripts and pictures are not text a reader sees.
				text = regexp.MustCompile(`(?s)<svg.*?</svg>|<script.*?</script>`).ReplaceAllString(text, "")
			}
			if m := leftover.FindString(text); m != "" {
				t.Errorf("seed %d %s: %q", seed, name, m)
			}
		}
	}
}

func TestReleasesOnlyGrow(t *testing.T) {
	r := func(at time.Time) []release {
		return newProject(5, "example.org", born, at).Releases
	}
	before, after := r(now), r(now.AddDate(1, 0, 0))
	if len(after) <= len(before) {
		t.Fatalf("%d releases a year later, %d now", len(after), len(before))
	}
	old := after[len(after)-len(before):]
	for i := range before {
		if before[i].Version != old[i].Version || !before[i].Date.Equal(old[i].Date) {
			t.Fatalf("release %d changed: %v, then %v", i, before[i], old[i])
		}
	}
	if before[0].Date.After(now) {
		t.Errorf("a release in the future: %v", before[0].Date)
	}
}

func TestExpand(t *testing.T) {
	r := newRnd(1, "test")
	vars := map[string]string{"name": "Tern"}
	for _, c := range []struct{ in, want string }{
		{"{name} runs", "Tern runs"},
		{`{"a": 1}`, `{"a": 1}`},
		{"${HOME:-/root}", "${HOME:-/root}"},
		{"{x|x}", "x"},
		{"{{name}|{name}}", "Tern"},
	} {
		if got := expand(r, c.in, vars); got != c.want {
			t.Errorf("%q: %q", c.in, got)
		}
	}
}

func TestRender(t *testing.T) {
	d := render("## Hello `x`\n\nSome **bold** and [a link](/docs/).\n\n- one\n- two\n\n```shell\n$ run\n# note\n```\n\n| A | B |\n| --- | --- |\n| 1 | 2 |")
	for _, want := range []string{`<h2 id="hello-x">`, "<code>x</code>", "<strong>bold</strong>", `<a href="/docs/">a link</a>`, "<li>two</li>", `<span class="p">$ </span>run`, `<span class="c"># note</span>`, "<th>A</th>", "<td>2</td>"} {
		if !strings.Contains(d.HTML, want) {
			t.Errorf("no %q in\n%s", want, d.HTML)
		}
	}
	if len(d.Headings) != 1 || d.Headings[0].Text != "Hello x" {
		t.Errorf("headings %+v", d.Headings)
	}
}

func TestConfigExamples(t *testing.T) {
	for _, ext := range []string{"yaml", "toml"} {
		p := &project{kind: backupKind, Conf: "/etc/x/x." + ext, Ext: ext, vars: map[string]string{"bin": "x", "env": "X_"}}
		got := p.example()
		for _, want := range map[string][]string{
			"yaml": {"repository: s3://backups/web-01\n", "exclude:\n  - \"*.tmp\"\n", "retention:\n  keep_daily: 14\n"},
			"toml": {"repository = \"s3://backups/web-01\"\n", "exclude = [\"*.tmp\", \"/home/*/.cache\"]\n", "\n[retention]\nkeep_daily = 14\n", "\n[upload]\nlimit = \"20MiB/s\""},
		}[ext] {
			if !strings.Contains(got, want) {
				t.Errorf("%s: no %q in\n%s", ext, want, got)
			}
		}
	}
}

var _ fstest.MapFS
