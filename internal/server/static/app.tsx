type Theme = "dark" | "light";

type Snippet = {
  id: number;
  name: string;
  description: string;
  content: string;
  language: string;
  type: string;
  usage_count: number;
  created_at: string;
  updated_at: string;
  tags?: string[];
};

type SnippetFormState = Pick<Snippet, "name" | "description" | "content" | "language" | "type"> & {
  tagsText: string;
};

declare const React: {
  createElement: (...args: unknown[]) => unknown;
  useState: <T>(initial: T) => [T, (next: T | ((value: T) => T)) => void];
  useEffect: (effect: () => void, deps: unknown[]) => void;
};
declare const ReactDOM: { createRoot: (element: HTMLElement | null) => { render: (component: unknown) => void } };

const h = React.createElement;

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

function iconAttrs(viewBox: string) {
  return {
    viewBox,
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

function IconButton(props: { label: string; icon: unknown; className?: string; onClick: () => void }) {
  return h("button", {
    className: `icon-button ${props.className || ""}`,
    type: "button",
    title: props.label,
    "aria-label": props.label,
    onClick: props.onClick
  }, props.icon);
}

function ThemeButton(props: { theme: Theme; onChange: (theme: Theme) => void }) {
  const nextTheme = props.theme === "dark" ? "light" : "dark";
  return h(IconButton, {
    className: "theme-toggle",
    label: `Switch to ${nextTheme} mode`,
    icon: props.theme === "dark" ? h(SunIcon) : h(MoonIcon),
    onClick: () => props.onChange(nextTheme)
  });
}

function NavBar(props: { theme: Theme; onThemeChange: (theme: Theme) => void; onNew: () => void }) {
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

function Sidebar(props: {
  snippets: Snippet[];
  selected: Snippet | null;
  query: string;
  onQuery: (query: string) => void;
  onLoad: (query: string) => void;
  onSelect: (name: string) => void;
}) {
  return h("section", { className: "sidebar" },
    h("div", { className: "row" },
      h("input", {
        "data-focus-key": "search",
        placeholder: "Search snippets",
        value: props.query,
        onInput: (event: Event) => {
          const value = (event.target as HTMLInputElement).value;
          props.onQuery(value);
          props.onLoad(value);
        }
      })
    ),
    h(SnippetList, props)
  );
}

function SnippetList(props: { snippets: Snippet[]; selected: Snippet | null; onSelect: (name: string) => void }) {
  return h("div", { className: "list" }, props.snippets.map((snippet) =>
    h("button", {
      className: `item${props.selected?.name === snippet.name ? " active" : ""}`,
      onClick: () => props.onSelect(snippet.name)
    }, snippet.name, h("br"),
      h("span", { className: "muted" }, snippet.description || "No description"),
      snippet.tags?.length ? h("span", { className: "tag-row" }, snippet.tags.join(", ")) : null
    )
  ));
}

function SnippetDetail(props: { snippet: Snippet; copied: boolean; onCopy: () => void }) {
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
    props.snippet.tags?.length ? h("p", { className: "tags" }, props.snippet.tags.map((tag) =>
      h("span", { className: "tag-pill" }, tag)
    )) : null,
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

function SnippetForm(props: {
  form: SnippetFormState;
  onUpdate: (key: keyof SnippetFormState, value: string) => void;
  onSave: (event: Event) => void;
}) {
  return h("form", { onSubmit: props.onSave },
    h("h2", null, "Add snippet"),
    h("div", { className: "row" },
      h("input", { "data-focus-key": "name", required: true, placeholder: "Name", value: props.form.name, onInput: (e: Event) => props.onUpdate("name", (e.target as HTMLInputElement).value) }),
      h("input", { "data-focus-key": "language", placeholder: "Language", value: props.form.language, onInput: (e: Event) => props.onUpdate("language", (e.target as HTMLInputElement).value) })
    ),
    h("input", { "data-focus-key": "description", placeholder: "Description", value: props.form.description, onInput: (e: Event) => props.onUpdate("description", (e.target as HTMLInputElement).value) }),
    h("input", { "data-focus-key": "tags", placeholder: "Tags, comma separated", value: props.form.tagsText, onInput: (e: Event) => props.onUpdate("tagsText", (e.target as HTMLInputElement).value) }),
    h("p", null, h("textarea", { "data-focus-key": "content", required: true, placeholder: "Snippet content", value: props.form.content, onInput: (e: Event) => props.onUpdate("content", (e.target as HTMLTextAreaElement).value) })),
    h("button", { type: "submit" }, "Save")
  );
}

function App() {
  const [theme, setTheme] = React.useState<Theme>((localStorage.getItem("snape-theme") as Theme) || "dark");
  const [snippets, setSnippets] = React.useState<Snippet[]>([]);
  const [selected, setSelected] = React.useState<Snippet | null>(null);
  const [query, setQuery] = React.useState("");
  const [copied, setCopied] = React.useState(false);
  const [form, setForm] = React.useState<SnippetFormState>({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });

  function load(q: string) {
    fetch(`/api/snippets${q ? `?q=${encodeURIComponent(q)}` : ""}`).then((res) => res.json()).then(setSnippets);
  }

  function selectSnippet(name: string) {
    fetch(`/api/snippets/${encodeURIComponent(name)}`).then((res) => res.json()).then((snippet) => {
      setCopied(false);
      setSelected(snippet);
    });
  }

  function newSnippet() {
    setSelected(null);
    setCopied(false);
  }

  React.useEffect(() => { load(""); }, []);
  React.useEffect(() => {
    document.body.setAttribute("data-theme", theme);
    localStorage.setItem("snape-theme", theme);
  }, [theme]);

  function update(key: keyof SnippetFormState, value: string) {
    setForm({ ...form, [key]: value });
  }

  function save(event: Event) {
    event.preventDefault();
    const { tagsText, ...snippet } = form;
    const payload = { ...snippet, tags: splitTags(tagsText) };
    fetch("/api/snippets", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    }).then((res) => res.json()).then((snippet) => {
      setSelected(snippet);
      setCopied(false);
      setForm({ name: "", description: "", language: "", type: "inline", content: "", tagsText: "" });
      load(query);
    });
  }

  function copySelected() {
    if (!selected) return;
    writeClipboard(selected.content).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 1800);
    });
  }

  return h("div", { className: "shell" },
    h(NavBar, { theme, onThemeChange: setTheme, onNew: newSnippet }),
    h(Sidebar, { snippets, selected, query, onQuery: setQuery, onLoad: load, onSelect: selectSnippet }),
    h("section", { className: "content" },
      selected
        ? h(SnippetDetail, { snippet: selected, copied, onCopy: copySelected })
        : h(SnippetForm, { form, onUpdate: update, onSave: save })
    )
  );
}

function splitTags(value: string) {
  return value.split(",").map((tag) => tag.trim()).filter(Boolean);
}

function writeClipboard(text: string) {
  if (navigator.clipboard?.writeText) {
    return navigator.clipboard.writeText(text);
  }
  const textarea = document.createElement("textarea");
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
