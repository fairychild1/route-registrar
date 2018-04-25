package main

import (
	"flag"
	"strconv"
	//"io/ioutil"
	"log"
	"os"
	"fmt"
	//"os/signal"
	//"strconv"
	//"syscall"
	"code.cloudfoundry.org/lager"
	//"code.cloudfoundry.org/lager/lagerflags"
	"route-registrar/config"
	//"route-registrar/healthchecker"
	"route-registrar/messagebus"
	//"route-registrar/registrar"
	//"github.com/tedsuo/ifrit"
	_ "route-registrar/routers"
	"github.com/astaxie/beego"
	"route-registrar/controllers"
	"route-registrar/utility"
	"route-registrar/models"

)

//func main() {
//	var configPath string
//	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
//
//	pidfile := flags.String("pidfile", "", "Path to pid file")
//	lagerflags.AddFlags(flags)
//
//	flags.StringVar(&configPath, "configPath", "", "path to configuration file with json encoded content")
//	flags.Set("configPath", "registrar_settings.yml")
//
//	flags.Parse(os.Args[1:])
//
//	logger, _ := lagerflags.New("Route Registrar")
//
//	logger.Info("Initializing")
//
//	configSchema, err := config.NewConfigSchemaFromFile(configPath)
//	if err != nil {
//		logger.Fatal("error parsing file: %s\n", err)
//	}
//
//	c, err := configSchema.ToConfig()
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	hc := healthchecker.NewHealthChecker(logger)
//
//	logger.Info("creating nats connection")
//	messageBus := messagebus.NewMessageBus(logger)
//
//	r := registrar.NewRegistrar(*c, hc, logger, messageBus)
//
//	if *pidfile != "" {
//		pid := strconv.Itoa(os.Getpid())
//		err := ioutil.WriteFile(*pidfile, []byte(pid), 0644)
//		logger.Info("Writing pid", lager.Data{"pid": pid, "file": *pidfile})
//		if err != nil {
//			logger.Fatal(
//				"error writing pid to pidfile",
//				err,
//				lager.Data{
//					"pid":     pid,
//					"pidfile": *pidfile,
//				},
//			)
//		}
//	}
//
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
//
//	logger.Info("Running")
//
//	process := ifrit.Invoke(r)
//	for {
//		select {
//		case s := <-sigChan:
//			logger.Info("Caught signal", lager.Data{"signal": s})
//			process.Signal(s)
//		case err := <-process.Wait():
//			if err != nil {
//				logger.Fatal("Exiting with error", err)
//			}
//			logger.Info("Exiting without error")
//			os.Exit(0)
//		}
//	}
//}

func main() {
	var port int
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("can't find the PORT parameter")
		return
	}

	beego.BConfig.Listen.HTTPPort = port
	//beego.BConfig.RunMode = "dev"
	beego.BConfig.CopyRequestBody = true



	var configPath string
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&configPath, "configPath", "", "path to configuration file with json encoded content")
	flags.Set("configPath", "registrar_settings.yml")
	flags.Parse(os.Args[1:])
	logger := lager.NewLogger("Route Registrar")
	utility.Logger = logger
	logger.Info("Initializing")
	configSchema, err := config.NewConfigSchemaFromFile(configPath)
	if err != nil {
		logger.Fatal("error parsing file: %s\n", err)
	}
	c, err := configSchema.ToConfig()
	if err != nil {
		log.Fatalln(err)
	}
	utility.Config = *c
	//连接msgbus
	utility.Msgbus = messagebus.NewMessageBus(logger)
	err = utility.Msgbus.Connect(utility.Config.MessageBusServers)
	if err != nil {
		fmt.Printf("连接不上nats，请检查配置文件的nats信息是否正确")
		//return
	}

	//读取数据库的route信息，如果有，则启动协程，进行路由注册
	if proutes,err := models.LoadAllRoutes();err != nil {
		fmt.Printf("读取数据库的route信息失败，请检查连接\n")
		return
	}else {
		for _,v := range *proutes {
			if err := controllers.AddRoutes(v["RouteName"].(string),v["Port"].(string),v["Uri"].(string),v["HostIp"].(string));err !=nil {
				return
			}
			//r := controllers.NewRoute(v["RouteName"].(string),int(v["Port"].(int64)),v["Uri"].(string))
			//utility.AddRouteChan(v["HostIp"].(string)+":"+utility.Int64ToString(v["Port"].(int64)))
			//go utility.RegisterRoute(v["HostIp"].(string),utility.Int64ToString(v["Port"].(int64)),*r)
		}
	}


	beego.InsertFilter("/v1/*", beego.BeforeRouter, controllers.FilterToken)
	beego.Router("/login/", &controllers.LoginController{}, "post:Login")
	beego.Router("/v1/route_register/", &controllers.RouteController{}, "post:RouteRegister")
	beego.Router("/v1/route_deregister/:route_id", &controllers.RouteController{}, "delete:RouteDeregister")
	beego.Router("/v1/list", &controllers.RouteController{}, "get:ListRoutes")
	beego.Router("/v1/route", &controllers.RouteController{}, "get:GetRoute")
	beego.Router("/v1/route", &controllers.RouteController{}, "put:UpdateRoute")
	beego.Router("/v1/checkuri", &controllers.RouteController{}, "post:CheckUri")
	beego.SetStaticPath("/views", "views")
	beego.Run()
}
