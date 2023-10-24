// File: main.go
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/usrbinkat/pulumi-kind-native/sdk/v3/go/kind"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kind.NewKindCluster(ctx, "cluster", &kind.KindClusterArgs{
			ClusterName:   pulumi.String("my-cluster"),
			ConfigFile:    pulumi.String("kind/config.yaml"),
			KindConfigDir: pulumi.String("."),
			Purge:         pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		ctx.Export("TestClusterEndpoint", testCluster.Endpoint)
		return nil
	})
}
