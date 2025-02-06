/*
 * Copyright (c) 2025 NVIDIA CORPORATION.  All rights reserved.
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

package main

import (
	"context"
	"fmt"
	"sync"

	resourceapi "k8s.io/api/resource/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	nvapi "github.com/NVIDIA/k8s-dra-driver-gpu/api/nvidia.com/resource/v1beta1"
)

type DeviceClassManager struct {
	config        *ManagerConfig
	waitGroup     sync.WaitGroup
	cancelContext context.CancelFunc

	factory  informers.SharedInformerFactory
	informer cache.SharedIndexInformer
}

func NewDeviceClassManager(config *ManagerConfig) *DeviceClassManager {
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      computeDomainLabelKey,
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}

	factory := informers.NewSharedInformerFactoryWithOptions(
		config.clientsets.Core,
		informerResyncPeriod,
		informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
			opts.LabelSelector = metav1.FormatLabelSelector(labelSelector)
		}),
	)

	informer := factory.Resource().V1beta1().DeviceClasses().Informer()

	m := &DeviceClassManager{
		config:   config,
		factory:  factory,
		informer: informer,
	}

	return m
}

func (m *DeviceClassManager) Start(ctx context.Context) (rerr error) {
	ctx, cancel := context.WithCancel(ctx)
	m.cancelContext = cancel

	defer func() {
		if rerr != nil {
			if err := m.Stop(); err != nil {
				klog.Errorf("error stopping DeviceClass manager: %v", err)
			}
		}
	}()

	if err := addComputeDomainLabelIndexer[*resourceapi.DeviceClass](m.informer); err != nil {
		return fmt.Errorf("error adding indexer for MulitNodeEnvironment label: %w", err)
	}

	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()
		m.factory.Start(ctx.Done())
	}()

	if !cache.WaitForCacheSync(ctx.Done(), m.informer.HasSynced) {
		return fmt.Errorf("informer cache sync for DeviceClass failed")
	}

	return nil
}

func (m *DeviceClassManager) Stop() error {
	m.cancelContext()
	m.waitGroup.Wait()
	return nil
}

func (m *DeviceClassManager) Create(ctx context.Context, name string, cd *nvapi.ComputeDomain) (*resourceapi.DeviceClass, error) {
	dcs, err := getByComputeDomainUID[*resourceapi.DeviceClass](ctx, m.informer, string(cd.UID))
	if err != nil {
		return nil, fmt.Errorf("error retrieving DeviceClass: %w", err)
	}
	if len(dcs) > 1 {
		return nil, fmt.Errorf("more than one DeviceClass found with same ComputeDomain UID")
	}
	if len(dcs) == 1 {
		return dcs[0], nil
	}

	deviceClass := &resourceapi.DeviceClass{
		ObjectMeta: metav1.ObjectMeta{
			Finalizers: []string{computeDomainFinalizer},
			Labels: map[string]string{
				computeDomainLabelKey: string(cd.UID),
			},
		},
		Spec: resourceapi.DeviceClassSpec{
			Selectors: []resourceapi.DeviceSelector{
				{
					CEL: &resourceapi.CELDeviceSelector{
						Expression: fmt.Sprintf("device.driver == '%s'", DriverName),
					},
				},
				{
					CEL: &resourceapi.CELDeviceSelector{
						Expression: fmt.Sprintf("device.attributes['%s'].type == 'imex-channel'", DriverName),
					},
				},
				{
					CEL: &resourceapi.CELDeviceSelector{
						Expression: fmt.Sprintf("device.attributes['%s'].domain == '%v'", DriverName, cd.UID),
					},
				},
			},
		},
	}

	if name == "" {
		deviceClass.GenerateName = cd.Name
	} else {
		deviceClass.Name = name
	}

	dc, err := m.config.clientsets.Core.ResourceV1beta1().DeviceClasses().Create(ctx, deviceClass, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating DeviceClass: %w", err)
	}

	return dc, nil
}

func (m *DeviceClassManager) Delete(ctx context.Context, cdUID string) error {
	dcs, err := getByComputeDomainUID[*resourceapi.DeviceClass](ctx, m.informer, cdUID)
	if err != nil {
		return fmt.Errorf("error retrieving DeviceClass: %w", err)
	}
	if len(dcs) > 1 {
		return fmt.Errorf("more than one DeviceClass found with same ComputeDomain UID")
	}
	if len(dcs) == 0 {
		return nil
	}

	dc := dcs[0]

	if dc.GetDeletionTimestamp() != nil {
		return nil
	}

	err = m.config.clientsets.Core.ResourceV1beta1().DeviceClasses().Delete(ctx, dc.Name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("erroring deleting DeviceClass: %w", err)
	}

	return nil
}

func (m *DeviceClassManager) RemoveFinalizer(ctx context.Context, cdUID string) error {
	dcs, err := getByComputeDomainUID[*resourceapi.DeviceClass](ctx, m.informer, cdUID)
	if err != nil {
		return fmt.Errorf("error retrieving DeviceClass: %w", err)
	}
	if len(dcs) > 1 {
		return fmt.Errorf("more than one DeviceClass found with same ComputeDomain UID")
	}
	if len(dcs) == 0 {
		return nil
	}

	dc := dcs[0]

	if dc.GetDeletionTimestamp() == nil {
		return fmt.Errorf("attempting to remove finalizer before DeviceClass marked for deletion")
	}

	newDC := dc.DeepCopy()
	newDC.Finalizers = []string{}
	for _, f := range dc.Finalizers {
		if f != computeDomainFinalizer {
			newDC.Finalizers = append(newDC.Finalizers, f)
		}
	}
	if len(dc.Finalizers) == len(newDC.Finalizers) {
		return nil
	}

	if _, err = m.config.clientsets.Core.ResourceV1beta1().DeviceClasses().Update(ctx, newDC, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("error updating DeviceClass: %w", err)
	}

	return nil
}
