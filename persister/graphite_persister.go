package persister

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/marpaia/graphite-golang"

	"github.com/giantswarm/k8s-net-prober/types"
)

type GraphitePersisterConfig struct {
	Logger micrologger.Logger
}

type GraphitePersister struct {
	logger   micrologger.Logger
	graphite *graphite.Graphite
}

func NewGraphitePersister(config GraphitePersisterConfig) (*GraphitePersister, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	GraphitePersister := GraphitePersister{
		logger: config.Logger,
	}

	return &GraphitePersister, nil
}

func (p *GraphitePersister) Init(ctx context.Context) error {
	graphiteClient, err := graphite.NewGraphite("graphite.azure.gigantic.io", 2003)
	if err != nil {
		return microerror.Mask(err)
	}

	p.graphite = graphiteClient
	return nil
}

func (p *GraphitePersister) Persist(ctx context.Context, result types.ProbeResult) error {
	label := fmt.Sprintf("clusters.%s.%s.%s", result.Cluster, result.SrcPodInfo.NodeName, result.DstPodInfo.NodeName)
	if result.Success {
		err := p.graphite.SimpleSend(fmt.Sprintf("%s.success", label), "1")
		if err != nil {
			return microerror.Mask(err)
		}
		err = p.graphite.SimpleSend(fmt.Sprintf("%s.duration", label), fmt.Sprintf("%f", result.ProbeDurationMs))
		if err != nil {
			return microerror.Mask(err)
		}
	} else {
		err := p.graphite.SimpleSend(fmt.Sprintf("%s.success", label), "0")
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
