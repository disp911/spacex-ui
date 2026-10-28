package docsite

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// project is the made-up software a site documents.
type project struct {
	kind                 *kind
	Name, Bin, Env, Conf string
	// Ext is the format of the configuration file, yaml or toml.
	Ext    string
	Lang   string
	Domain string
	// Releases are newest first.
	Releases []release
	vars     map[string]string
}

type release struct {
	Version               string
	Date                  time.Time
	Added, Changed, Fixed []string
}

var stems = []string{
	"Alder", "Basalt", "Birch", "Bramble", "Cedar", "Cobalt", "Copper", "Drift",
	"Dune", "Ember", "Fern", "Fjord", "Flint", "Garnet", "Glade", "Harbor",
	"Hazel", "Heath", "Juniper", "Kelp", "Kestrel", "Knoll", "Larch", "Lark",
	"Lumen", "Marlin", "Mesa", "Moss", "Nimbus", "Onyx", "Opal", "Otter",
	"Pebble", "Petrel", "Pike", "Quarry", "Quartz", "Quill", "Reef", "Rook",
	"Rune", "Sable", "Slate", "Sorrel", "Talon", "Tern", "Thistle", "Tide",
	"Umber", "Vale", "Vesper", "Willow", "Wren", "Yarrow",
	"Acorn", "Amber", "Anvil", "Aspen", "Auk", "Badger", "Barley", "Beacon",
	"Bison", "Bolt", "Briar", "Brine", "Bristle", "Canyon", "Cairn", "Caper",
	"Cinder", "Clover", "Comet", "Coral", "Corvid", "Crane", "Crest", "Cypress",
	"Dapple", "Delta", "Dingo", "Dove", "Eddy", "Egret", "Elm", "Falcon",
	"Fathom", "Finch", "Fir", "Flax", "Floe", "Fox", "Gale", "Gannet",
	"Gecko", "Ginger", "Granite", "Grove", "Gull", "Halyard", "Hawthorn", "Hemlock",
	"Hollow", "Ibex", "Ivy", "Jackal", "Jay", "Kiln", "Lagoon", "Laurel",
	"Ledger", "Linden", "Lynx", "Magpie", "Maple", "Meadow", "Merlin", "Mica",
	"Minnow", "Mistral", "Nettle", "Newt", "Oak", "Ochre", "Orchid", "Osprey",
	"Panda", "Pewter", "Pine", "Plover", "Prism", "Puffin", "Quail", "Raven",
	"Ridge", "River", "Robin", "Rowan", "Saffron", "Sage", "Salt", "Sedge",
	"Shale", "Skiff", "Sparrow", "Spruce", "Stoat", "Summit", "Swift", "Tamarack",
	"Teal", "Thorn", "Timber", "Topaz", "Trout", "Tundra", "Vireo", "Walnut",
	"Wasp", "Wicket", "Wolf", "Zinc",
}

