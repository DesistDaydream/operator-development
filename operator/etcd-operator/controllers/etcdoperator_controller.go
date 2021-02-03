/*
Copyright 2020 DesistDaydream.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	testv1alpha1 "etcd-operator/api/v1alpha1"
)

// EtcdOperatorReconciler reconciles a EtcdOperator object
type EtcdOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=test.desistdaydream.ltd,resources=etcdoperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=test.desistdaydream.ltd,resources=etcdoperators/status,verbs=get;update;patch

// Reconcile 用于调谐对象状态
func (r *EtcdOperatorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("EtcdOperator", req.NamespacedName)

	// 首先我们获取 EtcdOperator 实例
	var EtcdOperator testv1alpha1.EtcdOperator
	if err := r.Client.Get(ctx, req.NamespacedName, &EtcdOperator); err != nil {
		// EtcdOperator was deleted，Ignore
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 得到 EtcdOperator 过后去创建对应的StatefulSet和Service
	// (就是观察的当前状态和期望的状态进行对比)
	// 调谐，获取到当前的一个状态，然后和我们期望的状态进行对比是不是就可以
	var svc corev1.Service
	svc.Name = EtcdOperator.Name
	svc.Namespace = EtcdOperator.Namespace
	or, err := ctrl.CreateOrUpdate(ctx, r, &svc, func() error {
		// 调谐必须在这个函数中去实现
		MutateHeadlessSvc(&EtcdOperator, &svc)
		return controllerutil.SetControllerReference(&EtcdOperator, &svc, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("CreateOrUpdate", "Service", or)

	// CreateOrUpdate StatefulSet
	var sts appsv1.StatefulSet
	sts.Name = EtcdOperator.Name
	sts.Namespace = EtcdOperator.Namespace
	or, err = ctrl.CreateOrUpdate(ctx, r, &sts, func() error {
		// 调谐必须在这个函数中去实现
		MutateStatefulSet(&EtcdOperator, &sts)
		return controllerutil.SetControllerReference(&EtcdOperator, &sts, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("CreateOrUpdate", "StatefulSet", or)

	// _ = context.Background()
	// _ = r.Log.WithValues("etcdoperator", req.NamespacedName)

	// // your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager is
func (r *EtcdOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.EtcdOperator{}).
		Complete(r)
}
