package config

import (
	"fmt"
	"log"
	"os"
)

var envVars = []string{
	"KUBECONFIG_BASE64",
	"MONGODB_URI",
	"MONGODB_DATABASE",
	"LAB_BASE_URL",
	"REDIS_ADDR",
	"JWKS_URL",
	"GITHUB_INIT_CONTAINER_USERNAME",
	"GITHUB_INIT_CONTAINER_PAT",
}

func CheckEnvVariables() {
	allEnvIsSet := true
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("%s is not set\n", envVar)
			allEnvIsSet = false
		}
	}
	if allEnvIsSet {
		fmt.Println("✅ All env variables are set")
	} else {
		log.Fatal("Some env variables are not set")
	}
}
