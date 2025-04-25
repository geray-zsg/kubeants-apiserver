package models

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Metadata struct {
	Name string `json:"name"`
	// Namespace   string            `json:"namespace"` // 集群级别资源没有
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// RoleBinding defines a role binding for a user in a namespace or set of namespaces
type RoleBinding struct {
	// Name of the role binding
	Name string `json:"name,omitempty"`
	// Namespaces specifies the namespaces for the role binding
	Namespaces []string `json:"namespaces,omitempty"`
	// NamespaceSelector selects namespaces based on labels
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

// RoleBinding defines a role binding for a user in a namespace or set of namespaces
type Roles struct {
	// Name of the role binding
	Name string `json:"name,omitempty"`
	// Namespaces specifies the namespaces for the role binding
	Namespaces []string `json:"namespaces,omitempty"`
	// NamespaceSelector selects namespaces based on labels
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

type UserSpec struct {
	// State indicates the status of the user (e.g., Active, Inactive)
	State string `json:"state,omitempty"`
	// Name is the display name of the user
	Name string `json:"name,omitempty"`
	// 绑定的平台角色
	PlatformRoles string `json:"platformRoles,omitempty"`
	// Email is the email address of the user
	Email string `json:"email,omitempty"`
	// Phone is the phone number of the user
	Phone string `json:"phone,omitempty"`
	// Password is the hashed password of the user (not typically stored in a CRD, for security reasons)
	Password string `json:"password,omitempty"`
	// ClusterRoles specifies the cluster-wide role binding for the user
	ClusterRoles []string `json:"clusterroles,omitempty"`
	// Roles specifies the cluster-wide role binding for the user
	Roles []Roles `json:"roles,omitempty"`
	// ClusterRoleBinding specifies the cluster-wide role binding for the user
	ClusterRoleBinding string `json:"clusterrolebinding,omitempty"`
	// RoleBindings specifies the role bindings for the user in specific namespaces
	RoleBindings []RoleBinding `json:"rolebindings,omitempty"`
}

// UserStatus defines the observed state of User.
type UserStatus struct {
	// ServiceAccount is the name of the service account associated with the user
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// LastLoginTime is the time of the user's last login
	LastLoginTime metav1.Time `json:"lastLoginTime,omitempty"`
	// LastUpdatedTime is the time of the last update to the user's information
	LastUpdatedTime metav1.Time `json:"lastUpdatedTime,omitempty"`
	// LastAppliedGeneration is the generation of the last applied configuration
	LastAppliedGeneration int64 `json:"lastAppliedGeneration,omitempty"`
}

// User is the Schema for the users API.
type User struct {
	metav1.TypeMeta `json:",inline"`
	// metav1.ObjectMeta `json:"metadata,omitempty"`
	Metadata Metadata `json:"Metadata"`

	Spec   UserSpec   `json:"spec,omitempty"`
	Status UserStatus `json:"status,omitempty"`
}
