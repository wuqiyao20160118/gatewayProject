package main

import (
	"flag"
	"os"
	"os/signal"
	"src/gatewayProject/dao"
	"src/gatewayProject/golang_common/lib"
	"src/gatewayProject/grpc_proxy_router"
	"src/gatewayProject/http_proxy_router"
	"src/gatewayProject/router"
	"src/gatewayProject/tcp_proxy_router"
	"syscall"
)

// 压测命令： ab -n requests -c concurrency url

/*
	Arguments:
		endpoint: dashboard / server
		config:   ./conf/prod/  ./conf/dev/
*/

// Package flag implements command-line flag parsing.
//
// Usage
//
// Define flags using flag.String(), Bool(), Int(), etc.
//
// This declares an integer flag, -flagname, stored in the pointer ip, with type *int.
// import "flag"
// var ip = flag.Int("flagname", 1234, "help message for flagname")
// If you like, you can bind the flag to a variable using the Var() functions.
// var flagvar int
// func init() {
//	 flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
// }
// Or you can create custom flags that satisfy the Value interface (with
// pointer receivers) and couple them to flag parsing by
// flag.Var(&flagVal, "name", "help message for flagname")
// For such flags, the default value is just the initial value of the variable.
//
// After all flags are defined, call
// flag.Parse()
// to parse the command line into the defined flags.
//
// Flags may then be used directly. If you're using the flags themselves,
// they are all pointers; if you bind to variables, they're values.
// fmt.Println("ip has value ", *ip)
// fmt.Println("flagvar has value ", flagvar)
//
// After parsing, the arguments following the flags are available as the
// slice flag.Args() or individually as flag.Arg(i).
// The arguments are indexed from 0 through flag.NArg()-1.

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
		_ = dao.ServiceManagerHandler.LoadOnce()
		_ = dao.AppManagerHandler.LoadOnce()

		go func() {
			http_proxy_router.HttpServerRun()
		}()
		go func() {
			http_proxy_router.HttpsServerRun()
		}()
		go func() {
			tcp_proxy_router.TcpServerRun()
		}()
		go func() {
			grpc_proxy_router.GrpcServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		tcp_proxy_router.TcpServerStop()
		grpc_proxy_router.GrpcServerStop()
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()

	}
}
