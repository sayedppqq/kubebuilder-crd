/*
Copyright 2023.

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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "k8s.io/api/apps/v1"
	sayedppqqdevv1 "sayedppqq.dev/mycrd/api/v1"
)

// CustomRReconciler reconciles a CustomR object
type CustomRReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sayedppqq.dev.sayedppqq.dev,resources=customrs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sayedppqq.dev.sayedppqq.dev,resources=customrs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sayedppqq.dev.sayedppqq.dev,resources=customrs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomR object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CustomRReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.WithValues("ReqName", req.Name, "ReqNamespace", req.Namespace)
	fmt.Println("Reconcile started")

	// TODO(user): your logic here

	// customR have all data of CustomR Resources
	var customR sayedppqqdevv1.CustomR
	if err := r.Get(ctx, req.NamespacedName, &customR); err != nil {
		fmt.Errorf("unable to fetch CustomR")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, nil
	}
	fmt.Println("CustomR fetched", req.NamespacedName)

	// deploymentInstance have all data of deployment in specific namespace and name under the controller
	var deploymentInstance appsv1.Deployment
	// making deployment name
	deploymentName := func() string {
		if customR.Spec.DeploymentName == "" {
			return customR.Name + "-" + "randomName"
		} else {
			return customR.Name + "-" + customR.Spec.DeploymentName
		}
	}()
	//Creating NamespacedName for deploymentInstance
	obk := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      deploymentName,
	}

	if err := r.Get(ctx, obk, &deploymentInstance); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Deployment for ", customR.Name, ", creating one...")
			err := r.Client.Create(ctx, newDeployment(&customR, deploymentName))
			if err != nil {
				fmt.Errorf("error while creating deployment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Deployments Created...\n", customR.Name)
			}
		} else {
			fmt.Errorf("error fetching deployment %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if customR.Spec.Replicas != nil && *customR.Spec.Replicas != *deploymentInstance.Spec.Replicas {
			fmt.Println(*customR.Spec.Replicas, *deploymentInstance.Spec.Replicas)
			fmt.Println("Deployment replica miss match.....updating")
			//As the replica count didn't match, we need to update it
			deploymentInstance.Spec.Replicas = customR.Spec.Replicas
			if err := r.Update(ctx, &deploymentInstance); err != nil {
				fmt.Errorf("error updating deployment %s\n", err)
				return ctrl.Result{}, err
			}
			fmt.Println("deployment updated")
		}
	}

	var serviceInstance corev1.Service
	// making service name
	serviceName := func() string {
		if customR.Spec.Service.ServiceName == "" {
			return customR.Name + "-" + "randomName-service"
		} else {
			return customR.Name + "-" + customR.Spec.Service.ServiceName + "-service"
		}
	}()
	//Creating NamespacedName for serviceInstance
	obk = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      serviceName,
	}
	if err := r.Get(ctx, obk, &serviceInstance); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Service for ", customR.Name, ", creating one...")
			err := r.Create(ctx, newService(&customR, serviceName))
			if err != nil {
				fmt.Errorf("error while creating service %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Service Created...\n", customR.Name)
			}
		} else {
			fmt.Errorf("error fetching service %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if customR.Spec.Replicas != nil && *customR.Spec.Replicas != customR.Status.AvailableReplicas {
			fmt.Println("Service replica miss match.....updating")
			customR.Status.AvailableReplicas = *customR.Spec.Replicas
			if err := r.Status().Update(ctx, &customR); err != nil {
				fmt.Errorf("error while updating service %s\n", err)
				return ctrl.Result{}, err
			}
			fmt.Println("service updated")
		}
	}
	fmt.Println("reconciliation done")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomRReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// Main part
	//return ctrl.NewControllerManagedBy(mgr).
	//	For(&sayedppqqdevv1.CustomR{}).
	//	Owns(&appsv1.Deployment{}).
	//	Owns(&corev1.Service{}).
	//	Complete(r)

	// Extra part, if want to implement with watches and custom eventHandler
	// if someone edit the resources(here example given for deployment resource) by kubectl
	handlerForDeployment := handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
		// List all the CR
		customRs := &sayedppqqdevv1.CustomRList{}
		if err := r.List(context.Background(), customRs); err != nil {
			return nil
		}
		// This func return a reconcile request array
		var req []reconcile.Request
		for _, c := range customRs.Items {
			// Find the deployment owned by the CR
			if c.Spec.DeploymentName == obj.GetName() && c.Namespace == obj.GetNamespace() {
				deploy := &appsv1.Deployment{}
				if err := r.Get(context.Background(), types.NamespacedName{
					Namespace: obj.GetNamespace(),
					Name:      obj.GetName(),
				}, deploy); err != nil {
					return nil
				}
				// Only append to the reconcile request array if replica count miss match.
				if deploy.Spec.Replicas != c.Spec.Replicas {
					req = append(req, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Namespace: c.Namespace,
							Name:      c.Name,
						},
					})
				}

			}
		}
		return req
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&sayedppqqdevv1.CustomR{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, handlerForDeployment).
		Owns(&corev1.Service{}).
		Complete(r)
}
