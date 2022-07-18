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

const (
	LabelFrom              = GroupName + "/from" // 区分是从 appstore 还是从应用 app 部署的argo
	LabelValueFromApp      = "app"
	LabelValueFromAppStore = "appstore"

	AnnotationCreator   = GroupName + "/creator"   // 创建人,仅用于当前部署实时更新，从kustomize部署的历史需要从gitcommit取得
	AnnotationRef       = GroupName + "/ref"       // 标志这个资源所属的项目环境，避免使用过多label造成干扰
	AnnotationCluster   = GroupName + "/cluster"   // 标志这个资源所属集群
	AnnotationNamespace = GroupName + "/namespace" // 标志这个资源所属namespace

	AnnotationGeneratedByPlatformKey           = GroupName + "/is-generated-by" // 表示该资源是为了{for}自动生成的，需要在特定的时刻被清理
	AnnotationGeneratedByPlatformValueRollouts = "rollouts"                     // 表示该资源是为了rollouts而生成的

	AnnotationImagePullSecretKeyPrefix = GroupName + "/imagePullSecrets-"
)
