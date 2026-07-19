package html

import (
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

func ManifestHtml() jsonnet.NativeFunction {
	return jsonnet.NativeFunction{
		Name:   "manifestHtml",
		Params: ast.Identifiers{"tree"},
		Func: func(input []any) (any, error) {
			if len(input) != 1 {
				return nil, fmt.Errorf("tree must be provided")
			}
			out, err := ManifestAny(input[0])
			if err != nil {
				return nil, err
			}
			return out, nil
		},
	}
}
