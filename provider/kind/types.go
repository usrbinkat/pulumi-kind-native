// File: provider/kind/types.go
package kind

import (
	"reflect"
)

type Kind struct {
	Name string `pulumi:"name"`
}

type KindState struct {
	Name string `pulumi:"name"`
}

type KindArgs struct {
	Name string `pulumi:"name"`
}

func (k *Kind) ElementType() reflect.Type {
	return reflect.TypeOf((*Kind)(nil)).Elem()
}

func (s *KindState) ElementType() reflect.Type {
	return reflect.TypeOf((*KindState)(nil)).Elem()
}

func (a *KindArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*KindArgs)(nil)).Elem()
}
