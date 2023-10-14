// File: provider/utils.go
package provider

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// Initialize performs initial setup and dependency checks.
// Accepts a slice of dependencies that are needed.
func Initialize(ctx *pulumi.Context, dependencies []string) error {
	ctx.Log.Info("Starting dependency checks...", nil)

	var wg sync.WaitGroup
	errChan := make(chan error, len(dependencies))

	// Perform dependency checks concurrently
	for _, dependency := range dependencies {
		wg.Add(1)
		go func(dep string) {
			defer wg.Done()
			if err := CheckDependencies(ctx, dep); err != nil {
				errChan <- fmt.Errorf("Dependency check failed for %s: %w", dep, err)
			}
		}(dependency)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check if any errors occurred during the dependency checks
	var errs []string
	for err := range errChan {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multiple errors: %s", strings.Join(errs, ", "))
	}

	ctx.Log.Info("Local dependencies satisfied...", nil)
	return nil
}

// LoadKindConfig loads kind-specific configurations from Pulumi's config context.
func LoadKindConfig(ctx *pulumi.Context) (*KindClusterArgs, error) {
	cfg := config.New(ctx, "kind")

	purge := cfg.GetBool("purge")
	workingDir := cfg.Get("workingDir")
	if workingDir == "" {
		workingDir = "./"
	}

	// Resolve workingDir to an absolute path
	absWorkingDir, err := filepath.Abs(workingDir)
	if err != nil {
		return nil, err
	}

	clusterName := cfg.Get("kindClusterName")
	if clusterName == "" {
		clusterName = "cluster"
	}

	configFile := cfg.Get("kindConfigFile")
	if configFile == "" {
		configFile = "kind.yaml"
	}

	// Join workingDir and configFile
	configFilePath := filepath.Join(absWorkingDir, configFile)

	return &KindClusterArgs{
		ClusterName: pulumi.String(clusterName),
		ConfigFile:  pulumi.String(configFilePath),
		WorkingDir:  pulumi.String(absWorkingDir),
		Purge:       pulumi.Bool(purge),
	}, nil
}

// CheckDependencies checks if a specific dependency is installed.
func CheckDependencies(ctx *pulumi.Context, dependency string) error {
	if _, err := exec.LookPath(dependency); err != nil {
		return fmt.Errorf("Dependency '%s' is not installed: %w", dependency, err)
	}
	return nil
}

// ValidateClusterName validates a cluster name
func ValidateClusterName(name string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	if !re.MatchString(name) {
		return fmt.Errorf("Invalid cluster name: %s. Only alphanumeric and hyphens are allowed", name)
	}
	return nil
}
