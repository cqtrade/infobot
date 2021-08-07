package ftxtrade

import (
	"fmt"
	"math"
	"time"

	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
)

func (ft *FtxTrade) checkPosition(client *ftx.FtxClient, market string) {
	position, err := ft.CheckFuturePosition(client, market)
	if err != nil {
		ft.notif.Log("ERROR", "checkPosition", err.Error())
		return
	}

	openTriggerOrders, err := client.GetOpenTriggerOrders(market)
	if err != nil {
		ft.notif.Log("ERROR", "checkPosition GetOpenTriggerOrders", err.Error())
		return
	} else if !openTriggerOrders.Success {
		ft.notif.Log("ERROR", "checkPosition GetOpenTriggerOrders UNSUFFESSFUL", openTriggerOrders.Result)
		return
	} else {

	}

	triggerOrdersLength := len(openTriggerOrders.Result)
	if position.Size == 0 && triggerOrdersLength == 0 {
		return
	}

	if position.Size == 0 && triggerOrdersLength > 0 {
		ft.notif.Log("INFO", "checkPosition", "no position, open trigger orders. cancel all open orders.")
		res, err := client.CancelAllOrders()
		if err != nil {
			ft.notif.Log("ERROR", "checkPosition CancelAllOrders", err.Error())
			return
		}
		if !res.Success {
			ft.notif.Log("ERROR", "checkPosition CancelAllOrders UNSUCCESSFUL", res.Result)
			return
		}
		return
	}

	// ft.notif.Log("", "position", position)
	var slOrder structs.TriggerOrder
	// var tpOrder structs.TriggerOrder
	for _, triggerOrder := range openTriggerOrders.Result {
		if triggerOrder.Type == "stop" {
			slOrder = triggerOrder
		}
		// else if triggerOrder.Type == "take_profit" {
		// 	tpOrder = triggerOrder
		// }
	}

	// ft.notif.Log("", "tpOrder", tpOrder)

	price, err := ft.appState.ReadLatestPriceForMarket(market)

	if err != nil {
		ft.notif.Log("ERROR", "checkPosition ReadLatestPriceForMarket. Abort.", err.Error())
		return
	}

	ft.notif.Log("", "latest price", market, price)

	if slOrder.Status != "open" {
		ft.notif.Log("ERROR", "GetOpenTriggerOrders Open position no SL")
	}
	// ft.notif.Log("", "slOrder", slOrder)
	diffAllowed := 0.0001
	diff := math.Abs((position.AverageOpenPrice / slOrder.TriggerPrice) - 1)
	fmt.Println("position.AverageOpenPrice", position.AverageOpenPrice)
	fmt.Println("slOrder.TriggerPrice", slOrder.TriggerPrice, slOrder.Size)
	fmt.Println("diff", diff)
	if position.Size < slOrder.Size && diff > diffAllowed {
		slOrder, err := client.ModifyTriggerOrder(slOrder.ID, position.Size, position.AverageOpenPrice)
		if err != nil {
			ft.notif.Log("ERROR", "checkPosition New SL.", err.Error())
			return
		}
		if !slOrder.Success {
			ft.notif.Log("ERROR", "checkPosition New SL UNSUCCESSFUL.", slOrder.Result)
			return
		}
		ft.notif.Log("INFO", "checkPosition SUCCESS moving SL.")
	}
}

func (ft *FtxTrade) RunPositionsCheck() {
	time.Sleep(time.Second * 3)
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	// clientBTCD := ftx.New(key, secret, ft.cfg.SubAccBTCD)
	clientETHD := ftx.New(key, secret, ft.cfg.SubAccETHD)
	for {
		// 0.001

		// ft.checkPosition(clientBTCD, ft.cfg.FutureBTC)
		// time.Sleep(time.Second)
		ft.checkPosition(clientETHD, ft.cfg.FutureETH)
		time.Sleep(time.Second * 15)
	}
}
