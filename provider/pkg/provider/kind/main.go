// File: provider/pkg/provider/kind/main.go

// Package provider orchestrates the lifecycle of Kind clusters as a Pulumi custom provider.
// It exposes CRUD operations for Kind clusters, encapsulating the underlying Kind CLI and additional logic.
package provider

import (
	"context"
	"fmt"

	// Pulumi SDK packages for building custom providers.
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	gcrp "github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin/provider"

	// Import the version package to set provider version
	pv "github.com/usrbinkat/pulumi-kind-native/provider/pkg/version"
)

// providerVersion contains the semantic version of this custom provider.
var providerVersion = pv.GetVersion()

// Check performs preliminary validation on resources before CRUD operations.
//
// It returns an error if the resource's properties do not pass schema or custom validation.
func (k *KindClusterProvider) Check(ctx context.Context, req *gcrp.CheckRequest, urn resource.URN) (*gcrp.CheckResponse, error) {
	// Convert incoming properties to a KindClusterArgs object.
	args, err := CreateKindClusterArgs(req.News)
	if err != nil {
		return &gcrp.CheckResponse{
			Failures: []plugin.CheckFailure{
				{Property: "kindConfigDir", Reason: fmt.Sprintf("Validation error: %s", err)},
			},
		}, err
	}

	// Perform custom validation on the KindClusterArgs object.
	if err := ValidateKindClusterArgs(k.PulumiCtx, args); err != nil {
		return &gcrp.CheckResponse{
			Failures: []plugin.CheckFailure{
				{Property: "KindClusterArgs", Reason: fmt.Sprintf("Validation error in Check method: %s", err)},
			},
		}, err
	}
	return &gcrp.CheckResponse{}, nil
}

// Create provisions a new Kind cluster.
//
// Takes the CreateRequest object's properties to define the cluster's configuration.
// Returns a CreateResponse with the actual state of the created resource.
func (k *KindClusterProvider) Create(ctx context.Context, req *gcrp.CreateRequest, urn resource.URN) (*gcrp.CreateResponse, error) {
	// Convert incoming properties to a KindClusterArgs object.
	args, err := CreateKindClusterArgs(req.Object)
	if err != nil {
		return nil, fmt.Errorf("Create arguments error in Create method: %s", err)
	}

	// Actual Kind cluster creation call.
	if err := CreateKindCluster(k.PulumiCtx, &args); err != nil {
		return nil, fmt.Errorf("Create cluster error in Create method: %s", err)
	}

	return &gcrp.CreateResponse{
		Id: urn,
	}, nil
}

// Delete handles the deletion of Kind clusters and Docker volumes.
//
// Takes a DeleteRequest object and performs the logic for cluster deletion.
// Returns a DeleteResponse indicating successful deletion.
func (k *KindClusterProvider) Delete(ctx context.Context, req *gcrp.DeleteRequest, urn resource.URN) (*gcrp.DeleteResponse, error) {
	// Extract 'clusterName' from the properties.
	props := req.Properties
	clusterNameProp, ok := props["clusterName"]
	if !ok {
		// 'clusterName' not found, idempotent delete.
		return &gcrp.DeleteResponse{}, nil
	}

	// Validate and cast 'clusterName' to string.
	clusterName, ok := clusterNameProp.(string)
	if !ok {
		return nil, fmt.Errorf("Failed to convert value of 'clusterName' to string in Delete method")
	}

	// Execute the deletion logic.
	if err := DeleteKindCluster(k.PulumiCtx, clusterName); err != nil {
		return nil, err
	}

	// Additional clean-up: Delete associated volumes.
	if err := DeleteVolumesForCluster(k.PulumiCtx, clusterName); err != nil {
		return nil, err
	}

	return &gcrp.DeleteResponse{}, nil
}

// Diff calculates the difference between the old and new resource properties to indicate what changes will be applied.
//
// Diff computes the set of property changes between the old and new states.
// It identifies which properties will be replaced or require recreation of the resource.
func (k *KindClusterProvider) Diff(ctx context.Context, req *gcrp.DiffRequest, urn resource.URN) (*gcrp.DiffResponse, error) {
	// Initialize the DiffResponse object.
	diff := &gcrp.DiffResponse{}

	// List of properties to check for differences.
	keysToCheck := []string{"clusterName", "configFile"}
	olds := req.Olds
	news := req.News

	// Loop through the keys and identify differences.
	for _, key := range keysToCheck {
		newProp, newExists := news[key]
		oldProp, oldExists := olds[key]

		// If both old and new properties exist, compare.
		if newExists && oldExists {
			newPropStr, ok1 := newProp.(string)
			oldPropStr, ok2 := oldProp.(string)

			if ok1 && ok2 && newPropStr != oldPropStr {
				diff.ReplaceKeys = append(diff.ReplaceKeys, key)
			}
		}
	}

	// Update the DiffResponse based on identified changes.
	if len(diff.ReplaceKeys) > 0 {
		diff.Changes = gcrp.DiffSome
	}

	return diff, nil
}

