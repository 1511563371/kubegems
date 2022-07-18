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

package application

import "testing"

func TestPathRef_FullName(t *testing.T) {
	type fields struct {
		Tenant  string
		Project string
		Env     string
		Name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Tenant:  "t",
				Project: "p",
				Env:     "e",
				Name:    "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglong",
			},
			want: "5a638767-t-p-e-longlonglonglonglonglonglonglong",
		},
		{
			fields: fields{
				Tenant:  "tdasgoucgdsaiugyfiur3",
				Project: "fewklgviwhigver.-kfdkjsgv",
				Env:     "e321gdiuqwkfwelkf",
				Name:    "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglong",
			},
			want: "ef492a8e-tda-few-e32-longlonglonglonglonglonglonglong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PathRef{
				Tenant:  tt.fields.Tenant,
				Project: tt.fields.Project,
				Env:     tt.fields.Env,
				Name:    tt.fields.Name,
			}
			if got := p.FullName(); got != tt.want {
				t.Errorf("PathRef.FullName() = %v, want %v", got, tt.want)
			}
		})
	}
}
