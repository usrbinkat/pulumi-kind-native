// File: provider/pkg/provider/kind/volume.go

// Package provider encapsulates logic for managing Docker volumes for Kind clusters.
package provider

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	local "github.com/pulumi/pulumi/sdk/v3/go/pulumi/cmd"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider/cmd/pulumi-resource-command/schema"
)

// DefaultVolumeNames provides the default Docker volume names that should be created for Kind clusters.
func DefaultVolumeNames() []string {
	return []string{"kind-worker1-containerd", "kind-control1-containerd"}
}

// validateInput performs basic validation on the VolumeArgs.
// For example, it could check if the cluster name is not empty and if at least one volume name is provided.
func ValidateVolumeInputs(args VolumeArgs) error {
	if args.ClusterName == "" {
		return fmt.Errorf("ClusterName cannot be empty")
	}
	if len(args.VolumeNames) == 0 {
		return fmt.Errorf("At least one volume name must be provided")
	}
	for _, name := range args.VolumeNames {
		if name == "" {
			return fmt.Errorf("Volume names cannot be empty")
		}
	}
	return nil
}

// ValidateVolumeArgs performs validation on the input parameters for Docker volumes.
// It uses a separate validation function and logs the outcome.
func ValidateVolumeArgs(ctx *pulumi.Context, args VolumeArgs) error {
	if err := ValidateVolumeInputs(args); err != nil {
		ctx.Log.Error(fmt.Sprintf("Validation failed for VolumeArgs: %s", err.Error()), nil)
		return err
	}
	ctx.Log.Info("Validation successful for VolumeArgs", nil)
	return nil
}

// CheckIfVolumesExist queries existing Docker volumes to confirm if the specified volumes exist.
// Returns true only if all desired volumes are found.
func CheckIfVolumesExist(ctx *pulumi.Context, volumeNames []string) (bool, error) {
	for _, volumeName := range volumeNames {
		// Construct the command to list a specific Docker volume by name
		cmd := local.NewCommand(ctx, "listVolume", &local.CommandArgs{
			Command:        []string{"docker", "volume", "ls", "--filter", fmt.Sprintf("name=%s", volumeName), "--format", "{{.Name}}"},
			StdoutBehavior: schema.CommandBehavior_CAPTURE,
		})
		out, err := cmd.Stdout(ctx)
		if err != nil {
			ctx.Log.Error(fmt.Sprintf("Failed to list Docker volume: %s", err.Error()), nil)
			return false, err
		}

		// Check if the volume exists
		if strings.TrimSpace(out) != volumeName {
			return false, nil
		}
	}

	return true, nil
}

// CreateVolumesForCluster orchestrates the creation of Docker volumes for the Kind cluster.
// It validates the arguments and then iterates over the default volume names to create each.
func CreateVolumesForCluster(ctx *pulumi.Context, args VolumeArgs, opts ...pulumi.ResourceOption) error {
	if err := ValidateVolumeArgs(ctx, args); err != nil {
		return err
	}

	for _, volumeName := range DefaultVolumeNames() {
		_, err := local.NewCommand(ctx, fmt.Sprintf("createVolume%s", volumeName), &local.CommandArgs{
			Command:        []string{"docker", "volume", "create", fmt.Sprintf("--name=%s", volumeName)},
			StderrBehavior: schema.CommandBehavior_LOG,
			StdoutBehavior: schema.CommandBehavior_LOG,
		}, opts...)

		if err != nil {
			ctx.Log.Error(fmt.Sprintf("Failed to create Docker volume: %s", err.Error()), nil)
			return err
		}
	}

	ctx.Log.Info("Docker volumes created successfully", nil)
	return nil
}

// DeleteVolumesForCluster orchestrates the deletion of Docker volumes.
// It iterates over the provided list of volume names and deletes each.
func DeleteVolumesForCluster(ctx *pulumi.Context, volumeNames []string, opts ...pulumi.ResourceOption) error {
	for _, volumeName := range volumeNames {
		_, err := local.NewCommand(ctx, fmt.Sprintf("deleteVolume%s", volumeName), &local.CommandArgs{
			Command:        []string{"docker", "volume", "rm", volumeName},
			StderrBehavior: schema.CommandBehavior_LOG,
			StdoutBehavior: schema.CommandBehavior_LOG,
		}, opts...)

		if err != nil {
			ctx.Log.Error(fmt.Sprintf("Failed to delete Docker volume: %s", err.Error()), nil)
			return err
		}
	}

	ctx.Log.Info("Docker volumes deleted successfully", nil)
	return nil
}

// UpdateVolumesForCluster updates Docker volumes for the Kind cluster.
// It does so by first deleting existing volumes and then creating new ones.
func UpdateVolumesForCluster(ctx *pulumi.Context, args VolumeArgs, opts ...pulumi.ResourceOption) error {
	if err := DeleteVolumesForCluster(ctx, DefaultVolumeNames(), opts...); err != nil {
		return err
	}
	if err := CreateVolumesForCluster(ctx, args, opts...); err != nil {
		return err
	}
	return nil
}
