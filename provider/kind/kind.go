// File: provider/kind/kind.go
package kind

import (
	"os/exec"

	p "github.com/pulumi/pulumi-go-provider"
)

func (*Kind) Create(ctx p.Context, name string, input KindArgs, preview bool) (string, KindState, error) {
	// Logic to create a KinD cluster
	cmd := exec.Command("kind", "create", "cluster", "--name", input.Name)
	err := cmd.Run()
	if err != nil {
		return "", KindState{}, err
	}

	state := KindState(input)

	return name, state, nil
}

func (*Kind) Delete(ctx p.Context, id string, props KindState) error {
	// Logic to delete a KinD cluster
	cmd := exec.Command("kind", "delete", "cluster", "--name", props.Name)
	return cmd.Run()
}
