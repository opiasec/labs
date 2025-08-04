package k8s

import (
	"appseclabs/types"
	"appseclabs/utils"
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func (k *K8s) CreateLab(namespace string, lab *types.Lab) (string, error) {

	labels := map[string]string{"app": lab.Slug, "namespace": namespace}

	err := k.createNamespace(namespace)
	if err != nil {
		return "", fmt.Errorf("error creating namespace: %s", err)
	}

	err = k.createCodeServerPVC(namespace)
	if err != nil {
		return "", fmt.Errorf("error creating code server pvc: %s", err)
	}

	codeServerPassword, err := utils.GenerateRandomPassword(15)
	if err != nil {
		return "", fmt.Errorf("error generating random password: %s", err)
	}

	err = k.createDeployment(namespace, lab, labels, codeServerPassword)
	if err != nil {
		return "", fmt.Errorf("error creating deployment: %s", err)
	}

	err = k.createServices(namespace, lab, labels)
	if err != nil {
		return "", fmt.Errorf("error creating services: %s", err)
	}

	err = k.createIngress(namespace, lab, labels)
	if err != nil {
		return "", fmt.Errorf("error creating ingress: %s", err)
	}

	return codeServerPassword, nil
}

func (k *K8s) createCodeServerPVC(nsName string) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "code-server-pvc",
			Namespace: nsName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			//StorageClassName: strPtr("gp2"), -- AWS EKS storage class
		},
	}
	if _, err := k.Client.CoreV1().PersistentVolumeClaims(nsName).Create(context.TODO(), pvc, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *K8s) createIngress(nsName string, lab *types.Lab, labels map[string]string) error {
	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "code-server-ingress",
			Namespace: nsName,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/$2",
				"nginx.ingress.kubernetes.io/app-root":       fmt.Sprintf("/%s", nsName),
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: strPtr("nginx"),
			Rules: []networkingv1.IngressRule{{
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     fmt.Sprintf("/%s(/|$)(.*)", nsName),
								PathType: pathTypePtr("ImplementationSpecific"),
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: "code-server",
										Port: networkingv1.ServiceBackendPort{Number: 80},
									},
								},
							},
						},
					},
				},
			}},
		},
	}
	_, err := k.Client.NetworkingV1().Ingresses(nsName).Create(context.TODO(), ing, metav1.CreateOptions{})
	return err
}

func (k *K8s) createServices(nsName string, lab *types.Lab, labels map[string]string) error {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "code-server",
			Namespace: nsName,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}
	if _, err := k.Client.CoreV1().Services(nsName).Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil {
		return err
	}

	for _, service := range lab.LabSpec.Services {
		appSvc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      service.Name,
				Namespace: nsName,
			},
			Spec: corev1.ServiceSpec{
				Selector: labels,
				Ports: []corev1.ServicePort{
					{
						Port:       80,
						TargetPort: intstr.FromInt(int(service.Port)),
					},
				},
			},
		}
		if _, err := k.Client.CoreV1().Services(nsName).Create(context.TODO(), appSvc, metav1.CreateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (k *K8s) createDeployment(nsName string, lab *types.Lab, labels map[string]string, codeServerPassword string) error {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: lab.Slug, Namespace: nsName},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{Name: "code-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "code-server-pvc",
									ReadOnly:  false,
								}}},
						{Name: "docker-graph",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{}}},
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
					Containers: []corev1.Container{
						{
							Name:  "code-server",
							Image: "public.ecr.aws/s0d4j7b0/opiasec/code-server-docker:latest",
							Env: func() []corev1.EnvVar {
								envs := make([]corev1.EnvVar, 0, len(lab.LabSpec.Env)+3)
								envs = append(envs, corev1.EnvVar{Name: "LAB_BASE_URL", Value: os.Getenv("LAB_BASE_URL")})
								envs = append(envs, corev1.EnvVar{Name: "NAMESPACE", Value: nsName})
								envs = append(envs, corev1.EnvVar{Name: "PASSWORD", Value: codeServerPassword})
								for k, v := range lab.LabSpec.Env {
									envs = append(envs, corev1.EnvVar{Name: k, Value: v})
								}
								return envs
							}(),
							SecurityContext: &corev1.SecurityContext{
								RunAsUser:  func() *int64 { u := int64(0); return &u }(),
								Privileged: boolPtr(true),
							},
							Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "code-volume", MountPath: "/home/coder/project"},
								{Name: "docker-graph", MountPath: "/var/lib/docker"},
							},
						},
					},
				},
			},
		},
	}
	if _, err := k.Client.AppsV1().Deployments(nsName).Create(context.TODO(), deploy, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *K8s) createNamespace(namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	if _, err := k.Client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func (k *K8s) DeleteLab(namespace string) error {
	err := k.Client.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	return err
}

func (k *K8s) ExecInContainer(namespace, podName, containerName string, command []string) (stdout string, stderr string, err error) {
	execRequest := k.Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   command,
			Container: containerName,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Stdin:     false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k.Config, "POST", execRequest.URL())
	if err != nil {
		return "", "", err
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	err = executor.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
		Tty:    false,
	})

	return stdoutBuf.String(), stderrBuf.String(), err
}

