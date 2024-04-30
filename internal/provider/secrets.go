package provider

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kp *KubernetesProvider) createSecret(ctx context.Context, name, token string) error {
	if _, err := kp.client.CoreV1().Secrets(kp.namespace).Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"token": token,
		},
	}, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func (kp *KubernetesProvider) removeSecret(ctx context.Context, name string) error {
	return kp.client.CoreV1().Secrets(kp.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
