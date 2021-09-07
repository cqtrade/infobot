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
	ftws         ftxwebsocket.FtxWebSocket = *ftxwebsocket.New(cfg, notif, appState)
	tvController tvcontroller.TvController = *tvcontroller.New(cfg, notif, ft)
	webServer    server.Server             = *server.New(cfg, tvController)
)

func Run() {
	// key := cfg.FTXKey
	// secret := cfg.FTXSecret
	// client := ftx.New(key, secret, "")
	// defer client.Client.CloseIdleConnections()

	// // candles, err := client.GetHistoricalPriceLatest("BULL/USD", 15, 1)
	// // candles, err := client.GetHistoricalPriceLatest("BULL/USD", 15, 1)
	// candles, err := client.GetPositions(true)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// if !candles.Success {
	// 	fmt.Println(fmt.Sprintf("HERE1 %d ", candles.HTTPCode) + candles.ErrorMessage)
	// } else {
	// 	fmt.Println(fmt.Sprintf("HERE2 %+v", candles))
	// }

	go notif.RunStateLogMessages()
	go notif.RunReadLogMessages()
	go appState.RunStateLatestPrices()
	go appState.RunPositionsInfo()
	go ftws.RunWebSocket()
	go ft.RunHealthPing()
	go notif.Log("INFO", "Boot")
	go ft.RunPositionsCheck()
	webServer.Run()

}
