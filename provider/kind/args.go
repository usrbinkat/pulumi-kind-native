// File: provider/kind/args.go
package kind

import "github.com/pulumi/pulumi-go-provider/infer"

var _ = (infer.Annotated)((*KindStateArgs)(nil))

func (ka *KindStateArgs) Annotate(a infer.Annotator) {
	a.Describe(&ka.Name, "The name of the KinD cluster.")
}
