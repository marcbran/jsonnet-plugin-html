package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestAny(t *testing.T) {
	tests := []struct {
		name    string
		tree    any
		want    string
		wantErr bool
	}{
		{
			name: "empty element",
			tree: map[string]any{"element": "div"},
			want: "<div></div>",
		},
		{
			name: "text child",
			tree: map[string]any{"element": "p", "children": []any{"hello"}},
			want: "<p>hello</p>",
		},
		{
			name: "single non-array child",
			tree: map[string]any{"element": "p", "children": "hello"},
			want: "<p>hello</p>",
		},
		{
			name: "escapes text",
			tree: map[string]any{"element": "p", "children": []any{"Tom & Jerry <3"}},
			want: "<p>Tom &amp; Jerry &lt;3</p>",
		},
		{
			name: "attributes rendered",
			tree: map[string]any{
				"element":    "a",
				"attributes": map[string]any{"href": "https://example.com"},
				"children":   []any{"click"},
			},
			want: `<a href="https://example.com">click</a>`,
		},
		{
			name: "attributes render in sorted order regardless of map order",
			tree: map[string]any{
				"element":    "a",
				"attributes": map[string]any{"href": "https://example.com", "class": "link", "id": "x"},
				"children":   []any{"click"},
			},
			want: `<a class="link" href="https://example.com" id="x">click</a>`,
		},
		{
			name:    "invalid tag name errors",
			tree:    map[string]any{"element": "div onclick=alert(1)"},
			wantErr: true,
		},
		{
			name: "invalid attribute name errors",
			tree: map[string]any{
				"element":    "div",
				"attributes": map[string]any{`x="y" onclick=alert(1)`: "z"},
			},
			wantErr: true,
		},
		{
			name: "boolean true attribute renders bare",
			tree: map[string]any{
				"element":    "input",
				"attributes": map[string]any{"disabled": true},
			},
			want: "<input disabled>",
		},
		{
			name: "boolean false attribute omitted",
			tree: map[string]any{
				"element":    "input",
				"attributes": map[string]any{"disabled": false},
			},
			want: "<input>",
		},
		{
			name: "null attribute omitted",
			tree: map[string]any{
				"element":    "div",
				"attributes": map[string]any{"title": nil},
			},
			want: "<div></div>",
		},
		{
			name: "numeric attribute stringified",
			tree: map[string]any{
				"element":    "input",
				"attributes": map[string]any{"maxlength": float64(5)},
			},
			want: `<input maxlength="5">`,
		},
		{
			name: "attribute value escaped",
			tree: map[string]any{
				"element":    "div",
				"attributes": map[string]any{"title": `say "hi"`},
			},
			want: `<div title="say &#34;hi&#34;"></div>`,
		},
		{
			name: "void element has no closing tag and ignores children",
			tree: map[string]any{"element": "br"},
			want: "<br>",
		},
		{
			name: "nested elements",
			tree: map[string]any{
				"element": "div",
				"children": []any{
					map[string]any{"element": "span", "children": []any{"a"}},
					map[string]any{"element": "span", "children": []any{"b"}},
				},
			},
			want: "<div><span>a</span><span>b</span></div>",
		},
		{
			name: "nested arrays are flattened",
			tree: map[string]any{
				"element": "ul",
				"children": []any{
					[]any{
						map[string]any{"element": "li", "children": []any{"1"}},
						map[string]any{"element": "li", "children": []any{"2"}},
					},
					map[string]any{"element": "li", "children": []any{"3"}},
				},
			},
			want: "<ul><li>1</li><li>2</li><li>3</li></ul>",
		},
		{
			name: "null and false children are dropped",
			tree: map[string]any{
				"element": "div",
				"children": []any{
					"a", nil, false, "b",
				},
			},
			want: "<div>ab</div>",
		},
		{
			name: "component object with html object is rendered",
			tree: map[string]any{
				"element": "div",
				"children": []any{
					map[string]any{
						"html":  map[string]any{"element": "span", "children": []any{"hi"}},
						"props": map[string]any{"ignored": true},
					},
				},
			},
			want: "<div><span>hi</span></div>",
		},
		{
			name: "html string is raw and unescaped",
			tree: map[string]any{
				"element": "div",
				"children": []any{
					map[string]any{"html": "<b>already html</b> & stuff"},
				},
			},
			want: "<div><b>already html</b> & stuff</div>",
		},
		{
			name:    "bare true node errors",
			tree:    true,
			wantErr: true,
		},
		{
			name:    "node with neither element nor html errors",
			tree:    map[string]any{"foo": "bar"},
			wantErr: true,
		},
		{
			name: "node with both element and html errors",
			tree: map[string]any{
				"element": "div",
				"html":    "raw",
			},
			wantErr: true,
		},
		{
			name: "invalid attribute value type errors",
			tree: map[string]any{
				"element":    "div",
				"attributes": map[string]any{"foo": []any{"bad"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ManifestAny(tt.tree)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
