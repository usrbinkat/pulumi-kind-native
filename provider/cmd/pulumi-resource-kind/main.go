// File: provider/cmd/pulumi-resource-kind/main.go
package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pk "github.com/usrbinkat/pulumi-kind-native/provider/kind"
)

func main() {
	err := p.RunProvider("kind", "0.1.0", provider())
	if err != nil {
		panic(err)
	}
}

func provider() p.Provider {
	return infer.Provider(infer.Options{
		Resources: []infer.InferredResource{infer.Resource[*pk.Kind, pk.KindStateArgs]()},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"kind": "index",
		},
	})
}
