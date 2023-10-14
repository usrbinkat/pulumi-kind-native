// File: provider/kind.go
package provider

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	local "github.com/pulumi/pulumi/sdk/v3/go/pulumi/cmd"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider/cmd/pulumi-resource-command/schema"
)

// ValidateKindClusterArgs validates the input arguments for a KindCluster.
func ValidateKindClusterArgs(ctx *pulumi.Context, args KindClusterArgs) error {
	if err := validateInput(args); err != nil {
		ctx.Log.Error("Validation failed for KindClusterArgs", err)
		return err
	}
	ctx.Log.Info("Validation successful for KindClusterArgs")
	return nil
}

// CheckIfClusterExists checks if a Kind cluster with the specified name exists.
func CheckIfClusterExists(ctx *pulumi.Context, clusterName string) (bool, error) {
	cmd := exec.Command("kind", "get", "clusters")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("Failed to retrieve Kind clusters: %w", err)
	}

	clusters := strings.Split(out.String(), "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return true, nil
		}
	}

	return false, nil
}

// CreateKindCluster creates a new Kind cluster using the Kind CLI.
func CreateKindCluster(ctx *pulumi.Context, name string, args *KindClusterArgs, opts ...pulumi.ResourceOption) error {
	// Validate input arguments
	if err := ValidateKindClusterArgs(ctx, *args); err != nil {
		return err
	}

	cmd := exec.Command("kind", "create", "cluster", "--name", args.ClusterName, "--config", args.ConfigFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to create Kind cluster: %w", err)
	}

	return nil
}

// DeleteKindCluster deletes a Kind cluster using native Pulumi resources.
func DeleteKindCluster(ctx *pulumi.Context, clusterName string, opts ...pulumi.ResourceOption) error {
	_, err := local.NewCommand(ctx, "deleteKindCluster", &local.CommandArgs{
		Delete: &schema.Command{
			Command:        []string{"kind", "delete", "cluster", "--name", clusterName},
			StderrBehavior: schema.CommandBehavior_LOG,
			StdoutBehavior: schema.CommandBehavior_LOG,
		},
	}, opts...)

	if err != nil {
		ctx.Log.Error("Failed to delete Kind cluster", err)
		return err
	}

	ctx.Log.Info("Kind cluster deleted successfully")
	return nil
}

// UpdateKindCluster updates an existing Kind cluster.
func UpdateKindCluster(ctx *pulumi.Context, args *KindClusterArgs) error {
	if err := DeleteKindCluster(ctx, args.ClusterName); err != nil {
		return err
	}
	if err := CreateKindCluster(ctx, args.ClusterName, args); err != nil {
		return err
	}
	return nil
}

// CreateOrUpdateKindCluster creates or updates a Kind cluster using native Pulumi resources.
func CreateOrUpdateKindCluster(ctx *pulumi.Context, name string, args *KindClusterArgs, opts ...pulumi.ResourceOption) error {
	// Validate input arguments
	if err := ValidateKindClusterArgs(ctx, *args); err != nil {
		return err
	}

	exists, err := CheckIfClusterExists(ctx, args.ClusterName)
	if err != nil {
		ctx.Log.Error("Failed to check if Kind cluster exists", err)
		return err
	}

	if exists {
		if err := UpdateKindCluster(ctx, args); err != nil {
			ctx.Log.Error("Failed to update Kind cluster", err)
			return err
		}
		ctx.Log.Info("Kind cluster updated successfully")
	} else {
		if err := CreateKindCluster(ctx, args.ClusterName, args, opts...); err != nil {
			ctx.Log.Error("Failed to create Kind cluster", err)
			return err
		}
		ctx.Log.Info("Kind cluster created successfully")
	}

	return nil
}
