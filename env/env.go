package env

import "os"

const (
	ClusterIDEnvVarName = "CLUSTER_ID"
	NodeNameEnvVarName  = "NODE_NAME"
	PodIPEnvVarName     = "POD_IP"
)

func ClusterID() string {
	return os.Getenv(ClusterIDEnvVarName)
}

func NodeName() string {
	return os.Getenv(NodeNameEnvVarName)
}

func PodIP() string {
	return os.Getenv(PodIPEnvVarName)
}
