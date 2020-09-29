package main

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/k8s-net-prober/env"
	"github.com/giantswarm/k8s-net-prober/persister"
	"github.com/giantswarm/k8s-net-prober/prober"
	"github.com/giantswarm/k8s-net-prober/types"
	"github.com/giantswarm/k8s-net-prober/watcher"
)

func main() {
	err := mainError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainError() error {
	ctx := context.Background()

	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Init persisters.
	var persisters []persister.Persister
	{
		logPersister, err := persister.NewLogPersister(persister.LogPersisterConfig{
			Logger: logger,
		})
		if err != nil {
			return microerror.Mask(err)
		}
		persisters = append(persisters, logPersister)
	}

	w, err := watcher.NewWatcher(watcher.Config{
		Logger: logger,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	// Create channel to retrieve destinations.
	destinationsWatcher := make(chan []types.PodInfo)

	go w.Watch(ctx, destinationsWatcher)

	// Create channel to retrieve probe results from probers.
	ch := make(chan types.ProbeResult)

	probers := map[string]prober.Prober{}

	// TODO populare src info.
	source := types.PodInfo{
		IP:       "",
		NodePool: nil,
		NodeName: "",
	}

	go func() {
		for {
			destinations := <-destinationsWatcher

			for _, d := range destinations {
				// Check if probe is already running or start it.
				_, found := probers[d.IP]
				if !found {
					pingProber, err := prober.NewPingProber(prober.PingProberConfig{
						Logger:      logger,
						ClusterID:   env.ClusterID(),
						Source:      &source,
						Destination: &d,
					})
					if err != nil {
						panic(err)
					}

					probers[d.IP] = pingProber

					go pingProber.Start(ctx, ch)
				}
			}

			// Check if any probe has to be stopped.
			for d, _ := range probers {
				if !inSlice(d, destinations) {
					probers[d].Stop(ctx)
					delete(probers, d)
				}
			}
		}
	}()

	// Send all probe results to all persisters.
	for {
		res := <-ch
		for _, p := range persisters {
			err = p.Persist(ctx, res)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}
}

func inSlice(needle string, haystack []types.PodInfo) bool {
	for _, s := range haystack {
		if needle == s.IP {
			return true
		}
	}

	return false
}
