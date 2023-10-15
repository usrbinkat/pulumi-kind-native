// File: provider/cmd/pulumi-resource-kind/main.go

// Package main serves as the entry point for the Pulumi resource provider for Kind.
//
// This package registers the KindClusterProvider resource and starts the Pulumi runtime.
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
)

// The main function serves as the entry point for the Pulumi plugin.
// It first defines the provider using the provider.Main function.
// This function takes the name of our provider and a callback function,
// which itself takes in a provider.HostClient and returns a provider
// that implements the ResourceProvider interface.
func main() {
	// Run the Pulumi runtime with a callback function to register the resource provider.
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Register the KindClusterProvider resource provider.
		return provider.Main("kind", func(host *provider.HostClient) (pulumi.ResourceProvider, error) {
			// Create a new instance of KindClusterProvider for each resource operation.
			return provider.NewProvider(host, "kind", func() pulumi.Resource {
				return &KindClusterProvider{}
			}), nil
		})
	})
}

// NewKindClusterProvider creates a new instance of the KindClusterProvider resource.
// It registers the component resource in the Pulumi runtime and sets its outputs.
// The ctx.RegisterComponentResource function registers the resource in the Pulumi runtime.
// This function takes the resource's Pulumi type token,
// a unique name for the resource, a pointer to the resource itself,
// and additional resource options such as parent and dependencies.
// The ctx.RegisterResourceOutputs function is used to register any output properties for the resource,
// which are those that can be accessed after the resource's creation has been completed.
func NewKindClusterProvider(ctx *pulumi.Context,
	name string, args *KindClusterProviderArgs, opts ...pulumi.ResourceOption) (*KindClusterProvider, error) {

	// Initialize a KindClusterProvider resource instance.
	var resource KindClusterProvider

	// Register the KindClusterProvider as a component resource in the Pulumi runtime.
	err := ctx.RegisterComponentResource("kind:index:KindClusterProvider", name, &resource, opts...)
	if err != nil {
		return nil, err
	}

	// Register resource outputs, currently empty in this case.
	if err := ctx.RegisterResourceOutputs(&resource, pulumi.Map{}); err != nil {
		return nil, err
	}

	return &resource, nil
}
