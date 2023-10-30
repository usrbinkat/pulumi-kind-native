package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

func main() {
	err := p.RunProvider("kind", "0.1.0", provider())
	if err != nil {
		// Handle error
	}
}

func provider() p.Provider {
	return infer.Provider(infer.Options{
		Resources: []infer.InferredResource{infer.Resource[*Kind, KindArgs, KindState]()},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"kind": "index",
		},
	})
}
