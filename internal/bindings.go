/*
Copyright © 2024 Furkan Pehlivan <furkanpehlivan34@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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

const (
	ClusterRoleBindingKind       = "ClusterRoleBinding"
	RoleBindingKind              = "RoleBinding"
	ClusterRoleBindingAPIVersion = "rbac.authorization.k8s.io/v1"
	RoleBindingAPIVersion        = "rbac.authorization.k8s.io/v1"
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
	var i int

	for _, crb := range bindings.ClusterRoleBindings.Items {
		crb.ManagedFields = nil

		data = append(data, Data{
			Name:     crb.Name,
			Id:       i,
			Kind:     ClusterRoleBindingKind,
			Subjects: crb.Subjects,
			RoleRef:  crb.RoleRef,
			Raw:      yamlParser(&crb, ClusterRoleBindingKind, ClusterRoleBindingAPIVersion),
		})
		i++
	}

	for _, rb := range bindings.RoleBindings.Items {
		rb.ManagedFields = nil
		data = append(data, Data{
			Name:     rb.Name,
			Id:       i,
			Kind:     RoleBindingKind,
			Subjects: rb.Subjects,
			RoleRef:  rb.RoleRef,
			Raw:      yamlParser(&rb, RoleBindingKind, RoleBindingAPIVersion),
		})
		i++
	}

	return data
}

func yamlParser(obj runtime.Object, kind string, apiVersion string) string {
	// Convert the object to YAML
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: true})
	o, err := runtime.Encode(s, obj)
	if err != nil {
		return ""
	}

	if len(o) == 0 {
		return "Could not parse the object"
	}

	// Prepend the kind and apiVersion to the YAML
	o = []byte(fmt.Sprintf("kind: %s\napiVersion: %s\n%s", kind, apiVersion, o))

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