func (k *K8s) ExecInPod(podName, namespace string, command []string) (stdout string, stderr string, err error) {
	execRequest := k.Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: command,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k.Config, "POST", execRequest.URL())
	if err != nil {
		return "", "", err
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	err = executor.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
		Tty:    false,
	})

	return stdoutBuf.String(), stderrBuf.String(), err
}

func (k *K8s) GetLabPod(namespace string) (corev1.Pod, error) {
	pods, err := k.Client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return corev1.Pod{}, err
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return pod, nil
		}
	}
	return corev1.Pod{}, fmt.Errorf("no running pod found")
}

func (k *K8s) ScaleCodeServer(namespace string, deploymentName string, replicas int32) error {
	deployment, err := k.Client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting deployment: %s", err)
	}
	deployment.Spec.Replicas = int32Ptr(replicas)
	_, err = k.Client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating deployment: %s", err)
	}
	return nil
}

func (k *K8s) GetLabPodContainerStatus(pod corev1.Pod) []ContainerStatusInfo {
	statuses := []ContainerStatusInfo{}

	for _, cs := range pod.Status.ContainerStatuses {
		state := ""
		reason := ""
		exitCode := int32(0)
		startedAt := time.Time{}

		if cs.State.Running != nil {
			state = "Running"
			startedAt = cs.State.Running.StartedAt.Time
		} else if cs.State.Waiting != nil {
			state = "Waiting"
			reason = cs.State.Waiting.Reason
		} else if cs.State.Terminated != nil {
			state = "Terminated"
			reason = cs.State.Terminated.Reason
			exitCode = cs.State.Terminated.ExitCode
			startedAt = cs.State.Terminated.StartedAt.Time
		}

		statuses = append(statuses, ContainerStatusInfo{
			Name:       cs.Name,
			Ready:      cs.Ready,
			State:      state,
			Reason:     reason,
			RestartCnt: cs.RestartCount,
			ExitCode:   exitCode,
			StartedAt:  startedAt,
		})
	}

	return statuses
}

func (k *K8s) GetLabStatusMessage(pod corev1.Pod) string {
	for _, c := range pod.Status.ContainerStatuses {
		if !c.Ready {
			if c.State.Waiting != nil {
				return fmt.Sprintf("waiting for container %s: %s", c.Name, c.State.Waiting.Reason)
			} else if c.State.Terminated != nil {
				return fmt.Sprintf("container %s failed: %s", c.Name, c.State.Terminated.Reason)
			}
		}
	}
	return "all ready"
}

func strPtr(s string) *string {
	return &s
}

func (k *K8s) GetAllNamespaces() error {

	namespaces, err := k.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}
	return nil
}
func (k *K8s) GetLab(namespace string) error {
	return nil
}
