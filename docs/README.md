# Create a new Provider

```bash
mkdir provider && cd provider
go mod init pulumi-kind-provider

mkdir kind
touch kind/kind.go kind/args.go
````

# Deploy

1. Customize the .github/workflows/release.yaml with the correct tokens using the format:

      `${{ secrets.MyTokenName }}`

1. Push a tag to your repo in the format "v0.0.0" to initiate a release

1. IMPORTANT: also add a tag in the format "sdk/v0.0.0" for the Go SDK
