package persister

import (
	"context"

	"github.com/giantswarm/k8s-net-prober/types"
)

type Persister interface {
	Init(ctx context.Context) error
	Persist(ctx context.Context, result types.ProbeResult) error
}
