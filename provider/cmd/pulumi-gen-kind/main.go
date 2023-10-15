// File: provider/cmd/pulumi-gen-kind/main.go

// provider/cmd/pulumi-gen-kind/main.go is the main entry point for the schema generator for the pulumi-kind-native provider.
//
// This program generates JSON schema files for the custom Kind cluster provider.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	kind "github.com/usrbinkat/pulumi-kind-native/provider/pkg/provider/kind"
	providerVersion "github.com/usrbinkat/pulumi-kind-native/provider/pkg/version"
)

// MarshalIndent is a utility function for marshalling JSON data with indentation.
// Unlike the standard json.Marshal function, this doesn't escape HTML characters.
// It takes an "any" type and returns the pretty-printed JSON byte array and any error.
func MarshalIndent(v any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false) // Disable HTML escaping
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	b := buffer.Bytes()

	// Serialize and pretty-print JSON
	var buf bytes.Buffer
	prefix, indent := "", "    "
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// main is the entry point for the schema generator.
// It parses command-line flags and generates the schema.
func main() {
	// Define and print usage information for the command-line utility.
	flag.Usage = func() {
		const usageFormat = "Usage: %s <schema-file>"
		fmt.Fprintf(flag.CommandLine.Output(), usageFormat, os.Args[0])
		flag.PrintDefaults()
	}

	var version string // Version of the provider
	// Parse version flag. If not set, it defaults to the provider's current version.
	flag.StringVar(&version, "version", providerVersion.Version, "the provider version to record in the generated code")

	// Parse any additional flags
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return
	}
	// Generate schema for the custom Kind cluster provider
	s, err := kind.Schema(version)
	if err != nil {
		panic(err)
	}

	// Sort keys in the JSON schema
	var arg map[string]any
	err = json.Unmarshal([]byte(s), &arg)
	if err != nil {
		panic(err)
	}

	// Remove the version key to avoid redundancy
	delete(arg, "version")

	// Marshal JSON schema with indentation and without HTML escaping
	out, err := MarshalIndent(arg)
	if err != nil {
		panic(err)
	}

	// Write the generated schema to a file
	schemaPath := args[0]
	err = os.WriteFile(schemaPath, out, 0600)
	if err != nil {
		panic(err)
	}
}
