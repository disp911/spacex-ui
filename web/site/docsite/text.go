package docsite

import (
	"html"
	"regexp"
	"strings"
)

// expand fills in a text: {name} takes the value of a variable and {a|b|c}
// one of the choices, which may hold more of both. Braces around anything
// else stay as they are.
func expand(r *rnd, s string, vars map[string]string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] != '{' {
			b.WriteByte(s[i])
			i++
			continue
		}
		end := closing(s, i)
		if end < 0 {
			b.WriteString(s[i:])
			break
		}
		inner := s[i+1 : end]
		if parts := splitChoices(inner); len(parts) > 1 {
			b.WriteString(expand(r, pick(r, parts), vars))
		} else if v, ok := vars[inner]; ok {
			b.WriteString(v)
		} else {
			b.WriteString("{" + expand(r, inner, vars) + "}")
		}
		i = end + 1
	}
	return b.String()
}

// closing finds the brace that closes the one at open.
func closing(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// splitChoices splits at the bars that are not inside nested braces.
func splitChoices(s string) []string {
	var parts []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
		case '|':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	return append(parts, s[start:])
}

type heading struct {
	ID, Text string
	Sub      bool
}

// doc is a rendered page body: its HTML, the headings for the page's table
// of contents and its plain text for the search.
type doc struct {
	HTML     string
	Headings []heading
	Text     string
}

// render turns the little markup the pages are written in into HTML:
// "## " and "### " headings, ``` fences, "- " and "1. " lists, "> " notes,
// "| " tables and paragraphs, with `code`, **bold** and [links](/path/)
// inside them.
func render(src string) doc {
	var out, text strings.Builder
	var heads []heading
	lines := strings.Split(strings.TrimSpace(src), "\n")
	for i := 0; i < len(lines); {
		line := strings.TrimRight(lines[i], " ")
		switch {
		case strings.TrimSpace(line) == "":
			i++
		case strings.HasPrefix(line, "```"):
			lang := strings.TrimPrefix(line, "```")
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(lines[j], "```") {
				j++
			}
			out.WriteString(codeBlock(lang, lines[i+1:min(j, len(lines))]))
			i = j + 1
		case strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### "):
			level, title := "h2", strings.TrimPrefix(line, "## ")
			if strings.HasPrefix(line, "### ") {
				level, title = "h3", strings.TrimPrefix(line, "### ")
			}
			id := slug(title)
			heads = append(heads, heading{ID: id, Text: plain(title), Sub: level == "h3"})
			out.WriteString("<" + level + ` id="` + id + `"><a class="anchor" href="#` + id + `" aria-hidden="true">#</a>` + inline(title) + "</" + level + ">\n")
			text.WriteString(plain(title) + ". ")
			i++
		case strings.HasPrefix(line, "- ") || ordered(line):
			tag := "ul"
			if ordered(line) {
				tag = "ol"
			}
			out.WriteString("<" + tag + ">\n")
			for i < len(lines) && (strings.HasPrefix(lines[i], "- ") || ordered(lines[i])) {
				item := lines[i][strings.Index(lines[i], " ")+1:]
				out.WriteString("<li>" + inline(item) + "</li>\n")
				text.WriteString(plain(item) + " ")
				i++
			}
			out.WriteString("</" + tag + ">\n")
		case strings.HasPrefix(line, "> "):
			var parts []string
			for i < len(lines) && strings.HasPrefix(lines[i], "> ") {
				parts = append(parts, strings.TrimPrefix(lines[i], "> "))
				i++
			}
			p := strings.Join(parts, " ")
			out.WriteString(`<div class="note"><p>` + inline(p) + "</p></div>\n")
			text.WriteString(plain(p) + " ")
		case strings.HasPrefix(line, "| "):
			var rows [][]string
			for i < len(lines) && strings.HasPrefix(lines[i], "| ") {
				cells := strings.Split(strings.Trim(strings.TrimSpace(lines[i]), "|"), " | ")
				if !strings.HasPrefix(strings.TrimSpace(cells[0]), "---") {
					rows = append(rows, cells)
				}
				i++
			}
			out.WriteString(table(rows))
			for _, row := range rows {
				for _, c := range row {
					text.WriteString(plain(strings.TrimSpace(c)) + " ")
				}
			}
		default:
			var parts []string
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !block(lines[i]) {
				parts = append(parts, strings.TrimSpace(lines[i]))
				i++
			}
			p := strings.Join(parts, " ")
			out.WriteString("<p>" + inline(p) + "</p>\n")
			text.WriteString(plain(p) + " ")
		}
	}
	return doc{HTML: out.String(), Headings: heads, Text: strings.TrimSpace(text.String())}
}

var orderedRe = regexp.MustCompile(`^\d+\. `)

func ordered(line string) bool { return orderedRe.MatchString(line) }

func block(line string) bool {
	for _, p := range []string{"```", "## ", "### ", "- ", "> ", "| "} {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return ordered(line)
}

func table(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="table"><table>` + "\n<thead><tr>")
	for _, c := range rows[0] {
		b.WriteString("<th>" + inline(strings.TrimSpace(c)) + "</th>")
	}
	b.WriteString("</tr></thead>\n<tbody>\n")
	for _, row := range rows[1:] {
		b.WriteString("<tr>")
		for _, c := range row {
			b.WriteString("<td>" + inline(strings.TrimSpace(c)) + "</td>")
		}
		b.WriteString("</tr>\n")
	}
	b.WriteString("</tbody></table></div>\n")
	return b.String()
}

// codeBlock dims comments and marks prompts in shell, keys in config files.
func codeBlock(lang string, lines []string) string {
	var b strings.Builder
	b.WriteString(`<pre class="code"><code`)
	if lang != "" {
		b.WriteString(` class="language-` + lang + `"`)
	}
	b.WriteString(">")
	for n, line := range lines {
		if n > 0 {
			b.WriteString("\n")
		}
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "--"):
			b.WriteString(`<span class="c">` + html.EscapeString(line) + "</span>")
		case lang == "shell" && strings.HasPrefix(line, "$ "):
			b.WriteString(`<span class="p">$ </span>` + html.EscapeString(line[2:]))
		case lang == "yaml" || lang == "toml" || lang == "ini":
			b.WriteString(configLine(lang, line))
		default:
			b.WriteString(html.EscapeString(line))
		}
	}
	b.WriteString("</code></pre>\n")
	return b.String()
}

func configLine(lang, line string) string {
	sep := ":"
	if lang != "yaml" {
		sep = "="
	}
	if strings.HasPrefix(strings.TrimSpace(line), "[") {
		return `<span class="k">` + html.EscapeString(line) + "</span>"
	}
	key, rest, ok := strings.Cut(line, sep)
	if !ok || strings.ContainsAny(key, `"'`) {
		return html.EscapeString(line)
	}
	value := html.EscapeString(rest)
	if v := strings.TrimSpace(rest); v != "" {
		value = strings.Replace(value, html.EscapeString(v), `<span class="s">`+html.EscapeString(v)+"</span>", 1)
	}
	return `<span class="k">` + html.EscapeString(key) + "</span>" + sep + value
}

var (
	linkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// inline renders the markup inside a paragraph.
func inline(s string) string {
	var b strings.Builder
	for n, part := range strings.Split(s, "`") {
		if n%2 == 1 {
			b.WriteString("<code>" + html.EscapeString(part) + "</code>")
			continue
		}
		part = html.EscapeString(part)
		part = linkRe.ReplaceAllString(part, `<a href="$2">$1</a>`)
		part = boldRe.ReplaceAllString(part, `<strong>$1</strong>`)
		b.WriteString(part)
	}
	return b.String()
}

// plain is the text of the markup without it.
func plain(s string) string {
	s = linkRe.ReplaceAllString(s, "$1")
	s = boldRe.ReplaceAllString(s, "$1")
	return strings.ReplaceAll(s, "`", "")
}

func slug(s string) string {
	s = strings.ToLower(plain(s))
	var b strings.Builder
	dash := false
	for _, c := range s {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
			b.WriteRune(c)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}
