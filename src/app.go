package app

import (
	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/server"
	"github.com/cqtrade/infobot/src/service"
)

var (
	cfg                   config.Config                    = config.New()
	discordService        service.DiscordService           = service.NewDiscordService(cfg)
	tradingviewController controller.TradingviewController = controller.NewTradingviewController(discordService)
	webServer             server.Server                    = server.New(cfg, tradingviewController)
)

func Run() {
	webServer.Run()
}
