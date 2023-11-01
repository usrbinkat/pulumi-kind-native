# Pulumi Native Provider Layout

This repository contains the Pulumi Native Provider for KinD ([Kubernetes-in-Docker](https://kind.sigs.k8s.io)). The provider allows you to manage local Kind kubernetes clusters using Pulumi.

## Repository Layout

The repository is organized as follows:

- `sdk/`: Contains the generated multi-language provider SDKs.
- `provider/`: Contains the implementation of the provider itself.
- `examples/`: Contains example Pulumi programs that use the provider.
- `tests/`: Contains the tests for the provider.

## Getting Started

To use the provider, you must have Pulumi installed. You can install it using [these instructions](https://www.pulumi.com/docs/get-started/install/).

Once you have Pulumi installed, you can use the provider in your Pulumi programs by adding the following import statement:
