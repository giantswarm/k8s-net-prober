package persister

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/jackc/pgx/v4"

	"github.com/giantswarm/k8s-net-prober/env"
	"github.com/giantswarm/k8s-net-prober/types"
)

type PostgresPersisterConfig struct {
	Logger micrologger.Logger
}

type PostgresPersister struct {
	logger   micrologger.Logger
	pgclient *pgx.Conn
}

func NewPostgresPersister(config PostgresPersisterConfig) (*PostgresPersister, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	PostgresPersister := PostgresPersister{
		logger: config.Logger,
	}

	return &PostgresPersister, nil
}

func (p *PostgresPersister) Init(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, env.PostgresURL())
	if err != nil {
		return microerror.Mask(err)
	}

	p.pgclient = conn

	return nil
}

func (p *PostgresPersister) Persist(ctx context.Context, result types.ProbeResult) error {
	err := p.runQuery(ctx, "insert into results(cluster, ts, src_node, dst_node, success, duration_ms) values($1, $2, $3, $4, $5, $6)", result.Cluster, result.Timestamp, result.SrcPodInfo.NodeName, result.DstPodInfo.NodeName, result.Success, result.ProbeDurationMs)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (p *PostgresPersister) runQuery(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.pgclient.Exec(ctx, query, args...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
