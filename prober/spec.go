package prober

import (
	"context"

	"github.com/giantswarm/k8s-net-prober/types"
)

type Prober interface {
	Start(ctx context.Context, collector chan types.ProbeResult) error
	Stop(ctx context.Context) error
}
