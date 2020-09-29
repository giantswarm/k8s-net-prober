package types

import "time"

type ProbeResult struct {
	Cluster string
	NodePool *string
	NodeName string
	ProbeType string
	Src string
	Dst string
	Success bool
	ProbeLengthMs int64
	Timestamp time.Time
}
