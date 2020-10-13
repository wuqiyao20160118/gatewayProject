package main

import (
	"flag"
	"os"
	"os/signal"
	"src/gatewayProject/dao"
	"src/gatewayProject/golang_common/lib"
	"src/gatewayProject/http_proxy_router"
	"src/gatewayProject/router"
	"syscall"
)

// 压测命令： ab -n requests -c concurrency url

/*
	Arguments:
		endpoint: dashboard / server
		config:   ./conf/prod/  ./conf/dev/
*/

var (
	endpoint = flag.String("endpoint", "", "input endpoint: dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		lib.InitModule(*config)
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		lib.InitModule(*config)
		defer lib.Destroy()
		// 服务启动时加载服务配置列表
		dao.ServiceManagerHandler.LoadOnce()

		go func() {
			http_proxy_router.HttpServerRun()
			http_proxy_router.HttpsServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}
}
