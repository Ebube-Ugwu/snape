var h = React.createElement;
var pageSize = 10;

function SunIcon() {
  return h("svg", iconAttrs("0 0 24 24"),
    h("circle", { cx: "12", cy: "12", r: "4" }),
    h("path", { d: "M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" })
  );
}

function MoonIcon() {
  return h("svg", iconAttrs("0 0 24 24"),
    h("path", { d: "M20.5 14.5A8.5 8.5 0 0 1 9.5 3.5 7 7 0 1 0 20.5 14.5Z" })
  );
}

function CopyIcon() {
  return h("svg", iconAttrs("0 0 24 24"),
    h("rect", { x: "9", y: "9", width: "11", height: "11", rx: "2" }),
    h("path", { d: "M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" })
  );
}

function EditIcon() {
  return h("svg", iconAttrs("0 0 24 24"),
    h("path", { d: "M12 20h9" }),
    h("path", { d: "M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" })
  );
}

function iconAttrs(viewBox) {
  return {
    viewBox: viewBox,
    width: "18",
    height: "18",
    fill: "none",
    stroke: "currentColor",
    "stroke-width": "2",
    "stroke-linecap": "round",
    "stroke-linejoin": "round",
    "aria-hidden": "true"
  };
}

function IconButton(props) {
  return h("button", {
    className: "icon-button " + (props.className || ""),
    type: "button",
    title: props.label,
    "aria-label": props.label,
    onClick: props.onClick
  }, props.icon);
}

function ThemeButton(props) {
  var nextTheme = props.theme === "dark" ? "light" : "dark";
  return h(IconButton, {
    className: "theme-toggle",
    label: "Switch to " + nextTheme + " mode",
    icon: props.theme === "dark" ? h(SunIcon) : h(MoonIcon),
    onClick: function () { props.onChange(nextTheme); }
  });
}

function NavBar(props) {
  return h("header", { className: "navbar" },
    h("button", { className: "brand", type: "button", onClick: props.onHome, title: "Snape" },
      h("span", { className: "brand-mark" }, "S"),
      h("span", { className: "brand-name" }, "Snape")
    ),
    h("input", {
      className: "nav-search",
      "data-focus-key": "nav-search",
      placeholder: "Search snippets",
      value: props.query,
      onInput: function (event) { props.onSearch(event.target.value); }
    }),
    h("div", { className: "nav-actions" },
      h("button", { className: "secondary-button", type: "button", onClick: props.onTags }, "Tags"),
      h("button", { className: "secondary-button", type: "button", onClick: props.onNew }, "Add a new snippet"),
      h(ThemeButton, { theme: props.theme, onChange: props.onThemeChange })
    )
  );
}

function HomeView(props) {
  var totalPages = Math.max(1, Math.ceil(props.snippets.length / pageSize));
  var page = Math.min(props.page, totalPages);
  var start = (page - 1) * pageSize;
  var visible = props.snippets.slice(start, start + pageSize);
  return h("div", { className: "home-view" },
    h("div", { className: "detail-header" },
      h("div", null,
        h("h2", null, props.query ? "Search results" : "All snippets"),
        h("p", { className: "muted" }, props.query ? props.snippets.length + " matches for " + props.query : "Newest snippets first")
      )
    ),
    visible.length
      ? h(SnippetList, { snippets: visible, selected: props.selected, onSelect: props.onSelect })
      : h("p", { className: "muted" }, "No snippets found"),
    h(Pagination, { page: page, totalPages: totalPages, onPage: props.onPage })
  );
}

function Pagination(props) {
  if (props.totalPages <= 1) return null;
  return h("div", { className: "pagination" },
    h("button", {
      className: "secondary-button",
      type: "button",
      disabled: props.page <= 1,
      onClick: function () { props.onPage(props.page - 1); }
    }, "Previous"),
    h("span", { className: "muted" }, "Page " + props.page + " of " + props.totalPages),
    h("button", {
      className: "secondary-button",
      type: "button",
      disabled: props.page >= props.totalPages,
      onClick: function () { props.onPage(props.page + 1); }
    }, "Next")
  );
}

function SnippetList(props) {
  return h("div", { className: "list" }, props.snippets.map(function (snippet) {
    return h("button", {
      className: "item" + (props.selected && props.selected.name === snippet.name ? " active" : ""),
      onClick: function () { props.onSelect(snippet.name); }
    }, snippet.name, h("br"),
      h("span", { className: "muted" }, snippet.description || "No description"),
      snippet.tags && snippet.tags.length ? h("span", { className: "tag-row" }, snippet.tags.join(", ")) : null
    );
  }));
}

