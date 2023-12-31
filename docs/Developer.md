# About

This repository is designed with [VS Code](https://code.visualstudio.com) integration to automate dependency and prerequisites as much as possible using [Dev Containers](https://containers.dev/)) to prepare your development environment, or even just run your development directly in the browser with [Github CodeSpaces](https://github.com/features/codespaces).

![CodeSpaces Screenshot](./assets/codespaces.png)

# Getting Started

To start, open this repository in Github CodeSpaces, or clone the repo locally and launch with VS Code using the devcontainer.

> Fig 1. How to open project in CodeSpaces
![How to open repository in CodeSpaces](./assets/gh-open-codespaces.png)

# Workflow

## Build & Install

Build and install this provider in your dev container.

```bash
make build install
```

> Fig 2.a Command `make build install` in terminal
![command: make build install](./assets/make-build-install.png)
![command complete: make build install](./assets/make-build-install-complete.png)

## Test

1. Pulumi Login

```bash
pulumi login
```

> Fig 2.b pulumi login
![Pulumi login](./assets/pulumi-login.png)
![Pulumi login complete](./assets/pulumi-login-complete.png)

2. Create a new Pulumi 'stack'

```bash
cd examples/yaml
pulumi stack init kind
pulumi stack select kind
```

> Fig 2.c pulumi stack
![pulumi stack](./assets/pulumi-stack.png)

3. Set pulumi config key values

```bash
pulumi config set name kind-provider-dev
```

4. Run `pulumi up`

```bash
pulumi up
```

> Fig 2.d pulumi up
![pulumi up](./assets/pulumi-up.png)

# Dependencies

If you are not using the included Def Container to meet dependencies and prerequisites, then make sure to satisfy the following requirements and ensure they are present in your `$PATH`:

* [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
* [Go 1.21](https://golang.org/dl/) or 1.latest
* [NodeJS](https://nodejs.org/en/) 14.x.  We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
* [Yarn](https://yarnpkg.com/)
* [TypeScript](https://www.typescriptlang.org/)
* [Python](https://www.python.org/downloads/) (called as `python3`).  For recent versions of MacOS, the system-installed version is fine.
* [.NET](https://dotnet.microsoft.com/download)
* [Docker]()
* [KinD]()

# Release

1. Push a tag to this repo in the format "v0.0.0" to initiate a release

1. IMPORTANT: also add a tag in the format "sdk/v0.0.0" for the Go SDK
