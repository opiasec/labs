package k8s

import (
	"appseclabs/types"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8s) CreateEvaluation(namespace string, lab *types.Lab) (types.LabFinishResult, error) {

	finishLabResult := types.LabFinishResult{}

	// Create each job for each evaluator
	for _, evaluator := range lab.Evaluators {
		err := k.createEvaluationJob(namespace, evaluator)
		if err != nil {
			log.Println("Error creating evaluation job:", err)
			finishLabResult.Status = "failed"
			finishLabResult.CriteriaResult = append(finishLabResult.CriteriaResult, types.LabFinishResultCriterion{
				Name:    evaluator.Slug,
				Status:  "failed",
				Message: fmt.Sprintf("Error creating job: %s", err.Error()),
			})
			continue
		}

		_, err = k.monitorEvaluationJob(namespace, evaluator)
		if err != nil {
			finishLabResult.Status = "failed"
			finishLabResult.CriteriaResult = append(finishLabResult.CriteriaResult, types.LabFinishResultCriterion{
				Name:    evaluator.Slug,
				Status:  "failed",
				Message: fmt.Sprintf("Job failed: %s", err.Error()),
			})
			continue
		}
		rawResult, err := k.GetResultFromJob(namespace, evaluator)
		if err != nil {
			finishLabResult.Status = "failed"
			finishLabResult.CriteriaResult = append(finishLabResult.CriteriaResult, types.LabFinishResultCriterion{
				Name:    evaluator.Slug,
				Status:  "failed",
				Message: fmt.Sprintf("Job failed: %s", err.Error()),
			})
			continue
		}

		parsedResult, err := k.Evaluator.Scorers[evaluator.Slug].Score(rawResult, evaluator)
		if err != nil {
			finishLabResult.Status = "failed"
			finishLabResult.CriteriaResult = append(finishLabResult.CriteriaResult, types.LabFinishResultCriterion{
				Name:    evaluator.Slug,
				Status:  "failed",
				Message: fmt.Sprintf("Job failed: %s", err.Error()),
			})
			continue
		}
		finishLabResult.CriteriaResult = append(finishLabResult.CriteriaResult, parsedResult)
	}
	totalScore := 0
	for _, criterion := range finishLabResult.CriteriaResult {
		totalScore += criterion.Score
	}
	finishLabResult.TotalScore = totalScore
	finishLabResult.Status = "completed"

	FilesDiff, err := k.GetCodeDiffJob(namespace, lab)
	if err != nil {
		finishLabResult.Status = "failed"
		finishLabResult.ErrorMessage = fmt.Sprintf("Error getting code diff: %s", err.Error())
	}
	finishLabResult.FilesDiff = FilesDiff

	return finishLabResult, nil
}

func (k *K8s) createEvaluationJob(namespace string, evaluator types.Evaluator) error {
	// Get the evaluation from the database
	dbEvaluation, err := k.Database.GetEvaluationBySlug(evaluator.Slug)
	if err != nil {
		return err
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      evaluator.Slug,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: func() []corev1.Container {
						containers := make([]corev1.Container, len(dbEvaluation.EvaluationSpec.Containers)+1)
						containers[0] = corev1.Container{
							Name:  "result-sidecar",
							Image: "alpine",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "result-volume",
									MountPath: "/results",
								},
							},
							Command: []string{"sh", "-c", `
							while [ ! -f /results/result.json ]; do sleep 1; done;
  							cat /results/result.json`},
						}
						for i, container := range dbEvaluation.EvaluationSpec.Containers {
							containers[i+1] = corev1.Container{
								Name:    container.Name,
								Image:   container.Image,
								Command: container.Commands,
								Args:    container.Args,
								Env: func() []corev1.EnvVar {
									envVars := make([]corev1.EnvVar, len(container.Env))
									for j, env := range container.Env {
										envVars[j] = corev1.EnvVar{
											Name:  env.Name,
											Value: env.Value,
										}
									}
									return envVars
								}(),
								VolumeMounts: func() []corev1.VolumeMount {
									volumeMounts := make([]corev1.VolumeMount, len(container.Volumes)+2)
									volumeMounts[0] = corev1.VolumeMount{
										Name:      "lab-volume",
										MountPath: "/workspace",
									}
									volumeMounts[1] = corev1.VolumeMount{
										Name:      "result-volume",
										MountPath: "/results",
									}
									for j, volume := range container.Volumes {
										volumeMounts[j+2] = corev1.VolumeMount{
											Name:      volume.Name,
											MountPath: volume.Path,
										}
									}
									return volumeMounts
								}(),
							}
						}
						return containers
					}(),
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: func() []corev1.Volume {
						volumes := make([]corev1.Volume, len(dbEvaluation.EvaluationSpec.Volumes)+2)
						volumes[0] = corev1.Volume{Name: "lab-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "code-server-pvc",
									ReadOnly:  false,
								}}}
						volumes[1] = corev1.Volume{
							Name: "result-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						}
						for i, volume := range dbEvaluation.EvaluationSpec.Volumes {
							volumes[i+2] = corev1.Volume{
								Name: volume.Name,
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							}
						}
						return volumes
					}(),
					InitContainers: func() []corev1.Container {
						initContainer := dbEvaluation.EvaluationSpec.InitContainer
						if initContainer.Name == "" {
							return nil
						}
						containers := make([]corev1.Container, 1)
						containers[0] = corev1.Container{
							Name:    initContainer.Name,
							Image:   initContainer.Image,
							Command: initContainer.Commands,
							Args:    initContainer.Args,
							Env: func() []corev1.EnvVar {
								envVars := make([]corev1.EnvVar, len(initContainer.Env))
								for j, env := range initContainer.Env {
									envVars[j] = corev1.EnvVar{
										Name:  env.Name,
										Value: env.Value,
									}
								}
								return envVars
							}(),
							VolumeMounts: func() []corev1.VolumeMount {
								volumeMounts := make([]corev1.VolumeMount, len(initContainer.Volumes)+2)
								volumeMounts[0] = corev1.VolumeMount{
									Name:      "lab-volume",
									MountPath: "/workspace",
								}
								volumeMounts[1] = corev1.VolumeMount{
									Name:      "result-volume",
									MountPath: "/results",
								}
								for j, volume := range initContainer.Volumes {
									volumeMounts[j+2] = corev1.VolumeMount{
										Name:      volume.Name,
										MountPath: volume.Path,
									}
								}
								return volumeMounts
							}(),
						}
						return containers
					}(),
				},
			},
		},
	}

	_, err = k.Client.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	return err

}

