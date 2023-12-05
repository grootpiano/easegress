/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
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

// Package function provides FaasController.
package function

import (
	"fmt"
	"strings"

	"github.com/megaease/easegress/v2/pkg/api"
	"github.com/megaease/easegress/v2/pkg/object/function/spec"
	"github.com/megaease/easegress/v2/pkg/object/function/worker"
	"github.com/megaease/easegress/v2/pkg/supervisor"
	"github.com/megaease/easegress/v2/pkg/v"
)

const (
	// Category is the category of FaasController.
	Category = supervisor.CategoryBusinessController

	// Kind is the kind of FaaSController.
	Kind = "FaaSController"
)

var aliases = []string{
	"faas",
}

func init() {
	supervisor.Register(&FaasController{})
	api.RegisterObject(&api.APIResource{
		Category: Category,
		Kind:     Kind,
		Name:     strings.ToLower(Kind),
		Aliases:  aliases,
	})
}

type (
	// FaasController is Function controller.
	FaasController struct {
		superSpec *supervisor.Spec
		spec      *spec.Admin

		worker *worker.Worker
	}
)

// Category returns the category of FaasController.
func (f *FaasController) Category() supervisor.ObjectCategory {
	return Category
}

// Kind returns the kind of FaasController.
func (f *FaasController) Kind() string {
	return Kind
}

// DefaultSpec returns the default spec of Function.
func (f *FaasController) DefaultSpec() interface{} {
	return &spec.Admin{
		SyncInterval: "10s",
		Provider:     spec.ProviderKnative,
		Knative: &spec.Knative{
			Namespace: "default",
			Timeout:   "2s",
		},
	}
}

// Validate validates the spec
func (f *FaasController) Validate() error {
	switch f.spec.Provider {
	case spec.ProviderKnative:
		//
	default:
		return fmt.Errorf("unknown FaaS provider: %s", f.spec.Provider)
	}

	vr := v.Validate(f.spec.HTTPServer)
	if !vr.Valid() {
		return fmt.Errorf("%s", vr.Error())
	}
	return nil
}

// Init initializes Function.
func (f *FaasController) Init(superSpec *supervisor.Spec) {
	f.superSpec, f.spec = superSpec, superSpec.ObjectSpec().(*spec.Admin)
	f.reload()
}

// Inherit inherits previous generation of Function.
func (f *FaasController) Inherit(superSpec *supervisor.Spec, previousGeneration supervisor.Object) {
	previousGeneration.Close()
	f.Init(superSpec)
}

func (f *FaasController) reload() {
	f.worker = worker.NewWorker(f.superSpec)
}

// Status returns Status generated by Runtime.
func (f *FaasController) Status() *supervisor.Status {
	return &supervisor.Status{}
}

// Close closes Function.
func (f *FaasController) Close() {
	f.worker.Close()
}