function SnippetDetail(props) {
  return h("div", null,
    h("div", { className: "detail-header" },
      h("div", null,
        h("h2", null, props.snippet.name),
        h("p", { className: "muted" }, props.snippet.description || "No description")
      ),
      h("div", { className: "detail-actions" },
        h("span", { className: "badge" }, props.snippet.language || props.snippet.type)
      )
    ),
    props.snippet.tags && props.snippet.tags.length ? h("p", { className: "tags" }, props.snippet.tags.map(function (tag) {
      return h("button", { className: "tag-pill", type: "button", onClick: function () { props.onTagSelect(tag); } }, tag);
    })) : null,
    props.copied ? h("p", { className: "status" }, "Copied to clipboard") : null,
    h("div", { className: "snippet-box" },
      h("div", { className: "snippet-tools" },
        h(IconButton, { label: "Edit snippet", icon: h(EditIcon), onClick: props.onEdit }),
        h(IconButton, { label: "Copy snippet", icon: h(CopyIcon), onClick: props.onCopy })
      ),
      h("pre", null, props.snippet.content)
    )
  );
}

function SnippetForm(props) {
  return h("form", { onSubmit: props.onSave },
    h("h2", null, props.title),
    h("div", { className: "row" },
      h("input", { "data-focus-key": "name", required: true, placeholder: "Name", value: props.form.name, onInput: function (e) { props.onUpdate("name", e.target.value); } }),
      h("input", { "data-focus-key": "language", placeholder: "Language", value: props.form.language, onInput: function (e) { props.onUpdate("language", e.target.value); } })
    ),
    h("input", { "data-focus-key": "description", placeholder: "Description", value: props.form.description, onInput: function (e) { props.onUpdate("description", e.target.value); } }),
    h("input", { "data-focus-key": "tags", placeholder: "Tags, comma separated", value: props.form.tagsText, onInput: function (e) { props.onUpdate("tagsText", e.target.value); } }),
    h("p", null, h("textarea", { "data-focus-key": "content", required: true, placeholder: "Snippet content", value: props.form.content, onInput: function (e) { props.onUpdate("content", e.target.value); } })),
    h("button", { type: "submit" }, props.submitLabel)
  );
}

function TagsView(props) {
  return h("div", { className: "tags-view" },
    h("div", { className: "detail-header" },
      h("div", null,
        h("h2", null, "Tags"),
        h("p", { className: "muted" }, props.selectedTag ? "Snippets tagged " + props.selectedTag : "Select a tag to filter snippets")
      )
    ),
    h("div", { className: "tag-browser" },
      h("div", { className: "tag-list" }, props.tags.length ? props.tags.map(function (tag) {
        return h("button", { className: "tag-filter" + (props.selectedTag === tag.name ? " active" : ""), type: "button", onClick: function () { props.onTagSelect(tag.name); } }, tag.name);
      }) : h("p", { className: "muted" }, "No tags yet")),
      h("div", { className: "tag-snippets" },
        props.selectedTag
          ? h(SnippetList, { snippets: props.snippets, selected: null, onSelect: props.onSnippetSelect })
          : h("p", { className: "muted" }, "Choose a tag to see matching snippets")
      )
    )
  );
}