func (k *K8s) monitorEvaluationJob(namespace string, evaluator types.Evaluator) (*batchv1.Job, error) {
	for {
		job, err := k.Client.BatchV1().Jobs(namespace).Get(context.TODO(), evaluator.Slug, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		podsExpected := int32(1)
		if job.Spec.Completions != nil {
			podsExpected = *job.Spec.Completions
		}

		if job.Status.Succeeded == podsExpected {
			return job, nil
		} else if job.Status.Failed > 0 {
			failedString := job.Status.Conditions[0].Message
			jobLog, err := k.Client.CoreV1().Pods(namespace).GetLogs(job.Name, &corev1.PodLogOptions{}).DoRaw(context.TODO())
			if err != nil {
				return nil, err
			}
			return nil, errors.New("job failed: " + failedString + "\n" + string(jobLog))
		} else {
			time.Sleep(time.Second * 5)
			continue
		}
	}
}

func (k *K8s) GetResultFromJob(namespace string, evaluator types.Evaluator) (string, error) {
	job, err := k.Client.BatchV1().Jobs(namespace).Get(context.TODO(), evaluator.Slug, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if job.Status.Succeeded == 0 {
		return "", fmt.Errorf("job %s not yet succeeded", evaluator.Slug)
	}

	pods, err := k.Client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", evaluator.Slug),
	})
	if err != nil {
		return "", err
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pods found for job %s", evaluator.Slug)
	}

	podName := pods.Items[0].Name

	logs, err := k.Client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: "result-sidecar",
	}).DoRaw(context.TODO())
	if err != nil {
		return "", err
	}
	return string(logs), nil
}

func (k *K8s) GetCodeDiffJob(namespace string, lab *types.Lab) (string, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "save-code-job",
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "save-code",
							Image:   "ghcr.io/appsec-digital/eval-save-code:dev",
							Command: []string{"sh", "-c"},
							Args:    []string{"exec diff -ruN /lab/code /workspace | grep -v '^diff -ruN' "},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "code-volume",
									MountPath: "/lab/code",
								},
								{
									Name:      "lab-volume",
									MountPath: "/workspace",
								},
							},
						},
					},
					InitContainers: []corev1.Container{{
						Name:  "clone-lab",
						Image: "alpine/git",
						Env: []corev1.EnvVar{
							{Name: "GIT_URL", Value: lab.LabSpec.CodeConfig.GitURL},
							{Name: "GIT_BRANCH", Value: lab.LabSpec.CodeConfig.GitBranch},
							{Name: "GIT_PATH", Value: lab.LabSpec.CodeConfig.GitPath},
							{Name: "GIT_USERNAME", Value: k.InitContainerConfig.GitUsername},
							{Name: "GIT_PASSWORD", Value: k.InitContainerConfig.GitPassword},
						},
						Command:      []string{"sh", "-c", "set -e; echo 'Cloning lab...'; git clone --depth 1 --branch $GIT_BRANCH https://$GIT_USERNAME:$GIT_PASSWORD@$GIT_URL /tmp/lab; echo 'Copying lab...'; cp -r /tmp/lab/$GIT_PATH/* /lab/code; echo 'Chowning lab...'; chown -R 1000:1000 /lab/code"},
						VolumeMounts: []corev1.VolumeMount{{Name: "code-volume", MountPath: "/lab/code"}},
					}},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "code-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "lab-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "code-server-pvc",
									ReadOnly:  false,
								},
							},
						},
					},
				},
			},
		},
	}

	if _, err := k.Client.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{}); err != nil {
		return "", fmt.Errorf("failed to create job: %w", err)
	}

	for {
		job, err := k.Client.BatchV1().Jobs(namespace).Get(context.TODO(), job.Name, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to get job: %w", err)
		}
		if job.Status.Succeeded == 1 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	podList, err := k.Client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "job-name=" + job.Name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list pods for job: %w", err)
	}
	if len(podList.Items) == 0 {
		return "", fmt.Errorf("no pods found for job: %s", job.Name)
	}

	podName := podList.Items[0].Name

	logs, err := k.Client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: "save-code",
	}).DoRaw(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}

	return string(logs), nil

}

func (k *K8s) DeleteEvaluation(namespace string, evaluation types.LabEvaluation) error {
	return k.Client.BatchV1().Jobs(namespace).Delete(context.TODO(), evaluation.Name, metav1.DeleteOptions{})
}
