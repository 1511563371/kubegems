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

package statistics

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	gemsv1beta1 "kubegems.io/kubegems/pkg/apis/gems/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const ResourceLimitsPrefix = "limits."

type ClusterResourceStatistics struct {
	// 集群资源的总容量，即物理资源总量
	Capacity corev1.ResourceList `json:"capacity"`
	// 集群资源的真实使用量
	Used corev1.ResourceList `json:"used"`
	// 集群下的租户资源分配总量
	TenantAllocated corev1.ResourceList `json:"tenantAllocated"`
}

type ClusterWorkloadStatistics map[string]int

func GetClusterResourceStatistics(ctx context.Context, cli client.Client) ClusterResourceStatistics {
	nodelist := &corev1.NodeList{}
	_ = cli.List(ctx, nodelist)

	allcapacity := corev1.ResourceList{}
	allfree := corev1.ResourceList{}
	// all node capacity and free
	for _, node := range nodelist.Items {
		AddResourceList(allcapacity, node.Status.Capacity)
		AddResourceList(allfree, node.Status.Allocatable)
	}
	// calculate used
	allused := allcapacity.DeepCopy()
	SubResourceList(allused, allfree)

	allTenantAllocated, _ := GetClusterTenantResourceQuotaLimits(ctx, cli)

	// back compat
	if val, ok := allTenantAllocated["requests.storage"]; ok {
		allTenantAllocated["limits.storage"] = val.DeepCopy()
	}

	allTenantAllocated = FilterResourceName(allTenantAllocated, func(name corev1.ResourceName) bool {
		return strings.HasPrefix(string(name), ResourceLimitsPrefix)
	})

	// tenaut allocated has resource name with limit. prefix
	// we add limit. prefix for capacity and used
	allcapacity = AppendResourceNamePrefix(ResourceLimitsPrefix, allcapacity)
	allused = AppendResourceNamePrefix(ResourceLimitsPrefix, allused)

	return ClusterResourceStatistics{
		Capacity:        allcapacity,
		Used:            allused,
		TenantAllocated: allTenantAllocated,
	}
}

func FilterResourceName(list corev1.ResourceList, keep func(name corev1.ResourceName) bool) corev1.ResourceList {
	ret := corev1.ResourceList{}
	for k, v := range list {
		if keep(k) {
			ret[k] = v.DeepCopy()
		}
	}
	return ret
}

func AppendResourceNamePrefix(prefix string, list corev1.ResourceList) corev1.ResourceList {
	ret := corev1.ResourceList{}
	for k, v := range list {
		ret[corev1.ResourceName(prefix)+k] = v.DeepCopy()
	}
	return ret
}

func GetClusterTenantResourceQuotaLimits(ctx context.Context, cli client.Client) (corev1.ResourceList, error) {
	tenantResourceQuotaList := &gemsv1beta1.TenantResourceQuotaList{}
	if err := cli.List(ctx, tenantResourceQuotaList); err != nil {
		return nil, err
	}
	total := corev1.ResourceList{}
	for _, tquota := range tenantResourceQuotaList.Items {
		AddResourceList(total, tquota.Spec.Hard)
	}
	return total, nil
}

type ClusterPodResourceStatistics struct {
	Limit   corev1.ResourceList `json:"limit"`
	Request corev1.ResourceList `json:"request"`
}

func GetAllPodResourceStatistics(ctx context.Context, cli client.Client) (ClusterPodResourceStatistics, error) {
	podList := &corev1.PodList{}
	if err := cli.List(ctx, podList); err != nil {
		return ClusterPodResourceStatistics{}, err
	}
	limitResource := corev1.ResourceList{}
	requestResource := corev1.ResourceList{}
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		for _, container := range pod.Spec.Containers {
			AddResourceList(limitResource, container.Resources.Limits)
			AddResourceList(requestResource, container.Resources.Requests)
		}
	}
	return ClusterPodResourceStatistics{Limit: limitResource, Request: requestResource}, nil
}

func AddResourceList(total corev1.ResourceList, add corev1.ResourceList) {
	ResourceListCollect(total, add, func(_ corev1.ResourceName, into *resource.Quantity, val resource.Quantity) {
		into.Add(val)
	})
}

func SubResourceList(total corev1.ResourceList, sub corev1.ResourceList) {
	ResourceListCollect(total, sub, func(_ corev1.ResourceName, into *resource.Quantity, val resource.Quantity) {
		into.Sub(val)
	})
}

type ResourceListCollectFunc func(corev1.ResourceName, *resource.Quantity, resource.Quantity)

func ResourceListCollect(into, vals corev1.ResourceList, collect ResourceListCollectFunc) corev1.ResourceList {
	for resourceName, quantity := range vals {
		lastQuantity := into[resourceName].DeepCopy()
		collect(resourceName, &lastQuantity, quantity)
		into[resourceName] = lastQuantity
	}
	return into
}
