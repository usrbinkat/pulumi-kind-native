// File: provider/kind/types.go
package kind

import (
	"reflect"
)

type Kind struct {
	Name string `pulumi:"name"`
}

type KindStateArgs struct {
	Name string `pulumi:"name"`
}

func (k *Kind) ElementType() reflect.Type {
	return reflect.TypeOf((*Kind)(nil)).Elem()
}

func (s *KindStateArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*KindStateArgs)(nil)).Elem()
}
