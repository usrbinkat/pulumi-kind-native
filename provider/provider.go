// File: provider/provider.go
package provider

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/resource"
)

// Check validates the resource properties.
func (p *KindClusterProvider) Check(ctx *pulumi.Context, olds, news resource.PropertyMap) (plugin.CheckResult, error) {
	args := KindClusterArgs{
		ClusterName: news["clusterName"].V.(string),
		ConfigFile:  news["configFile"].V.(string),
		WorkingDir:  news["workingDir"].V.(string),
		Purge:       news["purge"].V.(bool),
	}
	if err := ValidateKindClusterArgs(args); err != nil {
		return plugin.CheckResult{}, err
	}
	return plugin.CheckResult{}, nil
}

// Create allocates a new resource instance.
func (p *KindClusterProvider) Create(ctx *pulumi.Context, inputs resource.PropertyMap) (plugin.CreateResult, error) {
	args := KindClusterArgs{
		ClusterName: inputs["clusterName"].V.(string),
		ConfigFile:  inputs["configFile"].V.(string),
		WorkingDir:  inputs["workingDir"].V.(string),
		Purge:       inputs["purge"].V.(bool),
	}

	// Create a Kind cluster
	if err := CreateKindCluster(ctx, args.ClusterName, args); err != nil {
		return plugin.CreateResult{}, fmt.Errorf("Failed to create Kind cluster: %w", err)
	}

	// Create Docker volumes
	if err := CreateVolumesForCluster(ctx, args); err != nil {
		return plugin.CreateResult{}, fmt.Errorf("Failed to create Docker volumes: %w", err)
	}

	return plugin.CreateResult{ID: args.ClusterName, Outputs: inputs}, nil
}

// Delete tears down an existing resource.
func (p *KindClusterProvider) Delete(ctx *pulumi.Context, id pulumi.ID, props resource.PropertyMap) error {
	cluster := KindCluster{
		ClusterName: props["clusterName"].V.(string),
	}

	// Delete Kind cluster
	if err := DeleteKindCluster(ctx, cluster); err != nil {
		return fmt.Errorf("Failed to delete Kind cluster: %w", err)
	}

	// Delete Docker volumes
	args := KindClusterArgs{ClusterName: cluster.ClusterName}
	if err := DeleteVolumesForCluster(ctx, args); err != nil {
		return fmt.Errorf("Failed to delete Docker volumes: %w", err)
	}

	return nil
}

// Diff checks for the differences between the old and new states and returns a DiffResult.
func (p *KindClusterProvider) Diff(ctx *pulumi.Context, id pulumi.ID, olds, news resource.PropertyMap) (plugin.DiffResult, error) {
	diff := plugin.DiffResult{Changes: plugin.DiffNone}

	oldClusterName, newClusterName := olds["clusterName"].V.(string), news["clusterName"].V.(string)
	oldConfig, newConfig := olds["configFile"].V.(string), news["configFile"].V.(string)

	// Check if the cluster name or config file has changed, if so, flag for replacement
	if oldClusterName != newClusterName || oldConfig != newConfig {
		diff.Changes = plugin.DiffSome
		diff.ReplaceKeys = append(diff.ReplaceKeys, "clusterName", "configFile")
	}

	return diff, nil
}

// Read fetches the current state of the resource.
func (p *KindClusterProvider) Read(ctx *pulumi.Context, id pulumi.ID, props resource.PropertyMap) (plugin.ReadResult, error) {
	clusterName := props["clusterName"].V.(string)
	exists, err := CheckIfClusterExists(ctx, clusterName)
	if err != nil {
		return plugin.ReadResult{}, err
	}
	if !exists {
		return plugin.ReadResult{}, nil
	}
	return plugin.ReadResult{
		ID:         id,
		Properties: props,
	}, nil
}

// Update updates an existing resource instance.
func (p *KindClusterProvider) Update(ctx *pulumi.Context, id pulumi.ID, olds, news resource.PropertyMap) (plugin.UpdateResult, error) {
	args := KindClusterArgs{
		ClusterName: news["clusterName"].V.(string),
		ConfigFile:  news["configFile"].V.(string),
		WorkingDir:  news["workingDir"].V.(string),
		Purge:       news["purge"].V.(bool),
	}

	// Update Kind cluster
	if err := UpdateKindCluster(ctx, args.ClusterName, args); err != nil {
		return plugin.UpdateResult{}, fmt.Errorf("Failed to update Kind cluster: %w", err)
	}

	// Update Docker volumes
	if err := UpdateVolumesForCluster(ctx, args); err != nil {
		return plugin.UpdateResult{}, fmt.Errorf("Failed to update Docker volumes: %w", err)
	}

	return plugin.UpdateResult{Outputs: news}, nil
}
