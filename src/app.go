package app

import (
	"github.com/cqtrade/infobot/src/config"
	tvcontroller "github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/ftxwebsocket"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/server"
)

var (
	cfg          config.Config             = *config.New()
	notif        notification.Notification = *notification.New(cfg)
	ft           ftxtrade.FtxTrade         = *ftxtrade.New(cfg)
	ftws         ftxwebsocket.FtxWebSocket = *ftxwebsocket.New(cfg, ft, notif)
	tvController tvcontroller.TvController = *tvcontroller.New(notif, ft)
	webServer    server.Server             = *server.New(cfg, tvController)
)

func Run() {
	// go ftws.StateLatestPrices()
	// go ftws.ReadPriceState()
	// go ftws.Start()
	// go ft.StartHealthPing()
	webServer.Run()
}
