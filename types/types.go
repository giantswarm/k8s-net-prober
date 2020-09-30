package types

import "time"

type PodInfo struct {
	IP       string `json:"ip"`
	NodeName string `json:"node_name"`
}

type ProbeResult struct {
	Cluster         string    `json:"cluster"`
	SrcPodInfo      PodInfo   `json:"src_pod_info"`
	DstPodInfo      PodInfo   `json:"dst_pod_info"`
	ProbeType       string    `json:"probe_type"`
	Success         bool      `json:"success"`
	ProbeDurationMs float64   `json:"probe_duration_ms"`
	Timestamp       time.Time `json:"timestamp"`
}
