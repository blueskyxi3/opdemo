package resources

import (
	appv1beta1 "github.com/blueskyxi3/opdemo/v2/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func MutateService(app *appv1beta1.AppService, svc *corev1.Service) {
	svc.Spec = corev1.ServiceSpec{
		ClusterIP: svc.Spec.ClusterIP,
		Type:      corev1.ServiceTypeNodePort,
		Ports:     app.Spec.Ports,
		Selector: map[string]string{
			"app": app.Name,
		},
	}
}

func NewService(app *appv1beta1.AppService) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   appv1beta1.GroupVersion.Group,
					Version: appv1beta1.GroupVersion.Version,
					Kind:    appv1beta1.Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: app.Spec.Ports,
			Selector: map[string]string{
				"app": app.Name,
			},
		},
	}
}
