package cookiemonster

import (
	"context"
	"k8s.io/apimachinery/pkg/labels"
	"reflect"
	appsv1 "k8s.io/api/apps/v1"
	rbxov1alpha1 "github.com/rbxorkt12/coockiemonster2/pkg/apis/rbxo/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_cookiemonster")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Cookiemonster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCookiemonster{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cookiemonster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Cookiemonster
	err = c.Watch(&source.Kind{Type: &rbxov1alpha1.Cookiemonster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Cookiemonster
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &rbxov1alpha1.Cookiemonster{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCookiemonster implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCookiemonster{}

// ReconcileCookiemonster reconciles a Cookiemonster object
type ReconcileCookiemonster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Cookiemonster object and makes changes based on the state read
// and what is in the Cookiemonster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCookiemonster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Cookiemonster")

	// Cookiemonster instance
	instance := &rbxov1alpha1.Cookiemonster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	//

	// configmap create
	
	//
	// Deployment
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) { // 만약에 Memcached를 위한 Deployment가 없다면
		// 새로운 Deployment를 생성합니다. deploymentForMemcached 함수는 Deployment를 위한 spec을 반환합니다.
		dep := r.deploymentForCookiemonster(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment가 성공적으로 생성되었다면, 이 이벤트를 다시 Requeue 합니다.
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Requeue된 이벤트는  (Deployment가 이미 생성되어 있기 때문에 위의 if문을 그냥 지나칠테고, 아래의 소스코드를 실행합니다.
	// 배포된 deployment의 spec이 요구되는 spec과 다른 경우 update를 실시한다.
	// 근데 애 딥이퀄이 안먹을것 같다.
	if !reflect.DeepEqual(found, r.deploymentForCookiemonster(instance)){
		found = r.deploymentForCookiemonster(instance)
		reqLogger.Info("Change new option in cookiemonster deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Failed to change Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	//

	// Memcached의 Status를 각 파드의 이름으로 업데이트합니다.
	// kubectl describe memcached를 해보면 Status 항목이 정의되어 있을 것입니다.
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForCookiemonster(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "Memcached.Namespace", instance.Namespace, "Memcached.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Memcached status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
	// STATUS MAP UPDATE

}
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
func labelsForCookiemonster(name string) map[string]string {
	return map[string]string{"app": "cookiemonster", "cookiemonster_cr": name}
}

func (r *ReconcileCookiemonster) deploymentForCookiemonster(m *rbxov1alpha1.Cookiemonster) *appsv1.Deployment {
	ls := labelsForCookiemonster(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "cookiemonster",
							Image:   "seungkyua/cookiemonster:latest",
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "config",
									MountPath:"/app/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name:"config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "cookiemonster-cm-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}
