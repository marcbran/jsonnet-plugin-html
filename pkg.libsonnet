local p = import 'pkg/main.libsonnet';

p.pkg({
  source: 'https://github.com/marcbran/jsonnet-plugin-html',
  repo: 'https://github.com/marcbran/jsonnet.git',
  branch: 'plugin/html',
  path: 'plugin/html',
  target: 'html',
}, |||
  Object-based DSL for creating HTML documents.

  A node is either an element (`{element: 'div', attributes: {...}, children: [...]}`)
  or a reference to other markup (`{html: ...}`), where the `html` value is either
  a string (used verbatim, unescaped) or another node (recursed into). This is how
  components compose: a component is any object with an `html` key, all its other
  keys are ignored by the renderer.

  This library itself doesn't output any HTML strings, and it intentionally has no
  helper functions for building nodes - object literals are cheaper than jsonnet
  function calls. The `manifestHtml` native function takes any value that is valid
  according to this shape and outputs a string in HTML format.
|||, {
  manifestHtml: p.desc(|||
    Renders a node tree to an HTML string.
  |||),
})
