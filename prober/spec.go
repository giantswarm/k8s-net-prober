package prober

import (
	"github.com/giantswarm/k8s-net-prober/types"
)

type Prober interface {
	Probe(dst string) ProbeResult
}
