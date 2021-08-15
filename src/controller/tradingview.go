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
	case -2002:
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
		go tvc.ftxTrade.ArbStart("arb1", "BTC")
		return
	case -888:
		go tvc.ftxTrade.ArbEnd("arb1", "BTC")
		return
	case 1, -1, 2, -2:
		go tvc.ftxTrade.TradeLev(message)
		// todo low equity high risk accounts?
		// exit signals for spot from here
		return
	default:
		go tvc.notif.Log("ERROR", "unknown signal", message)
		return
	}
}
