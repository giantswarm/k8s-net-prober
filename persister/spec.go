package persister

import (
	"context"

	"github.com/giantswarm/k8s-net-prober/types"
)

type Persister interface {
	Persist(ctx context.Context, result types.ProbeResult) error
}
