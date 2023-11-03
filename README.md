# Pulumi KinD (Kubernetes-in-Docker) Native Provider

[![test](https://github.com/usrbinkat/pulumi-kind-native/actions/workflows/test.yaml/badge.svg)](https://github.com/usrbinkat/pulumi-kind-native/actions/workflows/test.yaml)
![License](https://img.shields.io/github/license/usrbinkat/pulumi-kind-native)

### Overview

The \`pulumi-kind-native\` provider is a specialized package aimed at the seamless provisioning and management of [Kind](https://kind.sigs.k8s.io/) clusters via Pulumi. This provider allows you to create, update, and delete Kind clusters and associated Docker volumes programmatically using Pulumi's infrastructure-as-code capabilities. 

### Features

- Provision and manage Kind clusters
- Manage associated Docker volumes
- Support for input validation
- Idempotent operations
- Extensible and customizable
- Comprehensive logging and debugging support

### Table of Contents

- [Kubernetes-in-Docker | Pulumi Native Provider](#kubernetes-in-docker--pulumi-native-provider)
  - [Pulumi-Kind-Native Provider](#pulumi-kind-native-provider)
    - [Overview](#overview)
    - [Features](#features)
    - [Table of Contents](#table-of-contents)
    - [Installation](#installation)
    - [Getting Started](#getting-started)
    - [Usage](#usage)
    - [API Documentation](#api-documentation)
    - [Examples](#examples)
    - [Troubleshooting](#troubleshooting)
    - [Contributing](#contributing)
    - [Testing](#testing)
    - [License](#license)
  - [References](#references)

### Installation

```bash
# Install the package (Node.js example)
$ pulumi plugin install resource kind v{version}
```

### Getting Started

Before you begin, make sure you have [Pulumi](https://www.pulumi.com/docs/get-started/install/) and [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) installed.

```bash
# Create a new Pulumi project
$ pulumi new typescript
```

### Usage

Here's a TypeScript example that shows how to create a Kind cluster:

```typescript
import * as pulumi from '@pulumi/pulumi';
import * as kind from '@usrbinkat/pulumi-kind-native';

const cluster = new kind.KindCluster("my-cluster", {
    clusterName: "pulumi-cluster",
    configFile: "./kind-config.yaml",
    purge: false,
});
```

### API Documentation

The API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/usrbinkat/pulumi-kind-native).

### Examples

For more comprehensive examples, please refer to the \`examples/\` directory in this repository.

### Troubleshooting

Check the logs for any errors or inconsistencies. The provider is designed to output comprehensive logs to assist in debugging.

### Contributing

Contributions are highly encouraged! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project. Don't forget to check the [Code of Conduct](CODE_OF_CONDUCT.md).

### Testing

To run the tests, execute the following command:

```bash
$ make test
```

### License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

## References

Other resources/examples for implementing providers:
* [Pulumi Command provider](https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go)
* [Pulumi Go Provider repository](https://github.com/pulumi/pulumi-go-provider)
