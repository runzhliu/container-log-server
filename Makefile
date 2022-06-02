all: push

swag:
	@echo ">> swag init 初始化接口文档"
	swag init

build: swag
	@echo ">> building code"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

docker: build
	@echo ">> building docker image"
	docker build -t runzhliu/container-log-server:latest .

push: docker
	@echo ">> pushing docker image"
	docker push runzhliu/container-log-server:latest

.PHONY: all build docker swag
