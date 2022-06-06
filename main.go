package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github/runzhliu/container-log-server/docs"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
)

type FileServerResponse []struct {
	Mtime string `json:"mtime"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Type  string `json:"type"`
}

var namespace = os.Getenv("POD_NAMESPACE")
var serverLabel = "app=file-server-ds"
var nodeField = "spec.nodeName="

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
	pod, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: serverLabel, FieldSelector: nodeField + nodeName})
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
// @Router /v1/log/download [get]
func download(c *gin.Context) {
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

// @Summary k8s日志清单接口
// @Description 注意参数格式
// @Accept json
// @Param host query string true "母机节点" default(10.9.70.1)
// @Param pod query string true "Pod名" default(test-a)
// @Param container query string true "容器名" default(test)
// @Success 200 object FileServerResponse
// @Success 404 object FileServerResponse
// @Router /v1/log/list [get]
func list(c *gin.Context) {
	host := c.Request.URL.Query().Get("host")
	pod := c.Request.URL.Query().Get("pod")
	container := c.Request.URL.Query().Get("container")

	fileServerPod := getFileServerPodIp(host)
	url := fmt.Sprintf("http://%s/%s/flume/%s/", fileServerPod, pod, container)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var rs FileServerResponse
	err = json.Unmarshal(body, &rs)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rs)
}

// @title container-log-server API
// @version 1.0
// @description k8s本地日志服务器
// @termsOfService http://swagger.io/terms/

// @contact.name runzhliu
// @contact.email runzhliu@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiV1 := r.Group("/v1/log")
	{
		apiV1.GET("/download", download)
		apiV1.GET("/list", list)
	}
	r.Run()
}