function App() {
  var initialTheme = localStorage.getItem("snape-theme") || "dark";
  var stateTheme = React.useState(initialTheme);
  var theme = stateTheme[0];
  var setTheme = stateTheme[1];
  var stateSnippets = React.useState([]);
  var snippets = stateSnippets[0];
  var setSnippets = stateSnippets[1];
  var stateSelected = React.useState(null);
  var selected = stateSelected[0];
  var setSelected = stateSelected[1];
  var stateView = React.useState("home");
  var view = stateView[0];
  var setView = stateView[1];
  var stateQuery = React.useState("");
  var query = stateQuery[0];
  var setQuery = stateQuery[1];
  var statePage = React.useState(1);
  var page = statePage[0];
  var setPage = statePage[1];
  var stateCopied = React.useState(false);
  var copied = stateCopied[0];
  var setCopied = stateCopied[1];
  var stateTags = React.useState([]);
  var tags = stateTags[0];
  var setTags = stateTags[1];
  var stateSelectedTag = React.useState("");
  var selectedTag = stateSelectedTag[0];
  var setSelectedTag = stateSelectedTag[1];
  var stateTagSnippets = React.useState([]);
  var tagSnippets = stateTagSnippets[0];
  var setTagSnippets = stateTagSnippets[1];
  var stateForm = React.useState({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });
  var form = stateForm[0];
  var setForm = stateForm[1];

  function load(q) {
    fetch("/api/snippets" + (q ? "?q=" + encodeURIComponent(q) : ""))
      .then(function (res) { return res.json(); })
      .then(function (items) {
        setSnippets(items);
        setPage(1);
      });
  }

  function goHome() {
    setSelected(null);
    setCopied(false);
    setView("home");
  }

  function search(q) {
    setQuery(q);
    setSelected(null);
    setCopied(false);
    setView("home");
    load(q);
  }

  function selectSnippet(name) {
    fetch("/api/snippets/" + encodeURIComponent(name))
      .then(function (res) { return res.json(); })
      .then(function (snippet) {
        setCopied(false);
        setSelected(snippet);
        setView("snippet");
      });
  }

  function newSnippet() {
    setSelected(null);
    setCopied(false);
    setForm({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });
    setView("form");
  }

  function editSnippet() {
    if (!selected) return;
    setForm({
      name: selected.name,
      description: selected.description || "",
      language: selected.language || "",
      type: selected.type || "inline",
      content: selected.content,
      tagsText: (selected.tags || []).join(", ")
    });
    setCopied(false);
    setView("edit");
  }

  function showTags() {
    setSelected(null);
    setCopied(false);
    setView("tags");
    loadTags();
  }

  function loadTags() {
    fetch("/api/tags").then(function (res) { return res.json(); }).then(setTags);
  }

  function selectTag(tag) {
    setSelectedTag(tag);
    setSelected(null);
    setCopied(false);
    setView("tags");
    loadTags();
    fetch("/api/tags/" + encodeURIComponent(tag) + "/snippets")
      .then(function (res) { return res.json(); })
      .then(setTagSnippets);
  }

  React.useEffect(function () { load(""); loadTags(); }, []);
  React.useEffect(function () {
    document.body.setAttribute("data-theme", theme);
    localStorage.setItem("snape-theme", theme);
  }, [theme]);

  function update(key, value) {
    setForm(Object.assign({}, form, Object.fromEntries([[key, value]])));
  }

  function save(event) {
    event.preventDefault();
    var payload = Object.assign({}, form, { tags: splitTags(form.tagsText) });
    delete payload.tagsText;
    var editing = view === "edit" && selected;
    fetch(editing ? "/api/snippets/" + encodeURIComponent(selected.name) : "/api/snippets", {
      method: editing ? "PUT" : "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    }).then(function (res) { return res.json(); }).then(function (snippet) {
      setSelected(snippet);
      setView("snippet");
      setCopied(false);
      setForm({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });
      load(query);
      loadTags();
    });
  }

  function copySelected() {
    if (!selected) return;
    writeClipboard(selected.content).then(function () {
      setCopied(true);
      setTimeout(function () { setCopied(false); }, 1800);
    });
  }

  return h("div", { className: "shell" },
    h(NavBar, { theme: theme, query: query, onThemeChange: setTheme, onHome: goHome, onNew: newSnippet, onTags: showTags, onSearch: search }),
    h("section", { className: "content" },
      view === "tags"
        ? h(TagsView, { tags: tags, selectedTag: selectedTag, snippets: tagSnippets, onTagSelect: selectTag, onSnippetSelect: selectSnippet })
        : view === "snippet" && selected
        ? h(SnippetDetail, { snippet: selected, copied: copied, onCopy: copySelected, onEdit: editSnippet, onTagSelect: selectTag })
        : view === "form" || view === "edit"
        ? h(SnippetForm, {
            form: form,
            title: view === "edit" ? "Edit snippet" : "Add snippet",
            submitLabel: view === "edit" ? "Update" : "Save",
            onUpdate: update,
            onSave: save
          })
        : h(HomeView, { snippets: snippets, selected: selected, query: query, page: page, onPage: setPage, onSelect: selectSnippet })
    )
  );
}

function splitTags(value) {
  return value.split(",").map(function (tag) { return tag.trim(); }).filter(Boolean);
}

function writeClipboard(text) {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    return navigator.clipboard.writeText(text);
  }
  var textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "true");
  textarea.style.position = "fixed";
  textarea.style.left = "-9999px";
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand("copy");
  document.body.removeChild(textarea);
  return Promise.resolve();
}

ReactDOM.createRoot(document.getElementById("root")).render(h(App));
