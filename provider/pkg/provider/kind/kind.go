// File: provider/pkg/provider/kind/kind.go

// Package provider encapsulates logic for managing Kind clusters using Pulumi.
package provider

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pulumi/pulumi/sdk/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	local "github.com/pulumi/pulumi/sdk/v3/go/pulumi/cmd"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider/cmd/pulumi-resource-command/schema"

	// Import the version package to obtain the current provider version.
	pv "github.com/usrbinkat/pulumi-kind-native/provider/pkg/version"
)

// kindProviderVersion holds the version information of this custom provider.
var kindProviderVersion = pv.GetVersion()

// CreateKindClusterArgs constructs and returns a KindClusterArgs struct
// based on the given PropertyMap inputs.
// It provides default values for unspecified properties.
func CreateKindClusterArgs(inputs resource.PropertyMap) (KindClusterArgs, error) {
	// Initialize arguments with default values.
	clusterName := "pulumi"
	configFile := "kind/config.yaml"
	kindConfigDir, err := os.Getwd()
	if err != nil {
		return KindClusterArgs{}, fmt.Errorf("Failed to get current working directory: %w", err)
	}
	purge := false

	// Override default values if specified in the inputs.
	if cn, ok := inputs["clusterName"].V.(string); ok {
		clusterName = cn
	}
	if cf, ok := inputs["configFile"].V.(string); ok {
		configFile = cf
	}
	if kcd, ok := inputs["kindConfigDir"].V.(string); ok {
		kindConfigDir = kcd
	}
	if p, ok := inputs["purge"].V.(bool); ok {
		purge = p
	}

	return KindClusterArgs{
		ClusterName:   pulumi.String(clusterName),
		ConfigFile:    pulumi.String(configFile),
		KindConfigDir: pulumi.String(kindConfigDir),
		Purge:         pulumi.Bool(purge),
	}, nil
}

// ValidateKindInputs validates the values inside a KindClusterArgs struct.
// It checks whether the provided cluster name, config file, and directory are valid.
//
// Parameters:
// - args: The KindClusterArgs struct containing all the required and optional arguments.
//
// Returns:
// - error: An error object if any of the validations fail.
func ValidateKindInputs(args KindClusterArgs) error {
	// Validate ClusterName
	if args.ClusterName == "" {
		return fmt.Errorf("ClusterName cannot be empty")
	}

	// Validate ConfigFile
	if string(args.ConfigFile) == "" || !strings.HasSuffix(string(args.ConfigFile), ".yaml") {
		return fmt.Errorf("ConfigFile must be a non-empty string with a .yaml extension")
	}

	// Validate KindConfigDir
	if args.KindConfigDir == "" {
		return fmt.Errorf("KindConfigDir cannot be empty")
	}

	// You could add additional file or directory existence checks here if necessary

	return nil
}

// ValidateKindClusterArgs validates the input arguments for creating or updating a Kind cluster.
// It logs an error and returns it if validation fails.
func ValidateKindClusterArgs(ctx *pulumi.Context, args KindClusterArgs) error {
	if err := ValidateKindInputs(args); err != nil {
		errorMsg := fmt.Sprintf("Validation failed for KindClusterArgs: %s", err.Error())
		ctx.Log.Error(errorMsg, nil)
		return err
	}
	ctx.Log.Info("Validation successful for KindClusterArgs", nil)
	return nil
}

// CheckIfClusterExists queries the existing Kind clusters and checks if a cluster with the given name exists.
// It returns a boolean indicating the existence of the cluster and an error if the query operation fails.
func CheckIfClusterExists(ctx *pulumi.Context, clusterName string) (bool, error) {
	// Execute the 'kind get clusters' command to retrieve the list of existing clusters.
	cmd := exec.Command("kind", "get", "clusters")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("Failed to retrieve Kind clusters: %w", err)
	}

	// Parse the output to see if the cluster with the specified name exists.
	clusters := strings.Split(out.String(), "\n")
	for _, cluster := range clusters {
		if cluster == clusterName {
			return true, nil
		}
	}
	return false, nil
}

