package tests

import (
	"testing"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kind "github.com/usrbinkat/pulumi-kind-native/provider/kind"
)

func TestKindClusterLifecycle(t *testing.T) {
	prov := provider()

	name := "test-cluster"
	props := resource.PropertyMap{
		"name": resource.NewStringProperty(name),
	}

	// Create a KinD cluster
	createResp, err := prov.Create(p.CreateRequest{
		Urn:        urn("Kind"),
		Properties: props,
		Preview:    false,
	})

	require.NoError(t, err)
	assert.Equal(t, name, createResp.Properties["name"].StringValue())

	// Delete the KinD cluster
	_, err = prov.Delete(p.DeleteRequest{
		Urn:        urn("Kind"),
		Id:         createResp.Id,
		Properties: createResp.Properties,
	})

	require.NoError(t, err)
}

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("test:index:"+typ), "name")
}

// Create a test server.
func provider() integration.Server {
	return integration.NewServer(kind.Name, semver.MustParse("1.0.0"), kind.Provider())
}
