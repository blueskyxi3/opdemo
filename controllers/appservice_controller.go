/*
Copyright 2021 blueskyxi3.

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
	"fmt"
	"github.com/blueskyxi3/opdemo/v2/resources"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appv1beta1 "github.com/blueskyxi3/opdemo/v2/api/v1beta1"
)

var (
	oldSpecAnnotation = "old/spec"
)

// AppServiceReconciler reconciles a AppService object
type AppServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.ydzs.io,resources=appservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.ydzs.io,resources=appservices/status,verbs=get;update;patch

func (r *AppServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	fmt.Println("start to Reconcile.......")
	ctx := context.Background()
	log := r.Log.WithValues("appservice", req.NamespacedName)

	// 业务逻辑实现
	// 获取 AppService 实例
	var appService appv1beta1.AppService
	err := r.Get(ctx, req.NamespacedName, &appService)
	if err != nil {
		// AppService 被删除的时候，忽略
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	log.Info("fetch appservice objects", "appservice", appService)

	// 得到AppService过后去创建对应的Deployment和Service
	// CreateOrUpdate Deployment
	// (就是观察的当前状态和期望的状态进行对比)

	// 调谐，获取到当前的一个状态，然后和我们期望的状态进行对比
	// CreateOrUpdate Deployment
	var deploy appsv1.Deployment
	deploy.Name = appService.Name
	deploy.Namespace = appService.Namespace
	or, err := ctrl.CreateOrUpdate(ctx, r, &deploy, func() error {
		fmt.Println("CreateOrUpdate Deployment!")
		// 调谐必需在这个函数中去实现
		resources.MutateDeployment(&appService, &deploy)
		return controllerutil.SetControllerReference(&appService, &deploy, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("CreateOrUpdate", "Deployment", or)

	// CreateOrUpdate Deployment
	var svc corev1.Service
	svc.Name = appService.Name
	svc.Namespace = appService.Namespace
	or, err = ctrl.CreateOrUpdate(ctx, r, &svc, func() error {
		fmt.Println("CreateOrUpdate Service!")
		// 调谐必需在这个函数中去实现
		resources.MutateService(&appService, &svc)
		return controllerutil.SetControllerReference(&appService, &svc, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("CreateOrUpdate", "Service", or)

	return ctrl.Result{}, nil
}

func (r *AppServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1beta1.AppService{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
