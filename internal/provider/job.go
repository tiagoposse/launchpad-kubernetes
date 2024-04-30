package provider

import (
	"context"
	"fmt"

	"github.com/orbit-ops/launchpad-core/ent"
	"github.com/orbit-ops/launchpad-core/providers"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kp *KubernetesProvider) createJob(ctx context.Context, cmd providers.ProviderCommand, name string, rocket *ent.Rocket) error {
	// rc, err := kp.EncodeRocketConfig(mission.)
	// if err != nil {
	// 	return err
	// }
	rc := ""

	job := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: kp.getManagerContainerSpec(cmd, name, rocket.Code, rc),
	}

	createdJob, err := kp.client.BatchV1().Jobs(kp.namespace).Create(context.Background(), job, v1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create Job: %v", err))
	}

	fmt.Printf("Job %s created\n", createdJob.Name)
	return nil
}

func (kp *KubernetesProvider) createCronjob(ctx context.Context, cmd providers.ProviderCommand, name string, rocket *ent.Rocket) error {
	// rc, err := kp.EncodeRocketConfig(req)
	// if err != nil {
	// 	return err
	// }
	rc := ""

	job := &batchv1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: "",
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Name: name,
				},
				Spec: kp.getManagerContainerSpec(cmd, name, rocket.Code, rc),
			},
		},
	}

	createdJob, err := kp.client.BatchV1().CronJobs(kp.namespace).Create(context.Background(), job, v1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create Job: %v", err))
	}

	fmt.Printf("Job %s created\n", createdJob.Name)
	return nil
}
