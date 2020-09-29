package env

import "os"

const (
	ClusterIDEnvVarName = "CLUSTER_ID"
)

func ClusterID() string {
	return os.Getenv(ClusterIDEnvVarName)
}