func newProject(seed uint64, domain string, born, now time.Time) *project {
	r := newRnd(seed, "project")
	k := pick(r, kinds)
	p := &project{kind: k, Domain: domain}

	stem := pick(r, stems)
	word := strings.ToLower(pick(r, k.words))
	if r.chance(0.5) && len(stem)+len(word) <= 11 {
		p.Name = stem + word
	} else {
		p.Name = stem
	}
	p.Bin = strings.ToLower(p.Name)
	p.Env = strings.ToUpper(p.Bin) + "_"
	p.Ext = pick(r, []string{"yaml", "toml"})
	p.Conf = "/etc/" + p.Bin + "/" + p.Bin + "." + p.Ext
	p.Lang = "Go"
	if r.chance(0.35) {
		p.Lang = "Rust"
	}
	p.Releases = releases(newRnd(seed, "releases"), k, p.Lang, born, now)

	d1 := p.Releases[0].Date.AddDate(0, 0, 3)
	if d1.After(now) {
		d1 = now
	}
	p.vars = map[string]string{
		"name": p.Name, "bin": p.Bin, "env": p.Env, "conf": p.Conf, "ext": p.Ext,
		"domain": domain, "ver": p.Releases[0].Version, "lang": p.Lang, "format": strings.ToUpper(p.Ext),
		"d1": d1.Format("2006-01-02"), "d2": d1.AddDate(0, 0, -1).Format("2006-01-02"),
	}
	vr := newRnd(seed, "vars")
	for _, name := range sortedKeys(k.snips) {
		snip := k.snips[name][0]
		if p.Ext == "toml" {
			snip = k.snips[name][1]
		}
		p.vars[name] = expand(vr, snip, p.vars)
	}
	p.vars["example"] = p.example()
	return p
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// text fills in a text of the kind with the project's words.
func (p *project) text(r *rnd, s string) string {
	return expand(r, s, p.vars)
}

// releases makes the project's history: a first release a while before the
// panel was installed and new ones every few weeks up to now. The history is
// made from its start, so it only grows as time passes.
func releases(r *rnd, k *kind, lang string, born, now time.Time) []release {
	day := time.Date(born.Year(), born.Month(), born.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -r.between(240, 900))
	major, minor, patch := 0, 1, 0
	deck := &changeDeck{r: r, lines: append(append([]string(nil), k.changes...), genericChanges...)}
	var out []release
	for !day.After(now) {
		at := day.Add(time.Duration(r.between(9*60, 20*60)) * time.Minute)
		rel := release{Version: fmt.Sprintf("%d.%d.%d", major, minor, patch), Date: at}
		if patch == 0 {
			rel.Added = deck.deal(r.between(1, 3), "added")
			rel.Changed = deck.deal(r.between(0, 2), "changed")
			rel.Fixed = deck.deal(r.between(0, 2), "fixed")
		} else {
			rel.Fixed = deck.deal(r.between(1, 3), "fixed")
			if r.chance(0.3) {
				rel.Changed = []string{"Updated dependencies"}
			}
		}
		if r.chance(0.12) {
			rel.Changed = append(rel.Changed, toolchain(lang, at))
		}
		out = append(out, rel)

		day = day.AddDate(0, 0, r.between(12, 64))
		switch {
		case r.chance(0.55):
			patch++
		case major == 0 && minor >= 6 && r.chance(0.3):
			major, minor, patch = 1, 0, 0
		default:
			minor, patch = minor+1, 0
		}
	}
	slices.Reverse(out)
	return out
}

// changeDeck deals lines for the changelog without repeating one until all
// of their sort have been dealt.
type changeDeck struct {
	r     *rnd
	lines []string
	dealt map[string][]string
}

func (d *changeDeck) deal(n int, sort string) []string {
	if d.dealt == nil {
		d.dealt = map[string][]string{}
	}
	var out []string
	for range n {
		if len(d.dealt[sort]) == 0 {
			var pool []string
			for _, l := range d.lines {
				if changeSort(l) == sort {
					pool = append(pool, l)
				}
			}
			d.dealt[sort] = shuffled(d.r, pool)
		}
		next := d.dealt[sort][0]
		d.dealt[sort] = d.dealt[sort][1:]
		if !slices.Contains(out, next) {
			out = append(out, next)
		}
	}
	return out
}

func changeSort(line string) string {
	switch {
	case strings.HasPrefix(line, "Fixed"):
		return "fixed"
	case strings.HasPrefix(line, "New") || strings.HasPrefix(line, "Added"):
		return "added"
	}
	return "changed"
}

var genericChanges = []string{
	"Smaller release archives",
	"Documentation fixes and clarifications",
	"Clearer error message when the configuration file cannot be read",
	"`--version` prints the commit the binary was built from",
	"Added release archives for linux/riscv64",
	"Fixed `--log-level` being ignored when set in the environment",
	"New shell completion for fish",
	"Paths in the configuration file may start with `~`",
	"Fixed a crash when `HOME` is not set",
	"Faster start on machines with many CPUs",
	"Log lines name the component that wrote them",
	"Added `--config` to every subcommand",
	"Fixed exit status 0 after a failed configuration reload",
}

// toolchain names the compiler release current at a date.
func toolchain(lang string, at time.Time) string {
	if lang == "Rust" {
		// Rust 1.76 came out on 8 February 2024, and a release follows
		// every six weeks.
		n := 76 + int(at.Sub(time.Date(2024, 2, 8, 0, 0, 0, 0, time.UTC)).Hours()/24/42)
		return fmt.Sprintf("Built with Rust 1.%d", max(n, 60))
	}
	// Go 1.22 came out in February 2024, and a release follows every six
	// months.
	months := (at.Year()-2024)*12 + int(at.Month()) - 2
	return fmt.Sprintf("Built with Go 1.%d", max(22+months/6, 19))
}

// cnode is a key of the configuration file with its value or its keys.
type cnode struct {
	key, value string
	kids       []*cnode
}

func (n *cnode) insert(path []string, value string) {
	for _, k := range n.kids {
		if k.key == path[0] && len(path) > 1 {
			k.insert(path[1:], value)
			return
		}
	}
	kid := &cnode{key: path[0]}
	n.kids = append(n.kids, kid)
	if len(path) == 1 {
		kid.value = value
		return
	}
	kid.insert(path[1:], value)
}

// example is a configuration file with the options that have an example
// value, in the project's format.
func (p *project) example() string {
	root := &cnode{}
	for _, o := range p.kind.options {
		if o.ex != "" {
			root.insert(strings.Split(o.key, "."), expand(nil, o.ex, p.vars))
		}
	}
	var b strings.Builder
	b.WriteString("# " + p.Conf + "\n")
	if p.Ext == "yaml" {
		yamlNode(&b, root, 0)
	} else {
		tomlNode(&b, root, "")
	}
	out := strings.TrimRight(b.String(), "\n")
	for _, name := range p.kind.example {
		out += "\n\n" + p.vars[name]
	}
	return out
}

func yamlNode(b *strings.Builder, n *cnode, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, k := range n.kids {
		if k.kids != nil {
			b.WriteString(indent + k.key + ":\n")
			yamlNode(b, k, depth+1)
			continue
		}
		if strings.HasPrefix(k.value, "[") {
			b.WriteString(indent + k.key + ":\n")
			for _, item := range strings.Split(strings.Trim(k.value, "[]"), ", ") {
				b.WriteString(indent + "  - " + yamlString(item) + "\n")
			}
			continue
		}
		b.WriteString(indent + k.key + ": " + yamlString(k.value) + "\n")
	}
}

// yamlString drops the quotes of a string that YAML reads the same without.
func yamlString(v string) string {
	if !strings.HasPrefix(v, `"`) {
		return v
	}
	s := strings.Trim(v, `"`)
	if s == "" || strings.ContainsAny(s[:1], "*&!%@`'{[>|#-?:,") || strings.Contains(s, ": ") || strings.Contains(s, " #") {
		return v
	}
	return s
}

func tomlNode(b *strings.Builder, n *cnode, table string) {
	for _, k := range n.kids {
		if k.kids == nil {
			b.WriteString(k.key + " = " + k.value + "\n")
		}
	}
	for _, k := range n.kids {
		if k.kids == nil {
			continue
		}
		name := k.key
		if table != "" {
			name = table + "." + k.key
		}
		leaves := false
		for _, kk := range k.kids {
			leaves = leaves || kk.kids == nil
		}
		if leaves {
			b.WriteString("\n[" + name + "]\n")
		}
		tomlNode(b, k, name)
	}
}
