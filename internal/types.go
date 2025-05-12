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
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rs/zerolog"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type App struct {
	KubeClient    *kubernetes.Clientset
	DynamicClient *dynamic.DynamicClient
	Logger        *zerolog.Logger
}

type Generator interface {
	GetBindings() (*Bindings, error)
}

type WhatIfGenerator interface {
	ProcessClusterRoleBinding(crb *v1.ClusterRoleBinding) struct {
		Nodes []Node `json:"nodes"`
		Links []Link `json:"links"`
	}
	ProcessRoleBinding(rb *v1.RoleBinding) struct {
		Nodes []Node `json:"nodes"`
		Links []Link `json:"links"`
	}
	ProcessClusterRoleTemplateBinding(crtb *v3.ClusterRoleTemplateBinding) struct {
		Nodes []Node `json:"nodes"`
		Links []Link `json:"links"`
	}
	ProcessProjectRoleTemplateBinding(crtb *v3.ProjectRoleTemplateBinding) struct {
		Nodes []Node `json:"nodes"`
		Links []Link `json:"links"`
	}
	ProcessGlobalRoleBinding(crtb *v3.GlobalRoleBinding) struct {
		Nodes []Node `json:"nodes"`
		Links []Link `json:"links"`
	}
}

type Bindings struct {
	ClusterRoleBindings         *v1.ClusterRoleBindingList         `json:"clusterRoleBindings"`
	RoleBindings                *v1.RoleBindingList                `json:"roleBindings"`
	ClusterRoleTemplateBindings *v3.ClusterRoleTemplateBindingList `json:"clusterRoleTemplateBindings"`
	ProjectRoleTemplateBindings *v3.ProjectRoleTemplateBindingList `json:"projectRoleTemplateBindings"`
	GlobalRoleBindings          *v3.GlobalRoleBindingList          `json:"globalRoleBindings"`
}

type Data struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Kind     string       `json:"kind"`
	Subjects []v1.Subject `json:"subjects"`
	RoleRef  v1.RoleRef   `json:"roleRef"`
	Raw      string       `json:"raw"`
}
