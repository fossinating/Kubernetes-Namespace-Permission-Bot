/*
Copyright 2020.

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
	//	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	//	"k8s.io/apimachinery/pkg/types"

	"io/ioutil"

	"encoding/json"

	"context"
	//"fmt"
	"strings"

	"github.com/go-logr/logr"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	//rbac "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

// Reconciles Namespaces
type NamespaceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

/* Again, these don't work yet so we're disabling them.
// Reconciles ConfigMaps
type ConfigMapReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

// Reconciles Namespaces
type RolebindingReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}*/

type Subject struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Group struct {
	Prefix   string               `json:"prefix,omitempty"`
	Label    string               `json:"label,omitempty"`
	LabelKey string               `json:"label_key,omitempty"`
	Roles    map[string][]Subject `json:"roles"`
}

type Conf struct {
	Mode      string  `json:"mode"`
	LabelKey  string  `json:"label_key,omitempty"`
	Delimiter string  `json:"delimiter"`
	Groups    []Group `json:"groups"`
}

// +kubebuilder:rbac:groups=permission-bot.mrriptide.github.io,resources=permissionbots,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=permission-bot.mrriptide.github.io,resources=permissionbots/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=permission-bot.mrriptide.github.io,resources=permissionbots/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PermissionBot object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("permissionbot", req.NamespacedName)

	// Get namespace

	namespace, err := GetNamespace(log, ctx, r, req)
	if namespace == nil {
		return ctrl.Result{}, err
	}

	err = ReconcileNamespace(log, ctx, r.Config, namespace)
	return ctrl.Result{}, err
}

func NamespacePredicate(r *NamespaceReconciler) predicate.Predicate {
	log := r.Log.WithValues("permissionbot", "predicate")
	config, _ := GetConfig(log)
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return config.Mode == "authoritative"
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return config.Mode == "authoritative"
		},
	}
}

// Sets up the namespace controller with the manager
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		WithEventFilter(NamespacePredicate(r)).
		Complete(r)
}

/*

Commenting out the code so nothing gets accidentally called

// Reconciler for configmap events
func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("permissionbot", req.NamespacedName)

	configName := types.NamespacedName{Namespace: "namespace-permission-bot-system", Name: "permission-bot-config"}
	if req.NamespacedName == configName {
		log.Info("My config map has been updated")

		// Update config.json

		configMap := &corev1.ConfigMap{}
		err := r.Get(ctx, req.NamespacedName, configMap)
		if err != nil {
			log.Error(err, "Failed to retrieve config map")
			return ctrl.Result{}, nil
		}

		bytes := []byte(configMap.Data["config.json"])
		err = ioutil.WriteFile("config.json", bytes, 0644)
		if err != nil {
			log.Error(err, "Failed to write config map")
			return ctrl.Result{}, nil
		}

		// Go through all namespaces and reconcile them


	} else {
		// Configmap doesnt match mine
		//log.Info(req.NamespacedName.String() + " != " + configName.String())
	}

	return ctrl.Result{}, nil
}

func ConfigMapPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	}
}

// Sets up the configmap controller with the manager
func (r *ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithEventFilter(ConfigMapPredicate()).
		Complete(r)
}

// Reconciler for rolebinding events
func (r *RolebindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("permissionbot", req.NamespacedName)

	log.Info("Test")
	// Get the rolebinding from the event
	rolebinding := &rbacv1.RoleBinding{}
	err := r.Get(ctx, req.NamespacedName, rolebinding)
	if err != nil {
		log.Error(err, "Could not retrieve the rolebinding")
		return ctrl.Result{}, err
	}

	if strings.HasSuffix(rolebinding.ObjectMeta.Name, "-auto-permission") {
		log.Info("This triggers the reconcile")
		// So I need to get a proper namespace object, but I don't have it yet??
		// More my problem is I need to prevent an infinite loop
		//ReconcileNamespace(log, ctx, r.Config, rolebinding.ObjectMeta.Namespace)
	}

	return ctrl.Result{}, nil
}

func RolebindingPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	}
}

// Sets up the configmap controller with the manager
func (r *RolebindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rbacv1.RoleBinding{}).
		WithEventFilter(RolebindingPredicate()).
		Complete(r)
}*/

func printJSON(log logr.Logger, obj interface{}) {
	val, err := json.Marshal(obj)

	if err != nil {
		log.Error(err, "Failed to marshal")
	} else {
		log.Info(string(val))
	}
}

func NamespaceContainsRolebinding(existing []rbacv1.RoleBinding, role_binding_name string) bool {
	for _, a := range existing {
		if a.ObjectMeta.Name == role_binding_name {
			return true
		}
	}
	return false
}

