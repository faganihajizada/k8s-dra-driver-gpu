/*
 * Copyright (c) 2023, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/NVIDIA/k8s-dra-driver-gpu/api/nvidia.com/resource/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeComputeDomains implements ComputeDomainInterface
type FakeComputeDomains struct {
	Fake *FakeResourceV1beta1
	ns   string
}

var computedomainsResource = v1beta1.SchemeGroupVersion.WithResource("computedomains")

var computedomainsKind = v1beta1.SchemeGroupVersion.WithKind("ComputeDomain")

// Get takes name of the computeDomain, and returns the corresponding computeDomain object, and an error if there is any.
func (c *FakeComputeDomains) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.ComputeDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(computedomainsResource, c.ns, name), &v1beta1.ComputeDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ComputeDomain), err
}

// List takes label and field selectors, and returns the list of ComputeDomains that match those selectors.
func (c *FakeComputeDomains) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.ComputeDomainList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(computedomainsResource, computedomainsKind, c.ns, opts), &v1beta1.ComputeDomainList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.ComputeDomainList{ListMeta: obj.(*v1beta1.ComputeDomainList).ListMeta}
	for _, item := range obj.(*v1beta1.ComputeDomainList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested computeDomains.
func (c *FakeComputeDomains) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(computedomainsResource, c.ns, opts))

}

// Create takes the representation of a computeDomain and creates it.  Returns the server's representation of the computeDomain, and an error, if there is any.
func (c *FakeComputeDomains) Create(ctx context.Context, computeDomain *v1beta1.ComputeDomain, opts v1.CreateOptions) (result *v1beta1.ComputeDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(computedomainsResource, c.ns, computeDomain), &v1beta1.ComputeDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ComputeDomain), err
}

// Update takes the representation of a computeDomain and updates it. Returns the server's representation of the computeDomain, and an error, if there is any.
func (c *FakeComputeDomains) Update(ctx context.Context, computeDomain *v1beta1.ComputeDomain, opts v1.UpdateOptions) (result *v1beta1.ComputeDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(computedomainsResource, c.ns, computeDomain), &v1beta1.ComputeDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ComputeDomain), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeComputeDomains) UpdateStatus(ctx context.Context, computeDomain *v1beta1.ComputeDomain, opts v1.UpdateOptions) (*v1beta1.ComputeDomain, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(computedomainsResource, "status", c.ns, computeDomain), &v1beta1.ComputeDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ComputeDomain), err
}

// Delete takes name of the computeDomain and deletes it. Returns an error if one occurs.
func (c *FakeComputeDomains) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(computedomainsResource, c.ns, name, opts), &v1beta1.ComputeDomain{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeComputeDomains) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(computedomainsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.ComputeDomainList{})
	return err
}

// Patch applies the patch and returns the patched computeDomain.
func (c *FakeComputeDomains) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ComputeDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(computedomainsResource, c.ns, name, pt, data, subresources...), &v1beta1.ComputeDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ComputeDomain), err
}
