package watcher

import (
	"context"
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
}

func NewWatcher(config Config) (*Interface, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	return &Interface{
		logger: config.Logger,
	}, nil
}

func (w *Interface) Watch(ctx context.Context, c chan []types.PodInfo) error {
	for {
		destinations, err := getPods(ctx)
		if err != nil {
			return microerror.Mask(err)
		}

		c <- destinations

		time.Sleep(30 * time.Second)
	}
}

func getPods(ctx context.Context) ([]types.PodInfo, error) {
	// setup configuration from env variables
	config, err := rest.InClusterConfig()
	if err != nil {
		return []types.PodInfo{}, microerror.Mask(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return []types.PodInfo{}, microerror.Mask(err)
	}

	allPods, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return []types.PodInfo{}, microerror.Mask(err)
	}

	var filteredPods []types.PodInfo

	for _, p := range allPods.Items {
		if p.Status.Phase != "Running" || p.Status.PodIP == "" {
			// Pod not ready yet.
			continue
		}
		for _, r := range p.ObjectMeta.OwnerReferences {
			if r.Kind == "DaemonSet" && r.Name == daemonSetName {
				filteredPods = append(filteredPods, types.PodInfo{
					IP: p.Status.PodIP,
					// TODO fill up node pool name.
					NodePool: nil,
					NodeName: p.Spec.NodeName,
				})
			}
		}
	}

	return filteredPods, nil
}
