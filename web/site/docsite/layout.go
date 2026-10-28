package docsite

import "html/template"

var layout = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.FullTitle}}</title>
<meta name="description" content="{{.Description}}">
<link rel="canonical" href="https://{{.Site.Domain}}{{.Path}}">
<link rel="icon" href="/favicon.ico" sizes="32x32">
<link rel="icon" href="/favicon.svg" type="image/svg+xml">
<link rel="stylesheet" href="{{.Site.CSS}}">
<meta name="theme-color" content="{{.Site.Color}}">
<meta property="og:type" content="website">
<meta property="og:site_name" content="{{.Site.Name}}">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
</head>
<body>
<header class="top">
<div class="wrap">
<a class="brand" href="/">{{.Site.Logo}}<span>{{.Site.Name}}</span></a>
<nav class="main-nav" aria-label="Main">
<a href="/docs/"{{if eq .Section "docs"}} aria-current="page"{{end}}>Docs</a>
<a href="/changelog/"{{if eq .Section "changelog"}} aria-current="page"{{end}}>Changelog</a>
<a href="/about/"{{if eq .Section "about"}} aria-current="page"{{end}}>About</a>
</nav>
<div class="search">
<input type="search" id="q" placeholder="Search docs" autocomplete="off" spellcheck="false" aria-label="Search the documentation" data-index="{{.Site.Index}}">
<kbd>/</kbd>
<ul class="hits" id="hits" hidden></ul>
</div>
{{if eq .Section "docs"}}<button class="menu" type="button" aria-label="Menu"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M4 7h16M4 12h16M4 17h16"/></svg></button>{{end}}
</div>
</header>
{{if eq .Kind "home"}}{{template "home" .}}{{else if eq .Kind "doc"}}{{template "doc" .}}{{else if eq .Kind "404"}}{{template "notfound" .}}{{else}}{{template "plain" .}}{{end}}
<footer class="foot">
<div class="wrap">
<span>© {{.Site.Years}} The {{.Site.Name}} authors. Released under the <a href="/about/#license">MIT License</a>.</span>
<nav aria-label="Footer"><a href="/docs/">Docs</a><a href="/docs/installation/">Install</a><a href="/changelog/">Changelog</a><a href="/docs/faq/">FAQ</a></nav>
</div>
</footer>
<script src="{{.Site.JS}}" defer></script>
</body>
</html>
{{define "home"}}<main>
<section class="hero">
<div class="wrap">
<div>
<a class="release" href="/changelog/"><b>v{{.Site.Latest.Version}}</b>{{.Site.ReleaseNote}}</a>
<h1>{{.Site.Tagline}}</h1>
<p class="lead">{{.Site.Lead}}</p>
<div class="actions"><a class="btn" href="/docs/quick-start/">Get started</a><a class="btn alt" href="/docs/">Read the docs</a></div>
<div class="install">{{.Site.Install}}</div>
</div>
<div class="art">{{.Site.Art}}</div>
</div>
</section>
<section class="band alt">
<div class="wrap">
<h2>{{.Site.FeaturesTitle}}</h2>
<p class="sub">{{.Site.FeaturesSub}}</p>
<div class="features">
{{range .Site.Features}}<div class="feature"><div class="icon">{{.Icon}}</div><h3>{{.Title}}</h3><p>{{.Text}}</p></div>
{{end}}</div>
</div>
</section>
<section class="band">
<div class="wrap example">
<div>
<h2>{{.Site.ExampleTitle}}</h2>
<p>{{.Site.ExampleText}}</p>
<p><a href="/docs/configuration/">Configuration reference →</a></p>
</div>
{{.Site.Example}}
</div>
</section>
</main>{{end}}
{{define "doc"}}<div class="wrap docs">
<nav class="side" aria-label="Documentation">
{{range .Site.Nav}}<p>{{.Title}}</p>
<ul>
{{range .Pages}}<li><a href="{{.Path}}"{{if eq .Path $.Path}} aria-current="page"{{end}}>{{.Nav}}</a></li>
{{end}}</ul>
{{end}}</nav>
<main class="doc">
<article class="prose">
<div class="crumbs"><a href="/docs/">Docs</a> / {{.Group}}</div>
<h1>{{.Title}}</h1>
{{.Body}}
</article>
<div class="feedback" data-endpoint="{{.Site.Feedback}}">
<span class="ask">Was this page helpful?</span>
<button type="button" data-helpful="yes">Yes</button>
<button type="button" data-helpful="no">No</button>
<span class="edited">Last updated {{.Updated}}</span>
</div>
<nav class="pager" aria-label="Pages">
{{with .Prev}}<a class="prev" href="{{.Path}}"><small>Previous</small>{{.Nav}}</a>{{end}}
{{with .Next}}<a class="next" href="{{.Path}}"><small>Next</small>{{.Nav}}</a>{{end}}
</nav>
</main>
{{if .Site.TOC}}<aside class="toc">{{if .Headings}}
<p>On this page</p>
<ul>
{{range .Headings}}<li{{if .Sub}} class="sub"{{end}}><a href="#{{.ID}}">{{.Text}}</a></li>
{{end}}</ul>
{{end}}</aside>{{end}}
</div>{{end}}
{{define "plain"}}<main class="wrap plain prose">
<h1>{{.Title}}</h1>
{{.Body}}
</main>{{end}}
{{define "notfound"}}<main class="notfound">
<div class="code">404</div>
<h1>Page not found</h1>
<p>The page you were looking for does not exist. It may have moved when the documentation was reorganised.</p>
<p><a href="/">Home</a> · <a href="/docs/">Documentation</a></p>
</main>{{end}}
`))
