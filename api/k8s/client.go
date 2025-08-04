package k8s

import (
	"appseclabs/database"
	"appseclabs/evaluation"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8s struct {
	Client              *kubernetes.Clientset
	Config              *rest.Config
	Database            *database.Database
	InitContainerConfig *InitContainerConfig
	Evaluator           *evaluation.Evaluator
}

type InitContainerConfig struct {
	GitUsername string
	GitPassword string
}

func NewK8sClient(db *database.Database) *K8s {
	config := getKubeConfig()
	clientset := getClientFromConfig(config)
	evaluator := evaluation.NewEvaluator()
	initContainerConfig := &InitContainerConfig{
		GitUsername: os.Getenv("GITHUB_INIT_CONTAINER_USERNAME"),
		GitPassword: os.Getenv("GITHUB_INIT_CONTAINER_PAT"),
	}
	return &K8s{Client: clientset, Config: config, Database: db, InitContainerConfig: initContainerConfig, Evaluator: evaluator}
}

func getKubeConfig() *rest.Config {
	encoded := os.Getenv("KUBECONFIG_BASE64")
	if encoded == "" {
		log.Fatal("KUBECONFIG_BASE64 environment variable not set")
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Failed to decode KUBECONFIG_BASE64: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(decoded)
	if err != nil {
		log.Fatalf("Failed to build config from kubeconfig content: %v", err)
	}

	return config

}

func getClientFromConfig(config *rest.Config) *kubernetes.Clientset {

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	fmt.Println("✅ Kubernetes client initialized")
	return clientset
}
