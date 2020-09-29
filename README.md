[![CircleCI](https://circleci.com/gh/giantswarm/k8s-net-prober.svg?style=shield)](https://circleci.com/gh/giantswarm/k8s-net-prober) [![Docker Repository on Quay](https://quay.io/repository/giantswarm/k8s-net-prober/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/k8s-net-prober)

# k8s-net-prober

This project aims at testing kubernetes clusters' CNI network health.
It is intended to be run on test clusters and NOT on real clusters as it is heavy in terms of generated data.

It is meant to be deployed as a `Daemonset`, and every Pod keeps executing an ICMP Ping every second towards every other 
Pod in the Daemonset (including itself).

The results of such pings (success/failure and the RTT) are stored in a SQL database.
