package persister

import (
	"context"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/giantswarm/k8s-net-prober/env"
	"github.com/giantswarm/k8s-net-prober/types"
)

type PrometheusPersisterConfig struct {
	Logger micrologger.Logger
}

type PrometheusPersister struct {
	logger micrologger.Logger
}

var labels = []string{
	"cluster",
	"src_node",
	"dst_node",
	"src_ip",
	"dst_ip",
}

// Gauge for RTT.
var rttMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "k8s_net_prober_probe_duration_ms",
		Help: "The duration of the probe in milliseconds",
	},
	labels,
)

// Gauge for success/failure.
var successFailureMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "k8s_net_prober_probe_success",
		Help: "The succeeded status of the probe (1 = succeeded, 0 = failure)",
	},
	labels,
)

func NewPrometheusPersister(config PrometheusPersisterConfig) (*PrometheusPersister, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	prometheusPersister := PrometheusPersister{
		logger: config.Logger,
	}

	return &prometheusPersister, nil
}

func (p *PrometheusPersister) Init(ctx context.Context) error {
	prometheus.MustRegister(rttMetric)
	prometheus.MustRegister(successFailureMetric)

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		err := http.ListenAndServe(fmt.Sprintf(":%d", env.PrometheusListenPort()), nil)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (p *PrometheusPersister) Persist(ctx context.Context, result types.ProbeResult) error {
	rttMetric.WithLabelValues(result.Cluster, result.SrcPodInfo.NodeName, result.DstPodInfo.NodeName, result.SrcPodInfo.IP, result.DstPodInfo.IP).Set(result.ProbeDurationMs)

	successInt := float64(0)
	if result.Success {
		successInt = 1
	}
	successFailureMetric.WithLabelValues(result.Cluster, result.SrcPodInfo.NodeName, result.DstPodInfo.NodeName, result.SrcPodInfo.IP, result.DstPodInfo.IP).Set(successInt)

	return nil
}
