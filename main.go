package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
	"github.com/orbit-ops/launchpad-core/config"
	"github.com/orbit-ops/launchpad-core/controller"
	"github.com/orbit-ops/launchpad-core/providers"
	"github.com/orbit-ops/launchpad-kubernetes/internal/provider"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var prov providers.Provider
	provCfg := &providers.ProviderConfig{
		ApiUrl: cfg.ApiUrl,
	}

	if cfg.Provider.Executable != "" {
		provCfg.ManagerExecutable = cfg.Provider.Executable
	} else {
		provCfg.ManagerExecutable = "docker.io/orbit-opts/launchpad-rocket-base-k8s:" + cfg.Version
	}

	var namespace string
	if cfg.Provider.Kubernetes != nil && cfg.Provider.Kubernetes.JobNamespace != "" {
		namespace = cfg.Provider.Kubernetes.JobNamespace
	} else if val, ok := os.LookupEnv("NAMESPACE"); ok {
		namespace = val
	} else {
		namespace = "launchpad"
	}

	prov, err = provider.NewKubernetesProvider(provCfg, namespace)
	if err != nil {
		log.Fatalf("Initializing %s provider: %v\n", cfg.Provider.Type, err)
	}

	_, err = controller.NewAccessController(prov)
	if err != nil {
		log.Fatal(err)
	}

}
