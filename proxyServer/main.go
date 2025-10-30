package main

import (
	"common/constants"
	"common/modules/db"
	"common/modules/db/redis"
	"flag"
	"fmt"
	"proxyServer/router"
	"proxyServer/service"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/acceptor"
	"github.com/topfreegames/pitaya/v2/component"
	"github.com/topfreegames/pitaya/v2/config"
)

var app pitaya.Pitaya

func main() {
	serverType := constants.ProxyServer

	port := flag.Int("port", 40000, "the port to listen")
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	config := config.NewDefaultBuilderConfig()
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(newAcceptor(*port)) //前端服务器必须指定一个端口，用于接收前端服务器的连接

	app = builder.Build()
	pitaya.DefaultApp = app

	/**
	停止 NATS（消息总线）连接；

	注销服务发现（etcd）；

	停止所有网络 acceptor（如 TCP、WebSocket 等）；

	停止正在运行的 goroutine；

	清理 metrics、tracer、logger 等资源。
	*/
	defer app.Shutdown()

	registerServices()
	registerModules()

	app.AddRoute(constants.LobbyServer, router.LobbyRouterFunc) //添加路由，用于处理前端服务器的请求

	app.Start()
}

func newAcceptor(port int) acceptor.Acceptor {
	tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", port))
	return tcp
}

func registerServices() {
	account := service.NewAccountService(app)
	app.Register(account,
		component.WithName("account"),
		component.WithNameFunc(strings.ToLower))
}

func registerModules() {
	r := redis.NewRedisStorage(redis.RedisConfig{
		Config: db.Config{
			Host:     "localhost",
			Port:     6379,
			Password: "",
		},
	})
	app.RegisterModule(r, constants.RedisModule)
}
