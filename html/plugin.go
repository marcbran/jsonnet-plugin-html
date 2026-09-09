package html

import (
	"github.com/google/go-jsonnet"
	"github.com/marcbran/jpoet/pkg/jpoet"
)

func Plugin(opts ...jpoet.PluginOption) *jpoet.Plugin {
	return jpoet.NewPlugin("html", []jsonnet.NativeFunction{
		ManifestHtml(),
	}, opts...)
}
