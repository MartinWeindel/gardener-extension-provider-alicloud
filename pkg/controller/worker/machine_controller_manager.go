// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"context"
	"fmt"
	"path/filepath"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/utils/chart"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/gardener/gardener-extension-provider-alicloud/charts"
	"github.com/gardener/gardener-extension-provider-alicloud/pkg/alicloud"
)

var (
	mcmChart = &chart.Chart{
		Name:       alicloud.MachineControllerManagerName,
		EmbeddedFS: &charts.InternalChart,
		Path:       filepath.Join(charts.InternalChartsPath, alicloud.MachineControllerManagerName, "seed"),
		Images:     []string{alicloud.MachineControllerManagerImageName, alicloud.MachineControllerManagerProviderAlicloudImageName},
		Objects: []*chart.Object{
			{Type: &appsv1.Deployment{}, Name: alicloud.MachineControllerManagerName},
			{Type: &corev1.Service{}, Name: alicloud.MachineControllerManagerName},
			{Type: &corev1.ServiceAccount{}, Name: alicloud.MachineControllerManagerName},
			{Type: &corev1.Secret{}, Name: alicloud.MachineControllerManagerName},
			{Type: extensionscontroller.GetVerticalPodAutoscalerObject(), Name: alicloud.MachineControllerManagerVpaName},
			{Type: &corev1.ConfigMap{}, Name: alicloud.MachineControllerManagerMonitoringConfigName},
		},
	}

	mcmShootChart = &chart.Chart{
		Name:       alicloud.MachineControllerManagerName,
		EmbeddedFS: &charts.InternalChart,
		Path:       filepath.Join(charts.InternalChartsPath, alicloud.MachineControllerManagerName, "shoot"),
		Objects: []*chart.Object{
			{Type: &rbacv1.ClusterRole{}, Name: fmt.Sprintf("extensions.gardener.cloud:%s:%s", alicloud.Name, alicloud.MachineControllerManagerName)},
			{Type: &rbacv1.ClusterRoleBinding{}, Name: fmt.Sprintf("extensions.gardener.cloud:%s:%s", alicloud.Name, alicloud.MachineControllerManagerName)},
		},
	}
)

func (w *workerDelegate) GetMachineControllerManagerChartValues(ctx context.Context) (map[string]interface{}, error) {
	namespace := &corev1.Namespace{}
	if err := w.client.Get(ctx, kutil.Key(w.worker.Namespace), namespace); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"providerName": alicloud.Name,
		"namespace": map[string]interface{}{
			"uid": namespace.UID,
		},
		"podLabels": map[string]interface{}{
			v1beta1constants.LabelPodMaintenanceRestart: "true",
		},
	}, nil
}

func (w *workerDelegate) GetMachineControllerManagerShootChartValues(_ context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"providerName": alicloud.Name,
	}, nil
}
