package ftxtrade

import (
	"math"
	"time"

	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
	"github.com/cqtrade/infobot/src/types"
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

	var subAcc string
	if client.Subaccount == "" {
		subAcc = "main"
	} else {
		subAcc = client.Subaccount
	}

	var sidePos string
	if position.Size > 0 {
		sidePos = position.Side
	} else {
		sidePos = ""
	}

	writePositionsInfo := types.WritePositionsInfo{
		Key: subAcc + "_" + market,
		PositionInfo: types.PositionInfo{
			Side:        sidePos,
			Stops:       0,
			TakeProfits: 0,
		},
		Resp: make(chan bool),
	}

	if position.Size == 0 && triggerOrdersLength == 0 {
		ft.appState.PositionsInfoWrites <- writePositionsInfo
		<-writePositionsInfo.Resp
		return
	}

	var slOrder structs.TriggerOrder
	lenghtOfStopOrders := 0
	lenghtOfTpOrders := 0
	for _, triggerOrder := range openTriggerOrders.Result {
		if triggerOrder.Type == "stop" {
			slOrder = triggerOrder
			lenghtOfStopOrders++
		} else if triggerOrder.Type == "take_profit" {
			lenghtOfTpOrders++
		}
	}

	writePositionsInfo.PositionInfo.Stops = lenghtOfStopOrders
	writePositionsInfo.PositionInfo.TakeProfits = lenghtOfTpOrders
	ft.appState.PositionsInfoWrites <- writePositionsInfo
	<-writePositionsInfo.Resp

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
		ft.notif.Log("INFO", "checkPosition SUCCESS moving SL:", slOrder.Result.Side, "@", slOrder.Result.TriggerPrice, "Position:", position.Side, "@", position.AverageOpenPrice)
	}
}

func (ft *FtxTrade) RunPositionsCheck() {
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret

	clientBTCD := ftx.New(key, secret, ft.cfg.SubAccBTCD)
	clientETHD := ftx.New(key, secret, ft.cfg.SubAccETHD)
	clientBTCDC := ftx.New(key, secret, ft.cfg.SubAccBTCDC)
	clientETHDC := ftx.New(key, secret, ft.cfg.SubAccETHDC)

	for {
		time.Sleep(time.Second * 10)
		ft.checkPosition(clientBTCD, ft.cfg.FutureBTC)
		time.Sleep(time.Second)
		ft.checkPosition(clientETHD, ft.cfg.FutureETH)
		time.Sleep(time.Second)
		ft.checkPosition(clientBTCDC, ft.cfg.FutureBTC)
		time.Sleep(time.Second)
		ft.checkPosition(clientETHDC, ft.cfg.FutureETH)
	}
}
