package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gcfg.v1"
	"os"
	"os/exec"
	"path"
	"runtime"
	"seeta_campus/action"
	"strconv"
	"strings"
)

func setupRouter(debug bool) *gin.Engine {

	if debug == false {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.POST("/seeta/guest", recognition.Guest)
	router.POST("/seeta/insert", recognition.Insert)
	router.POST("/seeta/findOne", recognition.FindOne)
	router.POST("/seeta/find", recognition.Find)
	router.POST("/seeta/update", recognition.Update)
	router.POST("/seeta/image", recognition.Image)
	router.POST("/seeta/accuracy", recognition.Accuracy)
	return router
}
func execShell(s string) int {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return -1
	}
	// 去除空格
	str := strings.Replace(out.String(), " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	number, _ := strconv.Atoi(str)
	return number
}
func main() {

	//get the path
	_, fulleFilename, _, _ := runtime.Caller(0)
	realDir := path.Dir(fulleFilename)
	filename := os.Args[0]
	file := strings.Split(filename, "/")

	num := len(file) - 1
	conf := realDir + "/" + file[num] + ".ini"
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		fmt.Printf("配置文件：%s.ini 不存在\n\n", file[num])
		return
	}
	config := struct {
		Section struct {
			Enabled bool
			Debug   bool
			Port    int
		}
		MongoDB struct {
			Ip   string
			Port int
		}
		Seeta struct {
			Httpurl  string
			Httpport int
			Tcpport  int
		}
	}{}

	//get the configuration
	err := gcfg.ReadFileInto(&config, file[num]+".ini")
	if err != nil {
		fmt.Printf("Failed to parse config file: %s\n", err)
		return
	}
	if config.Section.Enabled == false {
		fmt.Println("请开始服务!")
	} else {
		port := config.Section.Port

		cmd := "netstat -an|grep " + strconv.Itoa(port) + "|wc -l"
		status := execShell(cmd)
		if status == -1 {
			fmt.Println("查询端口失败，请重新运行!")
			return
		} else if status >= 1 {
			fmt.Printf("端口号：%d 被占用，请配置未使用端口\n", port)
			return
		}

		httpurl := ":" + strconv.Itoa(port)
		//start the server
		router := setupRouter(config.Section.Debug)
		router.Run(httpurl)
	}
}
