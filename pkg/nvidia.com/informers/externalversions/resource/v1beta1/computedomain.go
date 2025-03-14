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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	resourcev1beta1 "github.com/NVIDIA/k8s-dra-driver-gpu/api/nvidia.com/resource/v1beta1"
	versioned "github.com/NVIDIA/k8s-dra-driver-gpu/pkg/nvidia.com/clientset/versioned"
	internalinterfaces "github.com/NVIDIA/k8s-dra-driver-gpu/pkg/nvidia.com/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/NVIDIA/k8s-dra-driver-gpu/pkg/nvidia.com/listers/resource/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ComputeDomainInformer provides access to a shared informer and lister for
// ComputeDomains.
type ComputeDomainInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.ComputeDomainLister
}

type computeDomainInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewComputeDomainInformer constructs a new informer for ComputeDomain type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewComputeDomainInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredComputeDomainInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredComputeDomainInformer constructs a new informer for ComputeDomain type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredComputeDomainInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ResourceV1beta1().ComputeDomains(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ResourceV1beta1().ComputeDomains(namespace).Watch(context.TODO(), options)
			},
		},
		&resourcev1beta1.ComputeDomain{},
		resyncPeriod,
		indexers,
	)
}

func (f *computeDomainInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredComputeDomainInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *computeDomainInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&resourcev1beta1.ComputeDomain{}, f.defaultInformer)
}

func (f *computeDomainInformer) Lister() v1beta1.ComputeDomainLister {
	return v1beta1.NewComputeDomainLister(f.Informer().GetIndexer())
}
