package provider

import (
	"fmt"

	"k8s.io/client-go/tools/auth"
)

func (kp *KubernetesProvider) ValidateAuthToken(token string) {

	podName, serviceAccountName, namespace, err := auth.GetInfoFromToken(token, kp.client.Discovery())
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get information from token: %v", err)
	}

}
