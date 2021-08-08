package ftxtrade

import (
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
		ft.notif.Log("ERROR", "checkPosition GetOpenTriggerOrders", market, err.Error())
		return
	} else if !openTriggerOrders.Success {
		ft.notif.Log("ERROR", "checkPosition GetOpenTriggerOrders UNSUFFESSFUL", market, openTriggerOrders.Result)
		return
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

	var slOrder structs.TriggerOrder
	lenghtOfStopOrders := 0
	for _, triggerOrder := range openTriggerOrders.Result {
		if triggerOrder.Type == "stop" {
			slOrder = triggerOrder
			lenghtOfStopOrders++
		}
	}

	if lenghtOfStopOrders > 1 {
		ft.notif.Log("ERROR", market, "lenghtOfStopOrders > 1", lenghtOfStopOrders)
	}

	price, err := ft.appState.ReadLatestPriceForMarket(market)

	if err != nil {
		ft.notif.Log("ERROR", "checkPosition ReadLatestPriceForMarket. Abort.", err.Error())
		return
	}

	if slOrder.Status != "open" {
		ft.notif.Log("ERROR", "GetOpenTriggerOrders Open position no SL")
	}

	diffAllowed := 0.0001
	diff := math.Abs((position.AverageOpenPrice / slOrder.TriggerPrice) - 1)

	// ft.notif.Log("", "latest price", market, price)
	// fmt.Println("position.AverageOpenPrice", position.AverageOpenPrice)
	// fmt.Println("slOrder.TriggerPrice", slOrder.TriggerPrice, slOrder.Size)
	// fmt.Println("diff", diff)

	if position.Size < slOrder.Size && diff > diffAllowed {
		newSLtriggerPrice := position.AverageOpenPrice
		if position.Side == "buy" && newSLtriggerPrice > price {
			ft.notif.Log("ERROR", "checkPosition buy can't move SL to breakeven SL is higher than entry", newSLtriggerPrice, price)
			return
		}
		if position.Side == "sell" && newSLtriggerPrice < price {
			ft.notif.Log("ERROR", "checkPosition sell can't move SL to breakeven SL is lower than entry", newSLtriggerPrice, price)
			return
		}
		slOrder, err := client.ModifyTriggerOrder(slOrder.ID, position.Size, newSLtriggerPrice)
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
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	clientBTCD := ftx.New(key, secret, ft.cfg.SubAccBTCD)
	clientETHD := ftx.New(key, secret, ft.cfg.SubAccETHD)
	for {
		time.Sleep(time.Second * 15)
		ft.checkPosition(clientBTCD, ft.cfg.FutureBTC)
		time.Sleep(time.Second)
		ft.checkPosition(clientETHD, ft.cfg.FutureETH)
	}
}
