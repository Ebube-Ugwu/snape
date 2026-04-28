var h = React.createElement;

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
    h("button", {
      className: "brand",
      type: "button",
      onClick: props.onNew,
      title: "Snape"
    },
      h("span", { className: "brand-mark" }, "S"),
      h("span", { className: "brand-name" }, "Snape")
    ),
    h("div", { className: "nav-actions" },
      h("button", {
        className: "secondary-button",
        type: "button",
        onClick: props.onNew
      }, "Add a new snippet"),
      h(ThemeButton, { theme: props.theme, onChange: props.onThemeChange })
    )
  );
}

function Sidebar(props) {
  return h("section", { className: "sidebar" },
    h("div", { className: "row" },
      h("input", {
        "data-focus-key": "search",
        placeholder: "Search snippets",
        value: props.query,
        onInput: function (event) {
          props.onQuery(event.target.value);
          props.onLoad(event.target.value);
        }
      })
    ),
    h(SnippetList, {
      snippets: props.snippets,
      selected: props.selected,
      onSelect: props.onSelect
    })
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
      return h("span", { className: "tag-pill" }, tag);
    })) : null,
    props.copied ? h("p", { className: "status" }, "Copied to clipboard") : null,
    h("div", { className: "snippet-box" },
      h(IconButton, {
        className: "copy-button",
        label: "Copy snippet",
        icon: h(CopyIcon),
        onClick: props.onCopy
      }),
      h("pre", null, props.snippet.content)
    )
  );
}

function SnippetForm(props) {
  return h("form", { onSubmit: props.onSave },
    h("h2", null, "Add snippet"),
    h("div", { className: "row" },
      h("input", { "data-focus-key": "name", required: true, placeholder: "Name", value: props.form.name, onInput: function (e) { props.onUpdate("name", e.target.value); } }),
      h("input", { "data-focus-key": "language", placeholder: "Language", value: props.form.language, onInput: function (e) { props.onUpdate("language", e.target.value); } })
    ),
    h("input", { "data-focus-key": "description", placeholder: "Description", value: props.form.description, onInput: function (e) { props.onUpdate("description", e.target.value); } }),
    h("input", { "data-focus-key": "tags", placeholder: "Tags, comma separated", value: props.form.tagsText, onInput: function (e) { props.onUpdate("tagsText", e.target.value); } }),
    h("p", null, h("textarea", { "data-focus-key": "content", required: true, placeholder: "Snippet content", value: props.form.content, onInput: function (e) { props.onUpdate("content", e.target.value); } })),
    h("button", { type: "submit" }, "Save")
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
  var stateQuery = React.useState("");
  var query = stateQuery[0];
  var setQuery = stateQuery[1];
  var stateCopied = React.useState(false);
  var copied = stateCopied[0];
  var setCopied = stateCopied[1];
  var stateForm = React.useState({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });
  var form = stateForm[0];
  var setForm = stateForm[1];

  function load(q) {
    fetch("/api/snippets" + (q ? "?q=" + encodeURIComponent(q) : ""))
      .then(function (res) { return res.json(); })
      .then(setSnippets);
  }

  function selectSnippet(name) {
    fetch("/api/snippets/" + encodeURIComponent(name))
      .then(function (res) { return res.json(); })
      .then(function (snippet) {
        setCopied(false);
        setSelected(snippet);
      });
  }

  function newSnippet() {
    setSelected(null);
    setCopied(false);
  }

  React.useEffect(function () { load(""); }, []);
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
    fetch("/api/snippets", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    }).then(function (res) { return res.json(); }).then(function (snippet) {
      setSelected(snippet);
      setCopied(false);
      setForm({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });
      load(query);
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
    h(NavBar, { theme: theme, onThemeChange: setTheme, onNew: newSnippet }),
    h(Sidebar, {
      snippets: snippets,
      selected: selected,
      query: query,
      onQuery: setQuery,
      onLoad: load,
      onSelect: selectSnippet
    }),
    h("section", { className: "content" },
      selected
        ? h(SnippetDetail, { snippet: selected, copied: copied, onCopy: copySelected })
        : h(SnippetForm, { form: form, onUpdate: update, onSave: save })
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
