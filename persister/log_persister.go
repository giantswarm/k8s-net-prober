package persister

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/k8s-net-prober/types"
)

type LogPersisterConfig struct {
	Logger micrologger.Logger
}

type LogPersister struct {
	logger micrologger.Logger
}

func NewLogPersister(config LogPersisterConfig) (*LogPersister, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	logPersister := LogPersister{
		logger: config.Logger,
	}

	return &logPersister, nil
}

func (l *LogPersister) Init(ctx context.Context) error {
	return nil
}

func (l *LogPersister) Persist(ctx context.Context, result types.ProbeResult) error {
	l.logger.LogCtx(ctx, "level", "info", "message", result)
	return nil
}
