FROM centos:centos7

WORKDIR /app
ADD log-server /app

ENTRYPOINT ["/app/log-server"]