// Read retrieves the current state of a particular resource.
//
// Returns a ReadResponse containing the current properties of the resource.
// Returns an error if the cluster does not exist or cannot be read.
func (k *KindClusterProvider) Read(ctx context.Context, req *gcrp.ReadRequest, urn resource.URN) (*gcrp.ReadResponse, error) {
	// Assume CheckIfClusterExists is implemented elsewhere
	props := req.Properties
	clusterNameProp, ok := props["clusterName"]
	if !ok {
		return nil, fmt.Errorf("The 'clusterName' property was not found in the Read method")
	}

	clusterName, ok := clusterNameProp.(string)
	if !ok {
		return nil, fmt.Errorf("Failed to convert value of 'clusterName' to string in Read method")
	}

	// Check for the existence of the cluster using the Pulumi context.
	exists, err := CheckIfClusterExists(k.PulumiCtx, clusterName)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("Cluster does not exist in Read method")
	}

	return &gcrp.ReadResponse{
		Id:         req.Id,
		Properties: props,
	}, nil
}

// Update modifies an existing Kind cluster based on new properties.
//
// Returns an UpdateResponse object containing the new state of the resource.
// Returns an error if the update operation fails.
func (k *KindClusterProvider) Update(ctx context.Context, req *gcrp.UpdateRequest, urn resource.URN) (*gcrp.UpdateResponse, error) {
	// Convert incoming properties to a KindClusterArgs object.
	args, err := CreateKindClusterArgs(req.News)
	if err != nil {
		return nil, fmt.Errorf("Update arguments error in Update method: %s", err)
	}

	// Actual Kind cluster update logic using the Pulumi context.
	updateErr := UpdateKindCluster(k.PulumiCtx, &args)
	if updateErr != nil {
		return nil, fmt.Errorf("Update cluster error in Update method: %s", updateErr)
	}

	// Convert KindClusterArgs to VolumeArgs for the volume update logic.
	volumeArgs := VolumeArgs{
		ClusterName: string(args.ClusterName), // Assuming ClusterName is a string; adjust as needed.
		// Populate other VolumeArgs fields here if necessary.
	}

	// Additional logic: Update associated volumes using the Pulumi context.
	if err := UpdateVolumesForCluster(k.PulumiCtx, volumeArgs); err != nil {
		return nil, fmt.Errorf("Update volume error in Update method: %s", err)
	}

	return &gcrp.UpdateResponse{}, nil
}

// GetSchema returns the JSON-serializable schema of this provider.
func (k *KindClusterProvider) GetSchema(ctx context.Context) (*gcrp.GetSchemaResponse, error) {
	schema := `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"clusterName": {
				"type": "string",
				"description": "The name of the Kind cluster"
			},
			"configFile": {
				"type": "string",
				"description": "Path to the Kind configuration file"
			}
		},
		"required": ["clusterName"]
	}`

	return &gcrp.GetSchemaResponse{
		Schema: schema,
	}, nil
}

// main is the entry point for this custom Pulumi provider.
// It sets up the plugin service and registers the custom KindClusterProvider.
//func main() {
//	plugin.Serve(&plugin.ServeOpts{
//		// ProviderFunc provides an instance of the custom KindClusterProvider.
//		ProviderFunc: func(_ string) (gcrp.Provider, error) {
//			return &KindClusterProvider{}, nil
//		},
//	})
//}

//// Construct - create a new component resource.
//func (k *kindClusterProvider) Construct(ctx context.Context, req *gcrp.ConstructRequest, urn resource.URN) (*gcrp.ConstructResponse, error) {
//	// Example: Initialize a new Kind cluster component with default settings.
//	// Note: In a real implementation, you'd typically look at `req.Args` and `req.Options` to customize the behavior.
//	return &gcrp.ConstructResponse{
//		URN:  urn,
//		State: map[string]interface{}{
//			"defaultComponent": "kindCluster",
//		},
//	}, nil
//}
//
//// GetPluginInfo - return version info about this provider
//func (k *kindClusterProvider) GetPluginInfo(ctx context.Context) (*gcrp.GetPluginInfoResponse, error) {
//	return &gcrp.GetPluginInfoResponse{Version: providerVersion}, nil
//}
//
//// Invoke - invoke a method that performs an action but doesn't manage resources.
//func (k *kindClusterProvider) Invoke(ctx context.Context, req *gcrp.InvokeRequest) (*gcrp.InvokeResponse, error) {
//	// Example: Implement a "ping" method that checks if the provider is responsive.
//	if req.Tok == "ping" {
//		return &gcrp.InvokeResponse{
//			Outputs: map[string]interface{}{
//				"message": "pong",
//			},
//		}, nil
//	}
//	return nil, fmt.Errorf("Unknown function: %s", req.Tok)
//}
//
//// Call - call a method on a component resource.
//func (k *kindClusterProvider) Call(ctx context.Context, req *gcrp.CallRequest) (*gcrp.CallResponse, error) {
//	// Example: Implement a "status" method that provides some status information about a Kind cluster component.
//	if req.Tok == "status" {
//		return &gcrp.CallResponse{
//			Outputs: map[string]interface{}{
//				"status": "Cluster is healthy", // You'd fetch this information dynamically.
//			},
//		}, nil
//	}
//	return nil, fmt.Errorf("Unknown method: %s", req.Tok)
//}
