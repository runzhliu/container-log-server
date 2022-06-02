package main

import (
	"context"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github/runzhliu/log-server/docs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

func getFileServerPodIp(nodeName string) string {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pod, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{LabelSelector: "app=file-server-ds", FieldSelector: "spec.nodeName=" + nodeName})
	if err != nil {
		panic(err.Error())
	} else {
		return pod.Items[0].Status.PodIP
	}
}

// @Summary k8s日志下载接口
// @Description 注意参数格式
// @Accept json
// @Param host query string true "母机节点" default(10.9.70.1)
// @Param pod query string true "Pod名" default(test-a)
// @Param container query string true "容器名" default(test)
// @Param log query string true "日志名" default(test.log)
// @Router /v1/log [get]
func run(c *gin.Context) {
	host := c.Request.URL.Query().Get("host")
	pod := c.Request.URL.Query().Get("pod")
	container := c.Request.URL.Query().Get("container")
	log := c.Request.URL.Query().Get("log")

	fileServerPod := getFileServerPodIp(host)
	url := fmt.Sprintf("http://%s/%s/flume/%s/%s", fileServerPod, pod, container, log)

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", url)

	resp := client.Do(req)

	// set header
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.Filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	c.File(fmt.Sprintf("/app/%s", resp.Filename))
}

// @title log-server API
// @version 1.0
// @description k8s本地日志服务器
// @termsOfService http://swagger.io/terms/

// @contact.name drummerliu
// @contact.email runzhliu@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiV1 := r.Group("/v1")
	{
		apiV1.GET("/log", run)
	}
	r.Run()
}
