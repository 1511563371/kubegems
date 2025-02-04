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

package prometheus

import (
	"database/sql/driver"
	"encoding/json"
)

func (g *MonitorGraphs) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, g)
	}
	return nil
}

// 注意这里不是指针，下同
func (g MonitorGraphs) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g MonitorGraphs) GormDataType() string {
	return "json"
}

type MonitorGraphs []MetricGraph

type MetricGraph struct {
	// graph名
	Name string `json:"name"`
	// 查询目标
	*PromqlGenerator `json:"promqlGenerator"`
	Expr             string `json:"expr"`
	Unit             string `json:"unit"`
}
