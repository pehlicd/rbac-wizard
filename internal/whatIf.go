package internal

import (
	"context"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (app App) ProcessClusterRoleBinding(crb *v1.ClusterRoleBinding) (data struct {
	Nodes []interface{} `json:"nodes"`
	Links []interface{} `json:"links"`
}) {
	for _, subject := range crb.Subjects {
		data.Nodes = append(data.Nodes, map[string]string{
			"id":       subject.Name,
			"kind":     subject.Kind,
			"apiGroup": subject.APIGroup,
		})
	}

	data.Nodes = append(data.Nodes, map[string]string{
		"id":       crb.RoleRef.Name,
		"kind":     crb.RoleRef.Kind,
		"apiGroup": crb.RoleRef.APIGroup,
	})

	roleRefInfo := fetchRoleRefDetails(app.KubeClient, crb.RoleRef)
	data.Nodes = append(data.Nodes, roleRefInfo)

	for _, subject := range crb.Subjects {
		data.Links = append(data.Links, map[string]string{
			"source": crb.Name,
			"target": subject.Name,
		})
	}

	data.Links = append(data.Links, map[string]string{
		"source": crb.Name,
		"target": crb.RoleRef.Name,
	})

	return data
}

func (app App) ProcessRoleBinding(rb *v1.RoleBinding) (data struct {
	Nodes []interface{} `json:"nodes"`
	Links []interface{} `json:"links"`
}) {
	for _, subject := range rb.Subjects {
		subjectInfo := fetchSubjectDetails(subject)
		data.Nodes = append(data.Nodes, subjectInfo)
	}

	data.Nodes = append(data.Nodes, map[string]interface{}{
		"id":       rb.RoleRef.Name,
		"kind":     rb.RoleRef.Kind,
		"apiGroup": rb.RoleRef.APIGroup,
	})

	roleRefInfo := fetchRoleRefDetails(app.KubeClient, rb.RoleRef)
	data.Nodes = append(data.Nodes, roleRefInfo)

	for _, subject := range rb.Subjects {
		data.Links = append(data.Links, map[string]interface{}{
			"source": subject.Name,
			"target": rb.RoleRef.Name,
		})
	}

	// also add roleBinding to the links
	data.Links = append(data.Links, map[string]interface{}{
		"source": rb.RoleRef.Name,
		"target": rb.Name,
	})

	return data
}

func fetchSubjectDetails(subject v1.Subject) map[string]interface{} {
	return map[string]interface{}{
		"id":       subject.Name,
		"kind":     subject.Kind,
		"apiGroup": subject.APIGroup,
	}
}

func fetchRoleRefDetails(client *kubernetes.Clientset, roleRef v1.RoleRef) map[string]string {
	switch roleRef.Kind {
	case "ClusterRole":
		// Fetch the ClusterRole details
		clusterRole, err := client.RbacV1().ClusterRoles().Get(context.TODO(), roleRef.Name, metav1.GetOptions{})
		if err != nil {
			return map[string]string{}
		}
		return map[string]string{
			"id":       clusterRole.Name,
			"kind":     "ClusterRole",
			"apiGroup": clusterRole.Kind,
		}
	case "Role":
		// Fetch the Role details
		role, err := client.RbacV1().Roles("").Get(context.TODO(), roleRef.Name, metav1.GetOptions{})
		if err != nil {
			return map[string]string{}
		}
		return map[string]string{
			"id":       role.Name,
			"kind":     "Role",
			"apiGroup": role.Kind,
		}
	}

	return map[string]string{}
}
