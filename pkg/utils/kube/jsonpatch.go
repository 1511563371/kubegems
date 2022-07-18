// Copyright 2022 The kubegems.io Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kube

import (
	"encoding/json"

	"github.com/mattbaird/jsonpatch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type JsonPatchType struct {
	From client.Object
}

// Type is the PatchType of the patch.
func (j *JsonPatchType) Type() types.PatchType {
	return types.JSONPatchType
}

// Data is the raw data representing the patch.
func (j *JsonPatchType) Data(obj client.Object) ([]byte, error) {
	// ignore reversion
	if meta, ok := obj.(metav1.Object); ok {
		meta.SetResourceVersion("")
	}
	if meta, ok := j.From.(metav1.Object); ok {
		meta.SetResourceVersion("")
	}
	if kinded, ok := obj.(schema.ObjectKind); ok {
		kinded.SetGroupVersionKind(j.From.GetObjectKind().GroupVersionKind())
	}
	return CreatePatchContent(j.From, obj)
}

func CreatePatchContent(origin, modified interface{}) ([]byte, error) {
	o, e := json.Marshal(origin)
	if e != nil {
		return nil, e
	}
	m, e := json.Marshal(modified)
	if e != nil {
		return nil, e
	}
	patches, e := jsonpatch.CreatePatch(o, m)
	if e != nil {
		return nil, e
	}
	out, e := json.Marshal(patches)
	if e != nil {
		return nil, e
	}
	return out, nil
}
