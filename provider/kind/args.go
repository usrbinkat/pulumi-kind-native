// File: provider/kind/args.go
package kind

import "github.com/pulumi/pulumi-go-provider/infer"

var _ = (infer.Annotated)((*KindArgs)(nil))

func (ka *KindArgs) Annotate(a infer.Annotator) {
	a.Describe(&ka.Name, "The name of the KinD cluster.")
}
