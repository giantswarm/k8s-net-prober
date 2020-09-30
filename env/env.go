package env

import (
	"os"
	"strconv"
)

const (
	ClusterIDEnvVarName            = "CLUSTER_ID"
	NodeNameEnvVarName             = "NODE_NAME"
	PodIPEnvVarName                = "POD_IP"
	PostgresURLEnvVar              = "POSTGRES_URL"
	PrometheusListenPortEnvVarName = "PROMETHEUS_LISTEN_PORT"
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

func PostgresURL() string {
	return os.Getenv(PostgresURLEnvVar)
}

func PrometheusListenPort() int {
	str := os.Getenv(PrometheusListenPortEnvVarName)
	i, err := strconv.Atoi(str)
	if err != nil {
		return 9339
	}

	return i
}
