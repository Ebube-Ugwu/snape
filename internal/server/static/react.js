(function (global) {
  function createElement(type, props) {
    var children = Array.prototype.slice.call(arguments, 2).flat();
    return { type: type, props: props || {}, children: children };
  }

  function useState(initial) {
    var index = state.cursor++;
    if (state.values.length <= index) state.values.push(initial);
    function setValue(next) {
      state.values[index] = typeof next === "function" ? next(state.values[index]) : next;
      renderRoot();
    }
    return [state.values[index], setValue];
  }

  function useEffect(effect, deps) {
    var index = state.cursor++;
    var previous = state.effects[index];
    var changed = !previous || deps.some(function (dep, i) { return dep !== previous[i]; });
    if (changed) {
      state.effects[index] = deps;
      setTimeout(effect, 0);
    }
  }

  function render(vnode, namespace) {
    if (vnode == null || vnode === false) return document.createTextNode("");
    if (typeof vnode === "string" || typeof vnode === "number") return document.createTextNode(vnode);
    if (typeof vnode.type === "function") return render(vnode.type(Object.assign({}, vnode.props, { children: vnode.children })), namespace);
    var nextNamespace = vnode.type === "svg" ? "http://www.w3.org/2000/svg" : namespace;
    var node = nextNamespace ? document.createElementNS(nextNamespace, vnode.type) : document.createElement(vnode.type);
    Object.keys(vnode.props).forEach(function (key) {
      var value = vnode.props[key];
      if (key === "className") node.setAttribute("class", value);
      else if (key.startsWith("on")) node.addEventListener(key.slice(2).toLowerCase(), value);
      else if (key === "value") node.value = value;
      else if (key === "children" || value == null || value === false) return;
      else node.setAttribute(key, value);
    });
    vnode.children.map(function (child) { return render(child, nextNamespace); }).forEach(function (child) { node.appendChild(child); });
    return node;
  }

  function createRoot(element) {
    state.root = element;
    return {
      render: function (component) {
        state.component = component;
        renderRoot();
      }
    };
  }

  function renderRoot() {
    var active = document.activeElement;
    var focusKey = active && active.getAttribute && active.getAttribute("data-focus-key");
    var selection = null;
    if (focusKey && typeof active.selectionStart === "number") {
      selection = { start: active.selectionStart, end: active.selectionEnd };
    }

    state.cursor = 0;
    state.root.replaceChildren(render(state.component));

    if (focusKey) {
      var next = state.root.querySelector('[data-focus-key="' + focusKey + '"]');
      if (next) {
        next.focus();
        if (selection && typeof next.setSelectionRange === "function") {
          next.setSelectionRange(selection.start, selection.end);
        }
      }
    }
  }

  var state = { values: [], effects: [], cursor: 0, root: null, component: null };
  global.React = { createElement: createElement, useState: useState, useEffect: useEffect };
  global.ReactDOM = { createRoot: createRoot };
})(window);
