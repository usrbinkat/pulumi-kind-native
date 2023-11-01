# Pulumi Native Provider Layout

This repository contains the Pulumi Native Provider for KinD ([Kubernetes-in-Docker](https://kind.sigs.k8s.io)). The provider allows you to manage local Kind kubernetes clusters using Pulumi.

## Repository Layout

The repository is organized as follows:

### Provider artifacts

- `sdk/`: Contains the generated multi-language provider SDKs.
- `provider/`: Contains the implementation of the provider itself.
- `examples/`: Contains example Pulumi programs that use the provider.
- `tests/`: Contains the tests for the provider.
- `.github/workflows/`: Contains the GitHub Actions CI workflows.

### Developer tools

- `.devcontainer/`: Contains configuration for the VS Code development container.
- `.vscode/`: Contains helpful VS Code workspace settings.
- `docker/`: Contains the Dockerfile and build artifacts for the development container.
- `docs/`: Contains documentation for the provider.
- `Makefile`: Contains helpful commands for building and testing the provider.

## Additional Resources

### Concepts

- [Pulumi Package](https://www.pulumi.com/docs/using-pulumi/pulumi-packages/#pulumi-packages)
- [Pulumi Resource](https://www.pulumi.com/docs/concepts/resources/)
- [Pulumi Resource Provider](https://www.pulumi.com/docs/concepts/resources/providers/)
- [Pulumi Component Resource](https://www.pulumi.com/docs/concepts/resources/components/)

### Guides

- [How Pulumi Works](https://www.pulumi.com/docs/concepts/how-pulumi-works/#how-pulumi-works)
- [Author Pulumi Packages](https://www.pulumi.com/docs/using-pulumi/pulumi-packages/how-to-author/)
- [Publish Pulumi Packages](https://www.pulumi.com/docs/using-pulumi/pulumi-packages/how-to-author/)

### Additional

- [Pulumi Documentation](https://www.pulumi.com/docs/)
- [Pulumi Community Slack](https://slack.pulumi.com/)
