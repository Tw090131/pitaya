package main

import (
	"common/constants"
	"common/modules/db"
	"common/modules/db/mongodb"
	"lobbyServer/service"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"github.com/topfreegames/pitaya/v2/config"
)

var app pitaya.Pitaya

func main() {
	serverType := constants.LobbyServer

	logrus.SetLevel(logrus.DebugLevel)

	config := config.NewDefaultBuilderConfig()
	builder := pitaya.NewDefaultBuilder(false, serverType, pitaya.Cluster, map[string]string{}, *config)

	app = builder.Build()
	pitaya.DefaultApp = app

	defer app.Shutdown()

	registerServices()
	registerModules()

	app.Start()
}

func registerServices() {
	account := service.NewAccountService(app)
	app.Register(account,
		component.WithName("account"),
		component.WithNameFunc(strings.ToLower))
	app.RegisterRemote(account,
		component.WithName("account"),
		component.WithNameFunc(strings.ToLower))
}

func registerModules() {
	// TODO: 测试中 直接写死 后续需改成读配置文件
	mongo := mongodb.NewMongoStorage(mongodb.MongoConfig{
		Config: db.Config{
			Host:     "localhost",
			Port:     20777,
			Username: "",
			Password: "",
		},
		MaxPoolSize: 10,
	})
	app.RegisterModule(mongo, constants.MongoDBModule)
}
