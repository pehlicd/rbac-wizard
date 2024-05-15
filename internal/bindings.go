package internal

import (
	"context"
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Bindings struct {
	ClusterRoleBindings *v1.ClusterRoleBindingList `json:"clusterRoleBindings"`
	RoleBindings        *v1.RoleBindingList        `json:"roleBindings"`
}

type Data struct {
	Name     string       `json:"name"`
	Id       int          `json:"id"`
	Kind     string       `json:"kind"`
	Subjects []v1.Subject `json:"subjects"`
	RoleRef  v1.RoleRef   `json:"roleRef"`
	Raw      string       `json:"raw"`
}

func GetBindings() (*Bindings, error) {
	clientset, err := getClientset()
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return nil, err
	}

	crbs, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing cluster role bindings:", err)
		return nil, err
	}

	rbs, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing role bindings:", err)
		return nil, err
	}

	return &Bindings{
		ClusterRoleBindings: crbs,
		RoleBindings:        rbs,
	}, nil
}

func GenerateData(bindings *Bindings) []Data {
	var data []Data

	for i, crb := range bindings.ClusterRoleBindings.Items {
		data = append(data, Data{
			Name:     crb.Name,
			Id:       i,
			Kind:     "ClusterRoleBinding",
			Subjects: crb.Subjects,
			RoleRef:  crb.RoleRef,
			Raw:      crb.String(),
		})
	}

	for i, rb := range bindings.RoleBindings.Items {
		data = append(data, Data{
			Name:     rb.Name,
			Id:       i,
			Kind:     "RoleBinding",
			Subjects: rb.Subjects,
			RoleRef:  rb.RoleRef,
			Raw:      rb.String(),
		})
	}

	return data
}
