package app

import (
	"github.com/cqtrade/infobot/src/config"
	tvcontroller "github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/ftxwebsocket"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/server"
	"github.com/cqtrade/infobot/src/state"
	"github.com/cqtrade/infobot/src/tasignals"
)

var (
	cfg          config.Config             = *config.New()
	notif        notification.Notification = *notification.New(cfg)
	appState     state.State               = *state.New(cfg, notif)
	ft           ftxtrade.FtxTrade         = *ftxtrade.New(cfg, notif, appState)
	ftws         ftxwebsocket.FtxWebSocket = *ftxwebsocket.New(cfg, notif, appState)
	tvController tvcontroller.TvController = *tvcontroller.New(cfg, notif, ft)
	tasigs       tasignals.TaSignals       = *tasignals.New(cfg, tvController)
	webServer    server.Server             = *server.New(cfg, tvController)
)

func Run() {
	// tasigs.Indies()
	// go notif.RunStateLogMessages()
	// go notif.RunReadLogMessages()
	// go appState.RunStateLatestPrices()
	// go appState.RunPositionsInfo()
	// go ftws.RunWebSocket()
	// go ft.RunHealthPing()
	// go notif.Log("INFO", "Boot")
	// go ft.RunPositionsCheck()
	// // ft.Portfolio("p")
	// webServer.Run()

	// tasigs.CheckFlashSignals()
	// go notif.RunStateLogMessages()
	// go notif.RunReadLogMessages()
	// go appState.RunStateLatestPrices()
	// go appState.RunPositionsInfo()
	// go ftws.RunWebSocket()
	// go ft.RunHealthPing()
	// go notif.Log("INFO", "Boot")
	// go ft.RunPositionsCheck()
	webServer.Run()

	// ft.PortfolioFTX("portfolio")
}
