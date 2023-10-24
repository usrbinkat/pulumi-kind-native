// File: provider/pkg/provider/kind/types.go

// Package provider defines types and utilities for creating a Pulumi Kind provider.
package provider

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// kindClusterProvider is the main implementation of the custom Pulumi provider.
//
// This struct is used for creating and managing KindCluster resources.
type KindClusterProvider struct {
	pulumi.ProviderResourceState // Embedding ProviderResourceState struct to capture internal state.
	PulumiCtx                    *pulumi.Context
}

// KindCluster represents a specific instance of a Kind cluster.
//
// It holds the state of a created or to-be-created Kind cluster.
type KindCluster struct {
	pulumi.ResourceState // Embedding ResourceState for Pulumi's internal use.

	ClusterName pulumi.StringOutput `pulumi:"clusterName"` // The name of the Kind cluster.
	ConfigFile  pulumi.StringOutput `pulumi:"configFile"`  // The configuration file used for the Kind cluster.
}

// KindClusterArgs holds the arguments for creating a new Kind cluster.
//
// These arguments are passed in while creating a KindCluster resource.
type KindClusterArgs struct {
	ClusterName   pulumi.String `pulumi:"clusterName,optional"`   // The name of the Kind cluster. Optional.
	ConfigFile    pulumi.String `pulumi:"configFile,optional"`    // The configuration file for the Kind cluster. Optional.
	KindConfigDir pulumi.String `pulumi:"kindConfigDir,optional"` // The directory containing Kind configuration files. Optional.
	Purge         pulumi.Bool   `pulumi:"purge,optional"`         // Whether to purge existing settings. Optional.
}

// VolumeArgs holds the arguments for creating Docker volumes for a Kind cluster.
//
// It contains the cluster name and a list of volume names.
type VolumeArgs struct {
	ClusterName string   `pulumi:"clusterName"` // The name of the Kind cluster to which the volumes will be attached.
	VolumeNames []string `pulumi:"volumeNames"` // The list of Docker volume names.
}

// ClusterUtility defines the interface for utility functions to manage Kind clusters.
//
// It includes methods for creating, deleting, and updating Kind clusters.
type ClusterUtility interface {
	Create(ctx *pulumi.Context, args KindClusterArgs, preview bool) error // Creates a new Kind cluster.
	Update(ctx *pulumi.Context, args KindClusterArgs, preview bool) error // Deletes an existing Kind cluster.
	Delete(ctx *pulumi.Context, clusterName string) error                 // Updates an existing Kind cluster.
}

// KindClusterUtility is an empty struct that implements the ClusterUtility interface.
//
// This is used to provide utility methods for managing Kind clusters.
type KindClusterUtility struct{}
