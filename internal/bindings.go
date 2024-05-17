package internal

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (app App) GetBindings() (*Bindings, error) {
	clientset := app.KubeClient

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
