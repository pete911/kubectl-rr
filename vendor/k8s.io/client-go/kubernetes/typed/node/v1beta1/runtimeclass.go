/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	context "context"

	nodev1beta1 "k8s.io/api/node/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsnodev1beta1 "k8s.io/client-go/applyconfigurations/node/v1beta1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// RuntimeClassesGetter has a method to return a RuntimeClassInterface.
// A group's client should implement this interface.
type RuntimeClassesGetter interface {
	RuntimeClasses() RuntimeClassInterface
}

// RuntimeClassInterface has methods to work with RuntimeClass resources.
type RuntimeClassInterface interface {
	Create(ctx context.Context, runtimeClass *nodev1beta1.RuntimeClass, opts v1.CreateOptions) (*nodev1beta1.RuntimeClass, error)
	Update(ctx context.Context, runtimeClass *nodev1beta1.RuntimeClass, opts v1.UpdateOptions) (*nodev1beta1.RuntimeClass, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*nodev1beta1.RuntimeClass, error)
	List(ctx context.Context, opts v1.ListOptions) (*nodev1beta1.RuntimeClassList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *nodev1beta1.RuntimeClass, err error)
	Apply(ctx context.Context, runtimeClass *applyconfigurationsnodev1beta1.RuntimeClassApplyConfiguration, opts v1.ApplyOptions) (result *nodev1beta1.RuntimeClass, err error)
	RuntimeClassExpansion
}

// runtimeClasses implements RuntimeClassInterface
type runtimeClasses struct {
	*gentype.ClientWithListAndApply[*nodev1beta1.RuntimeClass, *nodev1beta1.RuntimeClassList, *applyconfigurationsnodev1beta1.RuntimeClassApplyConfiguration]
}

// newRuntimeClasses returns a RuntimeClasses
func newRuntimeClasses(c *NodeV1beta1Client) *runtimeClasses {
	return &runtimeClasses{
		gentype.NewClientWithListAndApply[*nodev1beta1.RuntimeClass, *nodev1beta1.RuntimeClassList, *applyconfigurationsnodev1beta1.RuntimeClassApplyConfiguration](
			"runtimeclasses",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *nodev1beta1.RuntimeClass { return &nodev1beta1.RuntimeClass{} },
			func() *nodev1beta1.RuntimeClassList { return &nodev1beta1.RuntimeClassList{} },
			gentype.PrefersProtobuf[*nodev1beta1.RuntimeClass](),
		),
	}
}