func CreateRoleBinding(ctx context.Context, log logr.Logger, mode string, existing *rbacv1.RoleBindingList, namespace_name string, clientset kubernetes.Interface, group_name string, role string, targets []Subject) error {
	// Set the name of the rolebinding

	meta := metav1.ObjectMeta{
		Name:      role + "-" + group_name + "-auto-permission",
		Namespace: namespace_name,
	}

	// Get the role reference

	role_ref := rbacv1.RoleRef{
		Kind:     "ClusterRole",
		Name:     role,
		APIGroup: "rbac.authorization.k8s.io",
	}

	// Select the subjects
	var subjects []rbacv1.Subject = []rbacv1.Subject{}
	for _, target := range targets {
		subjects = append(subjects, rbacv1.Subject{
			Kind:     target.Type,
			APIGroup: "rbac.authorization.k8s.io",
			Name:     target.Name,
		})
	}

	// Create Rolebinding

	binding := rbacv1.RoleBinding{
		Subjects:   subjects,
		RoleRef:    role_ref,
		ObjectMeta: meta,
	}

	// Check if rolebinding exists in namespace already and if it does, update instead of create

	if NamespaceContainsRolebinding(existing.Items, role+"-"+group_name+"-auto-permission") {
		if mode == "authoritative" {
			log.Info("Rolebinding \"" + role + "-" + group_name + "-auto-permission" + "\" already exists in namespace, updating it now")

			_, err := clientset.RbacV1().RoleBindings(namespace_name).Update(ctx, &binding, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		} else if mode == "passive" {
			log.Info("Rolebinding \"" + role + "-" + group_name + "-auto-permission" + "\" already exists in namespace, ignoring it since on passive mode")
		} else {
			log.Info("Rolebinding \"" + role + "-" + group_name + "-auto-permission" + "\" already exists in namespace, ignoring it since on unknown \"" + mode + "\" mode")
		}

	} else {
		_, err := clientset.RbacV1().RoleBindings(namespace_name).Create(ctx, &binding, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetConfig(log logr.Logger) (Conf, error) {
	// Read permissions config file, let's hope it's here

	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Error(err, "Could not find config file")
		return Conf{}, err
	}

	var config Conf
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Error(err, "Could not process config file")
		return Conf{}, err
	}

	return config, nil
}

func GetClientSet(log logr.Logger, kube_config *rest.Config) (kubernetes.Interface, error) {
	// Since this is running in the cluster, no config file path is needed
	/*kube_config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		log.Error(err, "Failed to get the kubernetes config")
		return nil, err
	}*/

	// Get ClientSet from config

	clientset, err := kubernetes.NewForConfig(kube_config)
	if err != nil {
		log.Error(err, "Failed to get the clientset")
		return nil, err
	}

	return clientset, nil
}

func GetNamespace(log logr.Logger, ctx context.Context, r *NamespaceReconciler, req ctrl.Request) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{}
	err := r.Get(ctx, req.NamespacedName, namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Namespace deleted")
			return nil, nil
		}
		log.Error(err, "Failed to get Namespace")
		return nil, err
	}
	return namespace, nil
}

func ReconcileNamespace(log logr.Logger, ctx context.Context, kube_config *rest.Config, namespace *corev1.Namespace) error {
	// Get kubernetes clientset

	clientset, err := GetClientSet(log, kube_config)
	if err != nil {
		return err
	}

	// Get configuration for the operator

	config, err := GetConfig(log)
	if err != nil {
		return err
	}

	/*var prefix string = ""
	if len(strings.Split(namespace.ObjectMeta.Name, "-")) > 1 {
		prefix = strings.Split(namespace.ObjectMeta.Name, "-")[0]
	}*/
	var label string = ""
	if namespace.ObjectMeta.Labels != nil {
		if namespace.ObjectMeta.Labels[config.LabelKey] != "" {
			label = namespace.ObjectMeta.Labels[config.LabelKey]
		}
	}

	log.Info("Namespace \"" + namespace.ObjectMeta.Name + "\" has label \"" + string(label) + "\"")

	// Get list of existing rolebindings in namespace

	options := metav1.ListOptions{} // Create listoptions

	existing, err := clientset.RbacV1().RoleBindings(namespace.ObjectMeta.Name).List(ctx, options)
	if err != nil {
		log.Error(err, "Failed to get list of existing rolebindings")
		return err
	}

	// Go through all groups to check for matches
	//
	// Currently, multiple groups can match one namespace
	for _, group := range config.Groups {
		// Use unique LabelKey if it exists
		if group.LabelKey != "" {
			label = namespace.ObjectMeta.Labels[group.LabelKey]
		}
		// Check if prefix or label value matchesg
		if (strings.HasPrefix(namespace.ObjectMeta.Name, group.Prefix+config.Delimiter) && group.Prefix != "") || (group.Label == label && group.Label != "") {
			log.Info("Namespace has group rules defined")

			for role, target := range group.Roles {
				err = CreateRoleBinding(ctx, log, config.Mode, existing, namespace.ObjectMeta.Name, clientset, group.Prefix, role, target)
				if err != nil {
					log.Error(err, "Failed to create rolebinding \""+role+"\"")
					return err
				}
			}
		}
	}

	return nil
}
