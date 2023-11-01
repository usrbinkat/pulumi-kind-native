package main

import (
	"github.com/pulumi/pulumi-kind/sdk/go/kind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kind.NewKind(ctx, "myArbitraryKindClusterResourceName", &kind.KindArgs{
			Name: pulumi.String("test"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
