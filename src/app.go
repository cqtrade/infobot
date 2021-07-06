package app

import (
	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/server"
	"github.com/cqtrade/infobot/src/service"
)

var (
	cfg                   config.Config                    = config.New()
	discordService        service.DiscordService           = service.NewDiscordService(cfg)
	tradingviewController controller.TradingviewController = controller.NewTradingviewController(discordService)
	webServer             server.Server                    = server.New(cfg, tradingviewController)
)

/*
TODO
HEALTH Logger to Discord

https://github.com/go-numb/go-ftx
https://github.com/grishinsana/goftx
https://github.com/cloudingcity/go-ftx
*/

func Run() {
	ftxtrade.StartStuff()

	// webServer.Run()
}
