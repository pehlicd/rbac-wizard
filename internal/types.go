package internal

import (
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

type App struct {
	KubeClient *kubernetes.Clientset
}

type Generator interface {
	GetBindings() (*Bindings, error)
	ProcessClusterRoleBinding(crb *v1.ClusterRoleBinding) struct {
		Nodes []interface{} `json:"nodes"`
		Links []interface{} `json:"links"`
	}
	ProcessRoleBinding(rb *v1.RoleBinding) struct {
		Nodes []interface{} `json:"nodes"`
		Links []interface{} `json:"links"`
	}
}

type Bindings struct {
	ClusterRoleBindings *v1.ClusterRoleBindingList `json:"clusterRoleBindings"`
	RoleBindings        *v1.RoleBindingList        `json:"roleBindings"`
}

type Data struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Kind     string       `json:"kind"`
	Subjects []v1.Subject `json:"subjects"`
	RoleRef  v1.RoleRef   `json:"roleRef"`
	Raw      string       `json:"raw"`
}