// CreateKindCluster orchestrates the creation of a new Kind cluster using Pulumi's local.Command.
// It validates the input arguments before proceeding with the cluster creation.
func CreateKindCluster(ctx *pulumi.Context, args *KindClusterArgs, opts ...pulumi.ResourceOption) error {
	// Validate input arguments.
	if err := ValidateKindClusterArgs(ctx, *args); err != nil {
		return err
	}

	// Run the 'kind create cluster' command.
	_, err := local.NewCommand(ctx, "createKindCluster", &local.CommandArgs{
		Create: &schema.Command{
			Command:        []string{"kind", "create", "cluster", "--name", string(args.ClusterName), "--config", string(args.ConfigFile)},
			StderrBehavior: schema.CommandBehavior_LOG,
			StdoutBehavior: schema.CommandBehavior_LOG,
		},
	}, opts...)

	if err != nil {
		ctx.Log.Error("Failed to create Kind cluster", err)
		return err
	}
	ctx.Log.Info("Kind cluster created successfully", nil)
	return nil
}

// DeleteKindCluster deletes an existing Kind cluster.
// It runs the 'kind delete cluster' command.
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
	ctx.Log.Info("Kind cluster deleted successfully", nil)
	return nil
}

// UpdateKindCluster updates an existing Kind cluster by deleting and recreating it.
// This is a simplistic update strategy and could be optimized in the future.
func UpdateKindCluster(ctx *pulumi.Context, args *KindClusterArgs) error {
	// Delete the existing cluster.
	if err := DeleteKindCluster(ctx, string(args.ClusterName)); err != nil {
		return err
	}
	// Recreate the cluster with the updated configuration.
	if err := CreateKindCluster(ctx, args); err != nil {
		return err
	}
	return nil
}

// CreateOrUpdateKindCluster is a high-level function that decides whether to create
// a new Kind cluster or update an existing one based on the given arguments.
// It performs validation, existence checks, and then delegates to the appropriate
// Create or Update method of the KindClusterUtility.
//
// Parameters:
// - ctx: The Pulumi Context object for logging and resource management.
// - args: The arguments required for Kind cluster creation or update.
// - preview: A boolean flag for preview mode.
//
// Returns:
// - error: An error object indicating any issues during the operation.
func CreateOrUpdateKindCluster(ctx *pulumi.Context, args KindClusterArgs, preview bool) error {
	// Initialize KindClusterUtility to use its Create, Update, and Delete methods.
	var clusterUtil ClusterUtility = &KindClusterUtility{}

	// Validate the input arguments before any CRUD operations.
	if err := ValidateKindClusterArgs(ctx, args); err != nil {
		return err
	}

	// Check for the existence of the cluster by its name.
	// This will decide whether to go for a Create or an Update operation.
	exists, err := CheckIfClusterExists(ctx, string(args.ClusterName))
	if err != nil {
		return err
	}

	// Delegate to the appropriate CRUD operation based on the existence of the cluster.
	if exists {
		// Update existing cluster if it already exists.
		return clusterUtil.Update(ctx, args, preview)
	} else {
		// Create a new cluster if it doesn't exist.
		return clusterUtil.Create(ctx, args, preview)
	}
}

// Create is responsible for creating a new Kind cluster.
// It delegates the actual cluster creation to CreateKindCluster.
//
// Parameters:
// - ctx: The Pulumi Context for logging and resource management.
// - args: The KindClusterArgs struct containing all the required and optional arguments.
// - preview: A boolean flag indicating if this is a dry-run.
//
// Returns:
// - error: An error object if the creation fails.
func (k *KindClusterUtility) Create(ctx *pulumi.Context, args KindClusterArgs, preview bool) error {
	// Delegate to the actual cluster creation function.
	return CreateKindCluster(ctx, &args)
}

// Update is responsible for updating an existing Kind cluster.
// It delegates the actual cluster update to UpdateKindCluster.
//
// Parameters:
// - ctx: The Pulumi Context for logging and resource management.
// - args: The KindClusterArgs struct containing all the required and optional arguments.
// - preview: A boolean flag indicating if this is a dry-run.
//
// Returns:
// - error: An error object if the update fails.
func (k *KindClusterUtility) Update(ctx *pulumi.Context, args KindClusterArgs, preview bool) error {
	// Delegate to the actual cluster update function.
	return UpdateKindCluster(ctx, &args)
}

// Delete is responsible for deleting an existing Kind cluster.
// It delegates the actual cluster deletion to DeleteKindCluster.
//
// Parameters:
// - ctx: The Pulumi Context for logging and resource management.
// - clusterName: The name of the cluster to be deleted.
//
// Returns:
// - error: An error object if the deletion fails.
func (k *KindClusterUtility) Delete(ctx *pulumi.Context, clusterName string) error {
	// Delegate to the actual cluster deletion function.
	return DeleteKindCluster(ctx, clusterName)
}
