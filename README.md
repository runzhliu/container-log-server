## Overview

需求是当容器退出后，在一定的窗口期内，用户还可以通过某些方法拿到容器内的文件

## 分析

在目前的业务场景下，用户会将日志打到容器本地，不管日志文件是否有落到某种类型的卷中，同时由于日志采集的链路以及容量的问题，存在一些场景是用户希望可以用某些方法，在容器被 `kill` 或者重启后可以拿到原来的容器内的日志文件

### 排除的方法

1. Rancher 2.4商业版有提供文件下载的功能，Rancher 2.6社区版是没有文件下载的功能的，但是都不提供容器 `terminated` 之后的文件下载
2. 通过overlay2的文件系统去或者diff下的日志文件需要延缓container的删除

## 计划的方案

需要结合平台的运营的可运营性和用户使用的易用性考虑，前者是为了合理的清理母机的空间，后者是方便用户的使用

1. 通过webhook给业务Pod注入Volume/VolumeMount，以hostpath给用户提供日志存储的目录
2. 宿主机对应的目录必须固定且与系统根目录隔离，目录以namespace-workload-pod命令，提供定时删除的Job作为管理
3. 通过webhook注入LOG_PATH环境变量，提供已经挂载好的容器内目录作为默认的日志地址
4. 集群中部署file-server，用户可以根据Pod的属性host/workload/pod/container，通过HTTP协议下载日志文件

## 开发工作

1. webhook: 注入集群管理员的对象
2. file-server: 提供日志文件的下载
3. Job: 宿主机轮替/清理hostpath的定时任务
4. 下载文件接口: SRE这边只提供接口，接口参数包括host/pod/container/filename等，其他待定

## Notes

## Nginx配置文件

### Helm部署

```shell
helm install log-server -n log-server .
```

### 手动部署

## TODO

- [ ] fileserver工程
- [x] helm部署
- [ ] 容器内编译
