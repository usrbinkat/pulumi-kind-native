// File: provider/volume.go
package provider

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	local "github.com/pulumi/pulumi/sdk/v3/go/pulumi/cmd"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider/cmd/pulumi-resource-command/schema"
)

// ValidateVolumeArgs validates the input arguments for a Docker volume.
func ValidateVolumeArgs(ctx *pulumi.Context, args VolumeArgs) error {
	if err := validateInput(args); err != nil {
		ctx.Log.Error("Validation failed for VolumeArgs", err)
		return err
	}
	ctx.Log.Info("Validation successful for VolumeArgs")
	return nil
}

// CheckIfVolumesExist checks if Docker volumes with the specified names exist.
func CheckIfVolumesExist(ctx *pulumi.Context, volumeNames []string) (bool, error) {
	cmd := local.NewCommand(ctx, "listVolumes", &local.CommandArgs{
		Command:        []string{"docker", "volume", "ls", "--format", "{{.Name}}"},
		StdoutBehavior: schema.CommandBehavior_CAPTURE,
	})
	out, err := cmd.Stdout(ctx)
	if err != nil {
		ctx.Log.Error("Failed to list Docker volumes", err)
		return false, err
	}

	existingVolumes := strings.Split(out, "\n")
	for _, desiredVolume := range volumeNames {
		found := false
		for _, existingVolume := range existingVolumes {
			if existingVolume == desiredVolume {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	return true, nil
}

// CreateVolumesForCluster creates Docker volumes for the given cluster.
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
			ctx.Log.Error("Failed to create Docker volume", err)
			return err
		}
	}

	ctx.Log.Info("Docker volumes created successfully")
	return nil
}

// DeleteVolumesForCluster deletes Docker volumes for the given cluster.
func DeleteVolumesForCluster(ctx *pulumi.Context, volumeNames []string, opts ...pulumi.ResourceOption) error {
	for _, volumeName := range volumeNames {
		_, err := local.NewCommand(ctx, fmt.Sprintf("deleteVolume%s", volumeName), &local.CommandArgs{
			Command:        []string{"docker", "volume", "rm", volumeName},
			StderrBehavior: schema.CommandBehavior_LOG,
			StdoutBehavior: schema.CommandBehavior_LOG,
		}, opts...)

		if err != nil {
			ctx.Log.Error("Failed to delete Docker volume", err)
			return err
		}
	}

	ctx.Log.Info("Docker volumes deleted successfully")
	return nil
}

// UpdateVolumesForCluster updates Docker volumes for the given cluster.
func UpdateVolumesForCluster(ctx *pulumi.Context, args VolumeArgs, opts ...pulumi.ResourceOption) error {
	if err := DeleteVolumesForCluster(ctx, DefaultVolumeNames(), opts...); err != nil {
		return err
	}
	if err := CreateVolumesForCluster(ctx, args, opts...); err != nil {
		return err
	}
	return nil
}

// DefaultVolumeNames returns the default volume names for Kind clusters.
func DefaultVolumeNames() []string {
	return []string{"kind-worker1-containerd", "kind-control1-containerd"}
}
