// File: provider/pkg/provider/kind/utils.go

// Package provider contains utilities for managing Kind clusters in Pulumi.
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
//
// It accepts a slice of dependencies that are needed and performs
// concurrent checks to ensure they are all installed.
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

// LoadKindConfig loads Kind-specific configurations from Pulumi's config context.
//
// It reads configurations like clusterName, configFile, etc., from Pulumi's context
// and returns a KindClusterArgs struct filled with these configurations.
func LoadKindConfig(ctx *pulumi.Context) (*KindClusterArgs, error) {
	cfg := config.New(ctx, "kind")

	// Load various settings from the Pulumi config
	purge := cfg.GetBool("purge")
	workingDir := cfg.Get("workingDir")
	if workingDir == "" {
		workingDir = "./"
	}

	// Set default kind cluster name if not supplied.
	clusterName := cfg.Get("kindClusterName")
	if clusterName == "" {
		clusterName = "pulumi"
	}

	// Set default kind config file if not supplied.
	configFile := cfg.Get("kindConfigFile")
	if configFile == "" {
		configFile = "kind/config.yaml"
	}

	// Resolve workingDir to an absolute path
	absKindConfigDir, err := filepath.Abs(workingDir + configFile)
	ctx.Log.Info("Selecting Kind config: "+absKindConfigDir, nil)
	if err != nil {
		return nil, err
	}

	// Join workingDir and configFile
	configFilePath := filepath.Join(absKindConfigDir, configFile)

	// Create and return KindClusterArgs based on loaded config
	return &KindClusterArgs{
		ClusterName:   pulumi.String(clusterName),
		ConfigFile:    pulumi.String(configFilePath),
		KindConfigDir: pulumi.String(absKindConfigDir),
		Purge:         pulumi.Bool(purge),
	}, nil
}

// CheckDependencies checks if a specific dependency is installed.
//
// The function checks whether the given dependency can be found in the system's PATH.
func CheckDependencies(ctx *pulumi.Context, dependency string) error {
	if _, err := exec.LookPath(dependency); err != nil {
		return fmt.Errorf("Dependency '%s' is not installed: %w", dependency, err)
	}
	return nil
}

// ValidateClusterName validates the provided cluster name against a regular expression.
//
// The cluster name should only contain alphanumeric characters and hyphens.
func ValidateClusterName(name string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	if !re.MatchString(name) {
		return fmt.Errorf("Invalid cluster name: %s. Only alphanumeric and hyphens are allowed", name)
	}
	return nil
}
