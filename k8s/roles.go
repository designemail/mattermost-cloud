// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package k8s

import (
	rbacv1 "k8s.io/api/rbac/v1"
	rbacbetav1 "k8s.io/api/rbac/v1beta1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kc *KubeClient) createOrUpdateClusterRoleV1(account *rbacv1.ClusterRole) (metav1.Object, error) {
	_, err := kc.Clientset.RbacV1().ClusterRoles().Get(account.GetName(), metav1.GetOptions{})
	if err != nil && !k8sErrors.IsNotFound(err) {
		return nil, err
	}

	if err != nil && k8sErrors.IsNotFound(err) {
		return kc.Clientset.RbacV1().ClusterRoles().Create(account)
	}

	return kc.Clientset.RbacV1().ClusterRoles().Update(account)
}

func (kc *KubeClient) createOrUpdateClusterRoleBetaV1(account *rbacbetav1.ClusterRole) (metav1.Object, error) {
	_, err := kc.Clientset.RbacV1beta1().ClusterRoles().Get(account.GetName(), metav1.GetOptions{})
	if err != nil && !k8sErrors.IsNotFound(err) {
		return nil, err
	}

	if err != nil && k8sErrors.IsNotFound(err) {
		return kc.Clientset.RbacV1beta1().ClusterRoles().Create(account)
	}

	return kc.Clientset.RbacV1beta1().ClusterRoles().Update(account)
}

func (kc *KubeClient) createOrUpdateRoleV1(account *rbacv1.Role) (metav1.Object, error) {
	_, err := kc.Clientset.RbacV1().Roles(account.GetNamespace()).Get(account.GetName(), metav1.GetOptions{})
	if err != nil && !k8sErrors.IsNotFound(err) {
		return nil, err
	}

	if err != nil && k8sErrors.IsNotFound(err) {
		return kc.Clientset.RbacV1().Roles(account.GetNamespace()).Create(account)
	}

	return kc.Clientset.RbacV1().Roles(account.GetNamespace()).Update(account)
}

func (kc *KubeClient) createOrUpdateRoleBetaV1(account *rbacbetav1.Role) (metav1.Object, error) {
	_, err := kc.Clientset.RbacV1beta1().Roles(account.GetNamespace()).Get(account.GetName(), metav1.GetOptions{})
	if err != nil && !k8sErrors.IsNotFound(err) {
		return nil, err
	}

	if err != nil && k8sErrors.IsNotFound(err) {
		return kc.Clientset.RbacV1beta1().Roles(account.GetNamespace()).Create(account)
	}

	return kc.Clientset.RbacV1beta1().Roles(account.GetNamespace()).Update(account)
}
