package app

import (
	"errors"

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
	ft           ftxtrade.FtxTrade         = *ftxtrade.New(cfg)
	ftws         ftxwebsocket.FtxWebSocket = *ftxwebsocket.New(cfg, appState)
	tvController tvcontroller.TvController = *tvcontroller.New(notif, ft)
	webServer    server.Server             = *server.New(cfg, tvController)
)

func Run() {
	/*
		go appState.StateLatestPrices()
		go appState.ReadPriceState()
		// go ftws.Start()
		go ftws.Start()
		// go ft.StartHealthPing()
		webServer.Run()
	*/

	type Person struct {
		name string
		age  int
	}
	var p Person
	p.name = "Motyar"
	p.age = 120

	notif.Log("error", nil, 2, "tere", 2.3, p, errors.New("OMG").Error())

}
