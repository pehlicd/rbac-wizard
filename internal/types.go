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
	"github.com/rs/zerolog"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

type App struct {
	KubeClient *kubernetes.Clientset
	Logger     *zerolog.Logger
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
