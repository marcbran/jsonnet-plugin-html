package html

import (
	"fmt"
	gohtml "html"
	"sort"
	"strconv"
	"strings"
)

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

func ManifestAny(tree any) (string, error) {
	builder := strings.Builder{}
	err := manifestRec(&builder, tree)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func manifestRec(b *strings.Builder, node any) error {
	switch v := node.(type) {
	case nil:
		return nil
	case bool:
		if v {
			return fmt.Errorf("invalid html node: bare boolean true is not allowed")
		}
		return nil
	case string:
		b.WriteString(gohtml.EscapeString(v))
		return nil
	case []any:
		for _, child := range v {
			err := manifestRec(b, child)
			if err != nil {
				return err
			}
		}
		return nil
	case map[string]any:
		return manifestNode(b, v)
	default:
		return fmt.Errorf("invalid html node: %#v", node)
	}
}

func manifestNode(b *strings.Builder, node map[string]any) error {
	_, hasHtml := node["html"]
	_, hasElement := node["element"]
	_, hasDoctype := node["doctype"]
	kindCount := boolToInt(hasHtml) + boolToInt(hasElement) + boolToInt(hasDoctype)
	switch {
	case kindCount > 1:
		return fmt.Errorf("invalid html node: must have exactly one of 'element', 'html', or 'doctype' keys")
	case hasHtml:
		return manifestHtmlField(b, node["html"])
	case hasElement:
		return manifestElement(b, node)
	case hasDoctype:
		return manifestDoctype(b, node["doctype"])
	default:
		return fmt.Errorf("invalid html node: must have one of 'element', 'html', or 'doctype' keys, got %#v", node)
	}
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func manifestDoctype(b *strings.Builder, doctype any) error {
	name, ok := doctype.(string)
	if !ok {
		return fmt.Errorf("invalid 'doctype' value: must be a string, got %#v", doctype)
	}
	if !isValidName(name) {
		return fmt.Errorf("invalid 'doctype' value: %q is not a valid doctype name", name)
	}
	b.WriteString("<!doctype ")
	b.WriteString(name)
	b.WriteByte('>')
	return nil
}

func manifestHtmlField(b *strings.Builder, html any) error {
	if s, ok := html.(string); ok {
		b.WriteString(s)
		return nil
	}
	return manifestRec(b, html)
}

func manifestElement(b *strings.Builder, node map[string]any) error {
	tag, ok := node["element"].(string)
	if !ok {
		return fmt.Errorf("invalid 'element' value: must be a string, got %#v", node["element"])
	}
	if !isValidName(tag) {
		return fmt.Errorf("invalid 'element' value: %q is not a valid tag name", tag)
	}

	b.WriteByte('<')
	b.WriteString(tag)
	if attrs, ok := node["attributes"]; ok && attrs != nil {
		err := manifestAttributes(b, attrs)
		if err != nil {
			return err
		}
	}
	b.WriteByte('>')

	if voidElements[tag] {
		return nil
	}

	if children, ok := node["children"]; ok {
		err := manifestRec(b, children)
		if err != nil {
			return err
		}
	}

	b.WriteString("</")
	b.WriteString(tag)
	b.WriteByte('>')
	return nil
}

func manifestAttributes(b *strings.Builder, attrsAny any) error {
	attrs, ok := attrsAny.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid 'attributes' value: must be an object, got %#v", attrsAny)
	}
	names := make([]string, 0, len(attrs))
	for name := range attrs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !isValidName(name) {
			return fmt.Errorf("invalid attribute name: %q is not a valid attribute name", name)
		}
		value := attrs[name]
		switch v := value.(type) {
		case nil:
			continue
		case bool:
			if v {
				b.WriteByte(' ')
				b.WriteString(name)
			}
		case string:
			b.WriteByte(' ')
			b.WriteString(name)
			b.WriteString(`="`)
			b.WriteString(gohtml.EscapeString(v))
			b.WriteString(`"`)
		case float64:
			b.WriteByte(' ')
			b.WriteString(name)
			b.WriteString(`="`)
			b.WriteString(formatNumber(v))
			b.WriteString(`"`)
		default:
			return fmt.Errorf("invalid attribute value for %q: %#v", name, value)
		}
	}
	return nil
}

func isValidName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r <= ' ' || strings.ContainsRune(`<>"'=/&`, r) {
			return false
		}
	}
	return true
}

func formatNumber(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}
