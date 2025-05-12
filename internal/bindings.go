/*
Modified by Alessio Greggi © 2025. Based on work by Furkan Pehlivan <furkanpehlivan34@gmail.com>.

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

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	ClusterRoleBindingKind               = "ClusterRoleBinding"
	RoleBindingKind                      = "RoleBinding"
	ClusterRoleTemplateBindingKind       = "clusterroletemplatebindings"
	ProjectRoleTemplateBindingKind       = "projectroletemplatebindings"
	GlobalRoleBindingKind                = "globalrolebindings"
	ManagementAPI                        = "management.cattle.io"
	ManagementVersion                    = "v3"
	ClusterRoleBindingAPIVersion         = "rbac.authorization.k8s.io/v1"
	RoleBindingAPIVersion                = "rbac.authorization.k8s.io/v1"
	ProjectRoleTemplateBindingAPIVersion = "management.cattle.io/v3"
	ClusterRoleTemplateBindingAPIVersion = "management.cattle.io/v3"
	GlobalRoleBindingAPIVersion          = "management.cattle.io/v3"
)

func (app App) GetBindings() (*Bindings, error) {
	crbs, err := app.KubeClient.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing ClusterRoleBindings:", err)
		return nil, err
	}

	rbs, err := app.KubeClient.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing RoleBindings:", err)
		return nil, err
	}

	globalRoleBindingGVR := schema.GroupVersionResource{
		Group:    ManagementAPI,
		Version:  ManagementVersion,
		Resource: GlobalRoleBindingKind,
	}

	projectRoleTemplateBindingGVR := schema.GroupVersionResource{
		Group:    ManagementAPI,
		Version:  ManagementVersion,
		Resource: ProjectRoleTemplateBindingKind,
	}

	clusterRoleTemplateBindingGVR := schema.GroupVersionResource{
		Group:    ManagementAPI,
		Version:  ManagementVersion,
		Resource: ClusterRoleTemplateBindingKind,
	}

	grbs := &v3.GlobalRoleBindingList{}
	prtbs := &v3.ProjectRoleTemplateBindingList{}
	crtbs := &v3.ClusterRoleTemplateBindingList{}

	// missing rancher crtbs, prtbs, grbs
	grbUnstructured, err := app.DynamicClient.Resource(globalRoleBindingGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing GlobalRoleBindings:", err)
		return nil, err
	}

	prtbUnstructured, err := app.DynamicClient.Resource(projectRoleTemplateBindingGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing ProjectRoleTemplateBindings:", err)
		return nil, err
	}

	crtbUnstructured, err := app.DynamicClient.Resource(clusterRoleTemplateBindingGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing ClusterRoleTemplateBindings:", err)
		return nil, err
	}

	if grbUnstructured != nil {
		grbs.Items = make([]v3.GlobalRoleBinding, len(grbUnstructured.Items))
		for i, item := range grbUnstructured.Items {
			bindingObj := &v3.GlobalRoleBinding{}
			if err := convertUnstructuredToTyped(&item, bindingObj); err != nil {
				fmt.Printf("Warning: Error converting GlobalRoleBinding: %v\n", err)
				continue
			}
			grbs.Items[i] = *bindingObj
		}
	}

	if prtbUnstructured != nil {
		prtbs.Items = make([]v3.ProjectRoleTemplateBinding, len(prtbUnstructured.Items))
		for i, item := range prtbUnstructured.Items {
			bindingObj := &v3.ProjectRoleTemplateBinding{}
			if err := convertUnstructuredToTyped(&item, bindingObj); err != nil {
				fmt.Printf("Warning: Error converting ProjectRoleTemplateBinding: %v\n", err)
				continue
			}
			prtbs.Items[i] = *bindingObj
		}
	}

	if crtbUnstructured != nil {
		crtbs.Items = make([]v3.ClusterRoleTemplateBinding, len(crtbUnstructured.Items))
		for i, item := range crtbUnstructured.Items {
			bindingObj := &v3.ClusterRoleTemplateBinding{}
			if err := convertUnstructuredToTyped(&item, bindingObj); err != nil {
				fmt.Printf("Warning: Error converting ClusterRoleTemplateBinding: %v\n", err)
				continue
			}
			crtbs.Items[i] = *bindingObj
		}
	}

	return &Bindings{
		ClusterRoleBindings:         crbs,
		RoleBindings:                rbs,
		ClusterRoleTemplateBindings: crtbs,
		ProjectRoleTemplateBindings: prtbs,
		GlobalRoleBindings:          grbs,
	}, nil
}

// Helper function to convert unstructured objects to typed objects
func convertUnstructuredToTyped(unstructuredObj *unstructured.Unstructured, typedObj interface{}) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, typedObj)
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

	for _, grb := range bindings.GlobalRoleBindings.Items {
		grb.ManagedFields = nil
		user := v1.Subject{
			APIGroup: ManagementAPI,
			Kind:     "User",
			Name:     grb.UserName,
		}
		role := v1.RoleRef{
			APIGroup: ManagementAPI,
			Kind:     "GlobalRole",
			Name:     grb.GlobalRoleName,
		}
		data = append(data, Data{
			Name:     grb.Name,
			Id:       i,
			Kind:     GlobalRoleBindingKind,
			Subjects: []v1.Subject{user},
			RoleRef:  role,
			Raw:      yamlParser(&grb, GlobalRoleBindingKind, GlobalRoleBindingAPIVersion),
		})
		i++
	}

	for _, prtb := range bindings.ProjectRoleTemplateBindings.Items {
		prtb.ManagedFields = nil
		user := v1.Subject{
			APIGroup: ManagementAPI,
			Kind:     "User",
			Name:     prtb.UserName,
		}
		roleProject := v1.RoleRef{
			APIGroup: ManagementAPI,
			Kind:     "Project",
			Name:     prtb.ProjectName,
		}
		roleTemplate := v1.RoleRef{
			APIGroup: ManagementAPI,
			Kind:     "RoleTemplate",
			Name:     prtb.RoleTemplateName,
		}
		data = append(data, Data{
			Name:     prtb.Name,
			Id:       i,
			Kind:     ProjectRoleTemplateBindingKind,
			Subjects: []v1.Subject{user},
			RoleRef:  roleProject,
			Raw:      yamlParser(&prtb, ProjectRoleTemplateBindingKind, ProjectRoleTemplateBindingAPIVersion),
		})
		data = append(data, Data{
			Name:     prtb.Name,
			Id:       i,
			Kind:     ProjectRoleTemplateBindingKind,
			Subjects: []v1.Subject{user},
			RoleRef:  roleTemplate,
			Raw:      yamlParser(&prtb, ProjectRoleTemplateBindingKind, ProjectRoleTemplateBindingAPIVersion),
		})
		i++
	}

	for _, crtb := range bindings.ClusterRoleTemplateBindings.Items {
		crtb.ManagedFields = nil
		user := v1.Subject{
			APIGroup: ManagementAPI,
			Kind:     "User",
			Name:     crtb.UserName,
		}
		roleCluster := v1.RoleRef{
			APIGroup: ManagementAPI,
			Kind:     "Cluster",
			Name:     crtb.ClusterName,
		}
		roleTemplate := v1.RoleRef{
			APIGroup: ManagementAPI,
			Kind:     "RoleTemplate",
			Name:     crtb.RoleTemplateName,
		}
		data = append(data, Data{
			Name:     crtb.Name,
			Id:       i,
			Kind:     ClusterRoleTemplateBindingKind,
			Subjects: []v1.Subject{user},
			RoleRef:  roleCluster,
			Raw:      yamlParser(&crtb, ClusterRoleTemplateBindingKind, ClusterRoleTemplateBindingAPIVersion),
		})
		data = append(data, Data{
			Name:     crtb.Name,
			Id:       i,
			Kind:     ClusterRoleTemplateBindingKind,
			Subjects: []v1.Subject{user},
			RoleRef:  roleTemplate,
			Raw:      yamlParser(&crtb, ClusterRoleTemplateBindingKind, ClusterRoleTemplateBindingAPIVersion),
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
