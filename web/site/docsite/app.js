(function () {
  "use strict";
  var d = document;

  var menu = d.querySelector(".menu");
  if (menu) menu.addEventListener("click", function () { d.body.classList.toggle("menu-open"); });

  d.querySelectorAll("pre.code").forEach(function (pre) {
    if (!navigator.clipboard) return;
    var b = d.createElement("button");
    b.type = "button";
    b.className = "copy";
    b.textContent = "Copy";
    b.addEventListener("click", function () {
      var text = pre.querySelector("code").innerText.replace(/^\$ /gm, "");
      navigator.clipboard.writeText(text).then(function () {
        b.textContent = "Copied";
        setTimeout(function () { b.textContent = "Copy"; }, 1500);
      });
    });
    pre.appendChild(b);
  });

  var fb = d.querySelector(".feedback");
  if (fb) fb.addEventListener("click", function (e) {
    var v = e.target.getAttribute("data-helpful");
    if (!v) return;
    fetch(fb.getAttribute("data-endpoint"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ page: location.pathname, helpful: v === "yes" })
    }).catch(function () {});
    fb.querySelectorAll("button, .ask").forEach(function (el) { el.remove(); });
    var t = d.createElement("span");
    t.textContent = v === "yes" ? "Thanks for letting us know!" : "Thanks. We will try to make this page better.";
    fb.insertBefore(t, fb.firstChild);
  });

  var q = d.getElementById("q"), hits = d.getElementById("hits"), index = null, sel = -1;
  if (!q) return;

  function load() {
    if (index) return Promise.resolve(index);
    return fetch(q.getAttribute("data-index")).then(function (r) { return r.json(); }).then(function (j) { index = j; return j; });
  }

  function esc(s) {
    return s.replace(/[&<>"]/g, function (c) { return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]; });
  }

  function snippet(text, word) {
    var i = text.toLowerCase().indexOf(word);
    if (i < 0) return text.slice(0, 110);
    var from = Math.max(0, i - 40);
    return (from > 0 ? "…" : "") + text.slice(from, from + 110) + "…";
  }

  function search() {
    var words = q.value.toLowerCase().split(/\s+/).filter(Boolean);
    if (!words.length) { hits.hidden = true; return; }
    load().then(function (pages) {
      var found = [];
      pages.forEach(function (p) {
        var title = p.t.toLowerCase(), text = p.x.toLowerCase(), score = 0;
        for (var i = 0; i < words.length; i++) {
          var w = words[i], inTitle = title.indexOf(w) >= 0, times = text.split(w).length - 1;
          if (!inTitle && !times) return;
          score += (inTitle ? 10 : 0) + Math.min(times, 5);
        }
        found.push({ p: p, score: score });
      });
      found.sort(function (a, b) { return b.score - a.score; });
      sel = -1;
      hits.innerHTML = found.length
        ? found.slice(0, 8).map(function (f) {
            return '<li><a href="' + f.p.u + '"><b>' + esc(f.p.t) + "</b><span>" + esc(f.p.s + " · " + snippet(f.p.x, words[0])) + "</span></a></li>";
          }).join("")
        : '<li class="none">No results for “' + esc(q.value) + "”</li>";
      hits.hidden = false;
    });
  }

  function move(step) {
    var links = hits.querySelectorAll("a");
    if (!links.length) return;
    if (sel >= 0) links[sel].classList.remove("sel");
    sel = (sel + step + links.length) % links.length;
    links[sel].classList.add("sel");
  }

  q.addEventListener("focus", load);
  q.addEventListener("input", search);
  q.addEventListener("keydown", function (e) {
    if (e.key === "ArrowDown") { e.preventDefault(); move(1); }
    else if (e.key === "ArrowUp") { e.preventDefault(); move(-1); }
    else if (e.key === "Enter") {
      var a = hits.querySelector("a.sel") || hits.querySelector("a");
      if (a) location.href = a.getAttribute("href");
    } else if (e.key === "Escape") { q.value = ""; hits.hidden = true; q.blur(); }
  });
  d.addEventListener("click", function (e) { if (!e.target.closest(".search")) hits.hidden = true; });
  d.addEventListener("keydown", function (e) {
    if (e.key === "/" && d.activeElement !== q && !/input|textarea/i.test(d.activeElement.tagName)) {
      e.preventDefault();
      q.focus();
    }
  });
})();
