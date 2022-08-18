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

package observability

import (
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"kubegems.io/kubegems/pkg/utils"
)

func Test_newDefaultSamplePair(t *testing.T) {
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want []model.SamplePair
	}{
		{
			name: "1",
			args: args{
				start: utils.DayStartTime(time.Now()),
				end:   utils.NextDayStartTime(time.Now()),
			},
			want: []model.SamplePair{
				{
					Timestamp: model.Time(utils.DayStartTime(time.Now()).UnixMilli()),
					Value:     0,
				},
				{
					Timestamp: model.Time(utils.NextDayStartTime(time.Now()).UnixMilli()),
					Value:     0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDefaultSamplePair(tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDefaultSamplePair() = %v, want %v", got, tt.want)
			}
		})
	}
}
