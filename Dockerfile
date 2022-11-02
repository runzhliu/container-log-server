FROM centos:centos7

WORKDIR /app
ADD container-log-server /app

ENTRYPOINT ["/app/container-log-server"]