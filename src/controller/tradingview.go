package tvcontroller

import (
	"strings"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/types"
	"github.com/gin-gonic/gin"
)

type TvController struct {
	cfg      config.Config
	notif    notification.Notification
	ftxTrade ftxtrade.FtxTrade
}

func New(cfg config.Config, notif notification.Notification, ftxTrade ftxtrade.FtxTrade) *TvController {
	return &TvController{
		cfg:      cfg,
		notif:    notif,
		ftxTrade: ftxTrade,
	}
}

func (tvc *TvController) PostText(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})
	text, err := ctx.GetRawData()
	if err != nil {
		println("error " + err.Error()) // TODO needs logger
		return
	}
	go tvc.notif.SendTextMessage(string(text))
}

func (tvc *TvController) PostFlash(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})

	var message types.JSONMessageBody
	err := ctx.ShouldBindJSON(&message)

	if err != nil {
		go tvc.notif.Log("ERROR", "ctx.ShouldBindJSON(&message)", err.Error())
		return
	}

	go tvc.notif.SendFlashMessage(message)

	switch message.Signal {
	case 1001:
		go func() {
			tvc.ftxTrade.BuyCoinBull("ethbull", "ETHBULL/USD")
			time.Sleep(time.Second)
			tvc.ftxTrade.BuyCoinBull("bull", "BULL/USD")
		}()
		return
	case -2002: // TP
		go func(ticker string) {
			t := strings.ToUpper(ticker)
			if strings.HasPrefix(t, "BTC") || strings.HasPrefix(t, "XBT") {
				tvc.ftxTrade.TpCoinBull("bull", "BULL/USD", "BULL")
			} else if strings.HasPrefix(t, "ETH") {
				tvc.ftxTrade.TpCoinBull("ethbull", "ETHBULL/USD", "ETHBULL")
			}
		}(message.Ticker)
		return

	// case 888:
	// 	go tvc.ftxTrade.ArbStart("arbbtc", "BTC")
	// 	return
	// case -888:
	// 	go tvc.ftxTrade.ArbEnd("arbbtc", "BTC")
	// 	return

	case 1, -1, 2, -2:
		go func(msg types.JSONMessageBody) {
			var side string
			var sideOpposite string
			if message.Signal == 1 { // enter_buy 1
				side = "buy"
				sideOpposite = "sell"
			} else if message.Signal == -1 { // enter_sell -1
				side = "sell"
				sideOpposite = "buy"
			} else if message.Signal == 2 { // exit_buy 2
				side = "exitBuy"
				sideOpposite = "sell"
			} else if message.Signal == -2 { // exit_sell -2
				side = "exitSell"
				sideOpposite = "buy"
			}

			tvc.ftxTrade.TradeLevCrypto(msg, tvc.cfg.RiskDC, side, sideOpposite, "dc")
			time.Sleep(time.Second)
			tvc.ftxTrade.TradeLevCrypto(msg, tvc.cfg.RiskD, side, sideOpposite, "d")

			ticker := msg.Ticker
			t := strings.ToUpper(ticker)
			if side == "buy" {
				time.Sleep(time.Second)
				tvc.ftxTrade.BuyCoinBull("ethbull", "ETHBULL/USD")
				time.Sleep(time.Second)
				tvc.ftxTrade.BuyCoinBull("bull", "BULL/USD")
			} else if side == "sell" {
				if strings.HasPrefix(t, "BTC") || strings.HasPrefix(t, "XBT") {
					tvc.notif.Log("INFO", "TODO EXIT crypto, ARB BTC,ETH?", message)
				}
			}

		}(message)
		return
	default:
		go tvc.notif.Log("ERROR", "unknown signal", message)
		return
	}
}
