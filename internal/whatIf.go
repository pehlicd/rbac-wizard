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
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Node struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	ApiGroup string `json:"apiGroup"`
	Label    string `json:"label"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (app App) ProcessClusterRoleBinding(crb *v1.ClusterRoleBinding) (data struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}) {
	crbNodeID := crb.Kind + "-" + crb.Name
	data.Nodes = append(data.Nodes, Node{
		ID:       crbNodeID,
		Kind:     crb.Kind,
		ApiGroup: crb.APIVersion,
		Label:    crbNodeID,
	})

	for _, subject := range crb.Subjects {
		subjectInfo := fetchSubjectDetails(app.KubeClient, subject)
		if subjectInfo != nil {
			data.Nodes = append(data.Nodes, *subjectInfo)
			data.Links = append(data.Links, Link{
				Source: crbNodeID,
				Target: subjectInfo.ID,
			})
		}
	}

	roleRefInfo := fetchRoleRefDetails(app.KubeClient, crb.RoleRef)
	if roleRefInfo != nil {
		data.Nodes = append(data.Nodes, *roleRefInfo)
		data.Links = append(data.Links, Link{
			Source: crbNodeID,
			Target: roleRefInfo.ID,
		})
	}

	return data
}

func (app App) ProcessRoleBinding(rb *v1.RoleBinding) (data struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}) {
	rbNodeID := rb.Kind + "-" + rb.Name
	data.Nodes = append(data.Nodes, Node{
		ID:       rbNodeID,
		Kind:     rb.Kind,
		ApiGroup: rb.APIVersion,
		Label:    rbNodeID,
	})

	for _, subject := range rb.Subjects {
		subjectInfo := fetchSubjectDetails(app.KubeClient, subject)
		if subjectInfo != nil {
			data.Nodes = append(data.Nodes, *subjectInfo)
			data.Links = append(data.Links, Link{
				Source: rbNodeID,
				Target: subjectInfo.ID,
			})
		}
	}

	roleRefInfo := fetchRoleRefDetails(app.KubeClient, rb.RoleRef)
	if roleRefInfo != nil {
		data.Nodes = append(data.Nodes, *roleRefInfo)
		data.Links = append(data.Links, Link{
			Source: rbNodeID,
			Target: roleRefInfo.ID,
		})
	}

	return data
}

func fetchSubjectDetails(client *kubernetes.Clientset, subject v1.Subject) *Node {
	if subject.Kind == "ServiceAccount" {
		_, err := client.CoreV1().ServiceAccounts(subject.Namespace).Get(context.TODO(), subject.Name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
	}
	return &Node{
		ID:       subject.Name + "-" + subject.Kind,
		Kind:     subject.Kind,
		ApiGroup: subject.APIGroup,
		Label:    subject.Kind + ": " + subject.Name,
	}
}

func fetchRoleRefDetails(client *kubernetes.Clientset, roleRef v1.RoleRef) *Node {
	switch roleRef.Kind {
	case "ClusterRole":
		clusterRole, err := client.RbacV1().ClusterRoles().Get(context.TODO(), roleRef.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error fetching cluster role: %v\n", err)
			return nil
		}
		return &Node{
			ID:       clusterRole.Name + "-" + clusterRole.Kind,
			Kind:     clusterRole.Kind,
			ApiGroup: clusterRole.APIVersion,
			Label:    clusterRole.Kind + "-" + clusterRole.Name,
		}
	case "Role":
		return &Node{
			ID:       roleRef.Name + "-" + roleRef.Kind,
			Kind:     roleRef.Kind,
			ApiGroup: roleRef.APIGroup,
			Label:    roleRef.Kind + "-" + roleRef.Name,
		}
	}

	return nil
}
