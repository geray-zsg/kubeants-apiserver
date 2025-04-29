package models

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type UserBindingSpec struct {
	User  string `json:"user,omitempty"`
	Scope Scope  `json:"scope,omitempty"`
	Role  string `json:"role,omitempty"`
}

type Scope struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserBindingSpec   `json:"spec,omitempty"`
	Status UserBindingStatus `json:"status,omitempty"`
}

type UserBindingStatus struct {
	// 是否成功下发RBAC资源（RoleBinding/ClusterRoleBinding）
	Synced bool `json:"synced,omitempty"`
	// 是否要求撤销权限（即计划回收）由上游（比如User被禁用）主动设置true
	Revoked bool `json:"revoked,omitempty"`
	// 权限是否已经被成功回收， cleanupRBAC成功后设置为true
	// RevokeComplete bool `json:"revokeComplete,omitempty"`
	// 最后一次成功同步RBAC的时间戳，用来做审计、前端展示
	LastSyncTime metav1.Time `json:"lastSyncTime,omitempty"`
	// 上一次失败的错误信息，如果applyRBAC/cleanupRBAC出错，记录下来方便排查
	// LastErrorReason string `json:"lastErrorReason,omitempty"`
	// 用于记录最后同步信息
	LastTransitionMsg string `json:"lastTransitionMsg,omitempty"` // 🌟 新增字段
	// LastAppliedGeneration is the generation of the last applied configuration
	LastAppliedGeneration int64 `json:"lastAppliedGeneration,omitempty"`
}
