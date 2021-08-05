package app

import (
	"github.com/cqtrade/infobot/src/config"
	tvcontroller "github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/ftxwebsocket"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/server"
	"github.com/cqtrade/infobot/src/state"
)

var (
	cfg          config.Config             = *config.New()
	notif        notification.Notification = *notification.New(cfg)
	appState     state.State               = *state.New(cfg, notif)
	ft           ftxtrade.FtxTrade         = *ftxtrade.New(cfg, notif, appState)
	ftws         ftxwebsocket.FtxWebSocket = *ftxwebsocket.New(cfg, appState)
	tvController tvcontroller.TvController = *tvcontroller.New(notif, ft)
	webServer    server.Server             = *server.New(cfg, tvController)
)

func Run() {
	// go appState.StateLatestPrices()
	// go ftws.Start()
	// go ft.StartHealthPing()
	go notif.RunStateLogMessages()
	go notif.RunReadLogMessages()
	go notif.Log("INFO", "Boot")
	webServer.Run()
}
