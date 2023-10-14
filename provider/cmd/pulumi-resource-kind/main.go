package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
)

func main() {
	// Register the resource provider
	pulumi.Run(func(ctx *pulumi.Context) error {
		return provider.Main("kind", func(host *provider.HostClient) (pulumi.ResourceProvider, error) {
			return provider.NewProvider(host, "kind", func() pulumi.Resource {
				return &KindClusterProvider{}
			}), nil
		})
	})
}

// NewKindClusterProvider creates a new KindClusterProvider resource with the given unique name, arguments, and options.
func NewKindClusterProvider(ctx *pulumi.Context,
	name string, args *KindClusterProviderArgs, opts ...pulumi.ResourceOption) (*KindClusterProvider, error) {
	var resource KindClusterProvider
	err := ctx.RegisterComponentResource("kind:index:KindClusterProvider", name, &resource, opts...)
	if err != nil {
		return nil, err
	}
	if err := ctx.RegisterResourceOutputs(&resource, pulumi.Map{}); err != nil {
		return nil, err
	}
	return &resource, nil
}
