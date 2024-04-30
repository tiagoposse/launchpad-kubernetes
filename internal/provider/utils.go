package provider

import (
	"fmt"

	"github.com/orbit-ops/launchpad-core/ent"
	"github.com/orbit-ops/launchpad-core/providers"
	batchv1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
)

func computeJobName(req *ent.Request) string {
	return fmt.Sprintf("launchpad-access-%s", req.ID)
}

func (kp *KubernetesProvider) getContainerEnv() []core.EnvVar {
	return []core.EnvVar{
		{
			Name:  "LAUNCHPAD_API_URL",
			Value: kp.conf.ApiUrl,
		},
		{
			Name: "MY_POD_NAME",
			ValueFrom: &core.EnvVarSource{
				FieldRef: &core.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "MY_POD_NAMESPACE",
			ValueFrom: &core.EnvVarSource{
				FieldRef: &core.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
	}
}

func (kp *KubernetesProvider) getManagerContainerSpec(cmd providers.ProviderCommand, name, code, conf string) batchv1.JobSpec {
	return batchv1.JobSpec{
		Template: core.PodTemplateSpec{
			Spec: core.PodSpec{
				Containers: []core.Container{
					{
						Name:  "manager",
						Image: kp.conf.ManagerExecutable,
						Env:   kp.getContainerEnv(),
						EnvFrom: []core.EnvFromSource{
							core.EnvFromSource{
								SecretRef: &core.SecretEnvSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: name,
									},
								},
							},
						},
					},
					{
						Name:  "main",
						Image: code,
						Env: []core.EnvVar{
							{
								Name:  "LAUNCHPAD_CONFIG",
								Value: conf,
							},
							{
								Name:  "LAUNCHPAD_COMMAND",
								Value: string(cmd),
							},
						},
					},
				},
				RestartPolicy: "Never",
			},
		},
	}
}
