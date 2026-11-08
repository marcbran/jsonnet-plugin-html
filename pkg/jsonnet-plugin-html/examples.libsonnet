local html = import './main.libsonnet';
local p = import 'pkg/main.libsonnet';

p.ex({
  example:
    html.manifestHtml({
      element: 'div',
      attributes: { class: 'card' },
      children: [
        { element: 'h1', children: ['Title'] },
        { element: 'p', children: ['Hello World!'] },
      ],
    }),
  output: '<div class="card"><h1>Title</h1><p>Hello World!</p></div>',
}, {
  manifestHtml: p.ex([
    {
      name: 'void elements have no closing tag',
      example: html.manifestHtml({ element: 'br' }),
      output: '<br>',
    },
    {
      name: 'text children are escaped',
      example: html.manifestHtml({ element: 'p', children: ['Tom & Jerry'] }),
      output: '<p>Tom &amp; Jerry</p>',
    },
    {
      name: 'components are objects with an html key',
      example:
        local Greeting(name) = { html: { element: 'strong', children: ['Hello, ' + name + '!'] } };
        html.manifestHtml({ element: 'p', children: [Greeting('World')] }),
      output: '<p><strong>Hello, World!</strong></p>',
    },
    {
      name: 'html key with a string value is raw, unescaped markup',
      example: html.manifestHtml({ element: 'div', children: [{ html: '<b>raw</b>' }] }),
      output: '<div><b>raw</b></div>',
    },
    {
      name: 'null and false children are dropped, arrays are flattened',
      example: html.manifestHtml({
        element: 'ul',
        children: [
          [{ element: 'li', children: ['1'] }, { element: 'li', children: ['2'] }],
          if false then { element: 'li', children: ['3'] },
        ],
      }),
      output: '<ul><li>1</li><li>2</li></ul>',
    },
  ]),
})
