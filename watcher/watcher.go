package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/k8s-net-prober/types"
)

const (
	namespace     = "default"
	daemonSetName = "k8s-net-prober"
)

type Config struct {
	Logger micrologger.Logger
}

type Interface struct {
	logger micrologger.Logger

	clientset *kubernetes.Clientset
}

func NewWatcher(config Config) (*Interface, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// setup configuration from env variables
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &Interface{
		logger: config.Logger,

		clientset: clientset,
	}, nil
}

func (w *Interface) Watch(ctx context.Context, c chan []types.PodInfo) error {
	for {
		destinations, err := w.getPods(ctx)
		if err != nil {
			return microerror.Mask(err)
		}

		c <- destinations

		time.Sleep(30 * time.Second)
	}
}

func (w *Interface) getPods(ctx context.Context) ([]types.PodInfo, error) {
	allPods, err := w.clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return []types.PodInfo{}, microerror.Mask(err)
	}

	var filteredPods []types.PodInfo

	w.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Found %d total pods in %s namespace", len(allPods.Items), namespace))

	for _, p := range allPods.Items {
		if p.Status.Phase != "Running" || p.Status.PodIP == "" {
			// Pod not ready yet.
			continue
		}
		for _, r := range p.ObjectMeta.OwnerReferences {
			if r.Kind == "DaemonSet" && r.Name == daemonSetName {
				filteredPods = append(filteredPods, types.PodInfo{
					IP:       p.Status.PodIP,
					NodeName: p.Spec.NodeName,
				})
				break
			}
		}
	}

	w.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Found %d interesting pods in %s namespace", len(filteredPods), namespace))

	return filteredPods, nil
}
