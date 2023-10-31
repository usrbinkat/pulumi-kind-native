// File: provider/kind/kind.go
package kind

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"

	p "github.com/pulumi/pulumi-go-provider"
)

// Custom error to add context
type KindError struct {
	Msg string
	Err error
}

func (e *KindError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

// Define a constant for info-level diagnostic logs
const diagInfo = "info"

type Kind struct {
	Name string `pulumi:"name"`
}

type KindStateArgs struct {
	Name string `pulumi:"name"`
}

func (k *Kind) ElementType() reflect.Type {
	return reflect.TypeOf((*Kind)(nil)).Elem()
}

func (s *KindStateArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*KindStateArgs)(nil)).Elem()
}

func (*Kind) Create(ctx p.Context, name string, input KindStateArgs, preview bool) (string, KindStateArgs, error) {
	// Enhanced Logging
	ctx.Logf(diagInfo, "Creating Kind cluster with name: %s", input.Name)

	// Check if we're in the preview phase
	if preview {
		return name, input, nil
	}

	// Logic to create a KinD cluster
	cmd := exec.Command("kind", "create", "cluster", "--name", input.Name)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	// Enhanced Error Handling
	if err != nil {
		return "", KindStateArgs{}, &KindError{"Failed to create KinD cluster", err}
	}

	// For future: Populate KindStateArgs with additional state information if needed
	// Returning input for now
	return name, input, nil
}

func (*Kind) Delete(ctx p.Context, id string, props KindStateArgs) error {
	// Enhanced Logging
	ctx.Logf(diagInfo, "Deleting Kind cluster with name: %s", props.Name)

	// Logic to delete a KinD cluster
	cmd := exec.Command("kind", "delete", "cluster", "--name", props.Name)
	err := cmd.Run()

	// Enhanced Error Handling
	if err != nil {
		return &KindError{"Failed to delete KinD cluster", err}
	}

	return nil
}
