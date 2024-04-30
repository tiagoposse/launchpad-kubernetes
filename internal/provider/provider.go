package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/orbit-ops/launchpad-core/ent"
	"github.com/orbit-ops/launchpad-core/providers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesProvider struct {
	providers.BaseProvider

	client    *kubernetes.Clientset
	conf      *providers.ProviderConfig
	namespace string
}

func NewKubernetesProvider(c *providers.ProviderConfig, ns string) (*KubernetesProvider, error) {
	var config *rest.Config
	var err error

	if _, found := os.LookupEnv("KUBERNETES_SERVICE_HOST"); found {
		config, err = rest.InClusterConfig() // Use this if running inside Kubernetes cluster
	} else {
		var kubeConfigPath string
		if val, found := os.LookupEnv("KUBECONFIG"); found {
			kubeConfigPath = val
		} else {
			home, _ := os.UserHomeDir()
			kubeConfigPath = filepath.Join(home, ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create Kubernetes client: %v", err))
	}

	return &KubernetesProvider{
		client:    clientset,
		conf:      c,
		namespace: ns,
	}, nil
}

func (kp *KubernetesProvider) CreateAccess(ctx context.Context, token string, rocket *ent.Rocket, req *ent.Request) error {
	name := computeJobName(req)
	if err := kp.createSecret(ctx, name, token); err != nil {
		return err
	}

	if err := kp.createJob(ctx, providers.CreateAccess, name, rocket); err != nil {
		return err
	}

	return nil
}

func (kp *KubernetesProvider) PostCreateAccess(ctx context.Context, req *ent.Request) error {
	name := computeJobName(req)
	if err := kp.removeSecret(ctx, name); err != nil {
		return err
	}

	if err := kp.client.BatchV1().Jobs(kp.namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func (kp *KubernetesProvider) RemoveAccess(ctx context.Context, token string, rocket *ent.Rocket, req *ent.Request) error {
	name := computeJobName(req)

	if err := kp.createSecret(ctx, name, token); err != nil {
		return err
	}

	return kp.createCronjob(ctx, providers.RemoveAccess, name, rocket)
}

func (kp *KubernetesProvider) PostRemoveAccess(ctx context.Context, req *ent.Request) error {
	name := computeJobName(req)
	if err := kp.removeSecret(ctx, name); err != nil {
		return err
	}

	if err := kp.client.BatchV1().CronJobs(kp.namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}
