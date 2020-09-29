FROM alpine:3.8

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /opt/ignition

ADD ./k8s-net-prober /k8s-net-prober

ENTRYPOINT ["/k8s-net-prober"]
