package tvcontroller

import (
	"strings"
	"time"

	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/types"
	"github.com/gin-gonic/gin"
)

type TvController struct {
	notif    notification.Notification
	ftxTrade ftxtrade.FtxTrade
}

func New(notif notification.Notification, ftxTrade ftxtrade.FtxTrade) *TvController {
	return &TvController{
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
	case 888:
		go tvc.ftxTrade.ArbStart("arbbtc", "BTC")
		return
	case -888:
		go tvc.ftxTrade.ArbEnd("arbbtc", "BTC")
		return

		// enter_buy 1
		// enter_sell -1
		// exit_buy 2
		// exit_sell -2

	case 1, -1, 2, -2:
		go func(msg types.JSONMessageBody) {
			tvc.ftxTrade.TradeLev(msg)
			ticker := msg.Ticker
			t := strings.ToUpper(ticker)
			if message.Signal == 1 { // enter_buy 1
				time.Sleep(time.Second)
				tvc.ftxTrade.BuyCoinBull("ethbull", "ETHBULL/USD")
				time.Sleep(time.Second)
				tvc.ftxTrade.BuyCoinBull("bull", "BULL/USD")
			} else if message.Signal == 2 { // exit_buy 2
				if strings.HasPrefix(t, "BTC") || strings.HasPrefix(t, "XBT") {
					time.Sleep(time.Second)
					tvc.ftxTrade.TpCoinBull("bull", "BULL/USD", "BULL")
					time.Sleep(time.Second)
					tvc.ftxTrade.TpCoinBull("ethbull", "ETHBULL/USD", "ETHBULL")
				} else if strings.HasPrefix(t, "ETH") {
					time.Sleep(time.Second)
					tvc.ftxTrade.TpCoinBull("ethbull", "ETHBULL/USD", "ETHBULL")
				}
			} else if message.Signal == -1 { // TODO enter_sell -1 xbt
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
