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
	yaml "gopkg.in/yaml.v2"
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
	deployname := instance.Name + "-deployment"
	configmapname := instance.Name + "-configmap"

	// configmap create
	found_config:=&corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configmapname,Namespace:instance.Namespace}, found_config)
	if err!= nil && errors.IsNotFound(err){
		conm:= r.configForCookiemonster(instance)
		reqLogger.Info("Creating a new configmap", "Configmap.Namespace", conm.Namespace, "Configmap.Name",conm.Name)
		err = r.client.Create(context.TODO(), conm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new configmap", "Configmap.Namespace", conm.Namespace, "Configmap.Name", conm.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to create configmap")
		return reconcile.Result{}, err
	}
	//

	// configmap update
	instancespec, err := yaml.Marshal(instance.Spec.Data)
	if err != nil {
		reqLogger.Error(err, "Failed to get exist configmap option", "Configmap.Namespace", instance.Namespace, "Configmap.Name", instance.Name)
		return reconcile.Result{}, err
	}
	stringspec:= string(instancespec)
	if !reflect.DeepEqual(found_config.Data["config.yaml"],stringspec){
		found_config.Data["config.yaml"] = stringspec
		reqLogger.Info("Change new option in cookiemonster deployment", "Deployment.Namespace", found_config.Namespace, "Deployment.Name", found_config.Name)
		err = r.client.Update(context.TODO(), found_config)
		if err != nil {
			reqLogger.Error(err, "Failed to change configmap", "confimap.Namespace", found_config.Namespace, "Deployment.Name", found_config.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	//


	// Deployment create
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployname, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentForCookiemonster(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}
	//
	//Deployment update
	if !reflect.DeepEqual(found.Spec.Replicas,r.deploymentForCookiemonster(instance).Spec.Replicas) {
		found = r.deploymentForCookiemonster(instance)
		reqLogger.Info("Change new spec in cookiemonster deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Failed to change Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	//

	// cr의 pod describe section
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForCookiemonster(instance.Name))
	listops := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), podList, listops)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "Cookiemonster.Namespace", instance.Namespace, "Cookiemonster.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		reqLogger.Info("CR cookimonster node update", "pod.Namespace", found.Namespace, "pod.Name", found.Name)
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Cookiemonster status(pod)")
			return reconcile.Result{}, err
		}
	}

	// cr의 configmap decsribe section
	configmaplist := &corev1.ConfigMapList{}
	err = r.client.List(context.TODO(), configmaplist, listops)
	if err != nil {
		reqLogger.Error(err, "Failed to list configmap", "Cookiemonster.Namespace", instance.Namespace, "Cookiemonster.Name", instance.Name)
		return reconcile.Result{}, err
	}
	configNames := getConfigmapNames(configmaplist.Items)
	if !reflect.DeepEqual(configNames, instance.Status.Maps) {
		instance.Status.Maps = configNames
		err := r.client.Status().Update(context.TODO(), instance)
		reqLogger.Info("CR cookimonster configmap update", "pod.Namespace", found.Namespace, "pod.Name", found.Name)
		if err != nil {
			reqLogger.Error(err, "Failed to update Cookiemonster status(configmap)")
			return reconcile.Result{}, err
		}
	}
	//

	return reconcile.Result{}, nil
	// STATUS MAP UPDATE

}

func getConfigmapNames(configmaps []corev1.ConfigMap) []string {
	var configNames []string
	for _, configmap := range configmaps {
		configNames = append(configNames, configmap.Name)
	}
	return configNames
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
	replicas := m.Spec.Data.Size
	deployname := m.Name + "-deployment"
	configmapname := m.Name + "-configmap"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployname,
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
										Name: configmapname,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

func (r *ReconcileCookiemonster) configForCookiemonster(m *rbxov1alpha1.Cookiemonster) *corev1.ConfigMap {
	ls := labelsForCookiemonster(m.Name)
	configmapname := m.Name + "-configmap"
	specfile, _ := yaml.Marshal(&m.Spec.Data)
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapname,
			Namespace: m.Namespace,
			Labels: ls,
		},
		Data: map[string]string{
				"config.yaml": string(specfile),
		},
	}
	// Set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(m, configmap, r.scheme)
	return configmap

}
