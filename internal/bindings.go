package internal

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
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
			Raw:      yamlParser(&crb),
		})
	}

	for i, rb := range bindings.RoleBindings.Items {
		data = append(data, Data{
			Name:     rb.Name,
			Id:       i,
			Kind:     "RoleBinding",
			Subjects: rb.Subjects,
			RoleRef:  rb.RoleRef,
			Raw:      yamlParser(&rb),
		})
	}

	return data
}

func yamlParser(obj runtime.Object) string {
	// Convert the object to YAML
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: true})
	o, err := runtime.Encode(s, obj)
	if err != nil {
		return ""
	}

	// Unmarshal the JSON into a generic map
	var yamlObj map[string]interface{}
	err = yaml.Unmarshal(o, &yamlObj)
	if err != nil {
		return err.Error()
	}

	// Marshal the map back into YAML
	yamlData, err := yaml.Marshal(yamlObj)
	if err != nil {
		return err.Error()
	}

	return string(yamlData)
}
