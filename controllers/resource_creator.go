package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	sayedppqqdevv1 "sayedppqq.dev/mycrd/api/v1"
)

func newDeployment(customR *sayedppqqdevv1.CustomR, deploymentName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: customR.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(customR, sayedppqqdevv1.GroupVersion.WithKind("CustomR")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: customR.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  customR.Name,
					"Kind": "CustomR",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: customR.Namespace,
					Labels: map[string]string{
						"app":  customR.Name,
						"Kind": "CustomR",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  customR.Name,
							Image: customR.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: customR.Spec.Container.Port,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			},
		},
	}
}
func newService(customR *sayedppqqdevv1.CustomR, serviceName string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: customR.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(customR, sayedppqqdevv1.GroupVersion.WithKind("CustomR")),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       customR.Spec.Service.ServicePort,
					TargetPort: intstr.FromInt(int(customR.Spec.Container.Port)),
					NodePort:   customR.Spec.Service.ServiceNodePort,
				},
			},
			Selector: map[string]string{
				"app":  customR.Name,
				"Kind": "CustomR",
			},
			Type: func() corev1.ServiceType {
				if customR.Spec.Service.ServiceType == "NodePort" {
					return corev1.ServiceTypeNodePort
				} else {
					return corev1.ServiceTypeClusterIP
				}
			}(),
		},
	}
}
