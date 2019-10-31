package main

import (
	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"net/http"
	"runtime"
	"strings"
	"vlog/controllers"
	"vlog/server"
)

func Initialize() {
	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	viper.SetConfigFile("./config/App.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件出错")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if strings.Contains(e.Name, "App.json") {
			log.Info("Reload App.json Success")
		}
	})

	logger, err := log.LoggerFromConfigAsFile("./config/seelog.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)
	defer log.Flush()

	log.Info("cpu: ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

}

func main() {
	Initialize()

	server := server.Server{
		Controller: &controllers.MainController{},
	}

	server.Run(":8088")
}
