module github.com/giantswarm/k8s-net-prober

go 1.15

require (
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.3
	github.com/go-ping/ping v0.0.0-20200918120429-e8ae07c3cec8
	github.com/jackc/pgx/v4 v4.9.0
	github.com/marpaia/graphite-golang v0.0.0-20190519024811-caf161d2c2b1
	github.com/prometheus/client_golang v1.7.1
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
)

replace github.com/gorilla/websocket v0.0.0-20170926233335-4201258b820c => github.com/gorilla/websocket v1.4.2
