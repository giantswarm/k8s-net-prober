package prober

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/go-ping/ping"

	"github.com/giantswarm/k8s-net-prober/types"
)

const (
	probeType = "ping"
)

type PingProberConfig struct {
	Logger micrologger.Logger

	ClusterID   string
	Source      *types.PodInfo
	Destination *types.PodInfo
}

type PingProber struct {
	logger micrologger.Logger

	clusterID   string
	source      *types.PodInfo
	destination *types.PodInfo
	stopped     bool
}

func NewPingProber(config PingProberConfig) (*PingProber, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.ClusterID == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ClusterID must not be empty", config)
	}
	if config.Destination == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Destination must not be empty", config)
	}
	if config.Source == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Destination must not be empty", config)
	}

	PingProber := PingProber{
		logger: config.Logger,

		clusterID:   config.ClusterID,
		source:      config.Source,
		destination: config.Destination,
		stopped:     false,
	}

	return &PingProber, nil
}

func (p *PingProber) Start(ctx context.Context, collector chan types.ProbeResult) error {
	p.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("Starting probing %s", p.destination.IP))
	for {
		if p.stopped {
			return nil
		}

		pinger, err := ping.NewPinger(p.destination.IP)
		if err != nil {
			return microerror.Mask(err)
		}
		pinger.Timeout = 1 * time.Second
		pinger.SetPrivileged(true)
		pinger.Count = 1
		err = pinger.Run()
		if err != nil {
			return microerror.Mask(err)
		}

		stats := pinger.Statistics()

		res := types.ProbeResult{
			Cluster:       p.clusterID,
			SrcPodInfo:    *p.source,
			DstPodInfo:    *p.destination,
			ProbeType:     probeType,
			Success:       stats.PacketsRecv == stats.PacketsSent,
			ProbeLengthMs: float64(stats.AvgRtt.Microseconds() / 1000),
			Timestamp:     time.Now(),
		}

		collector <- res

		time.Sleep(1 * time.Second)
	}
}

func (p *PingProber) Stop(ctx context.Context) error {
	p.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("Stopping probing %s", p.destination.IP))

	p.stopped = true

	return nil
}
