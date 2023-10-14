// File: provider/types.go

package provider

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// KindClusterProvider struct represents the provider itself.
type KindClusterProvider struct {
	pulumi.ProviderResourceState
}

// KindCluster struct represents a Kind cluster.
type KindCluster struct {
	pulumi.ResourceState

	ClusterName pulumi.StringOutput `pulumi:"clusterName"`
	ConfigFile  pulumi.StringOutput `pulumi:"configFile"`
}

// KindClusterArgs struct holds the input arguments to create a KindCluster.
type KindClusterArgs struct {
	ClusterName pulumi.String `pulumi:"clusterName"`
	ConfigFile  pulumi.String `pulumi:"configFile"`
	WorkingDir  pulumi.String `pulumi:"workingDir"`
	Purge       pulumi.Bool   `pulumi:"purge"`
}

// VolumeProvider struct represents the resource provider for Docker volumes.
type VolumeProvider struct {
}

// VolumeArgs struct holds the input arguments for creating a Docker volume.
type VolumeArgs struct {
	ClusterName string   `pulumi:"clusterName"`
	VolumeNames []string `pulumi:"volumeNames"`
}
