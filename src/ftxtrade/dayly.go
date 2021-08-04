package ftxtrade

import (
	"fmt"
	"strings"
	"time"

	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
	"github.com/cqtrade/infobot/src/types"
)

func (ft *FtxTrade) closePosition(client *ftx.FtxClient, market string, position structs.Position) {
	var side string
	if position.Side == "buy" {
		side = "sell"
	} else if position.Side == "sell" {
		side = "buy"
	}
	orderMarketFuture, err := client.PlaceMarketOrder(market, side, "market", position.Size)
	if err != nil {
		ft.notif.Log("ERROR", position.Future, "closing", position.Side, err.Error())
	} else if orderMarketFuture.Success {
		ft.notif.Log("INFO", position.Future, "SUCCESS closing", position.Side)
		co, err := client.CancelAllOrders()
		if err != nil || !co.Success {
			ft.notif.Log("ERROR", position.Future, "cancel all open orders", err.Error(), co)
		} else {
			ft.notif.Log("INFO", position.Future, "SUCCESS cancel all open orders")
		}
	}
}

func (ft *FtxTrade) handlePosition(client *ftx.FtxClient, market string) {
	position, _ := ft.CheckFuturePosition(client, market)
	ft.notif.Log("", position)

	openOrders, err := client.GetOpenOrders(market)
	if err != nil || !openOrders.Success {
		ft.notif.Log("ERROR", "getting open orders", err.Error(), openOrders)
	}
	openTriggerOrders, err := client.GetOpenTriggerOrders(market, "stop")
	if err != nil || !openTriggerOrders.Success {
		ft.notif.Log("ERROR", "getting open trigger orders", err.Error(), openTriggerOrders)
	}

	fmt.Println("#####")
	fmt.Println("Open orders length", len(openOrders.Result))
	fmt.Println("Open trigger orders length", len(openTriggerOrders.Result))
	fmt.Println("#####")
}

func (ft *FtxTrade) CheckPositionBTC(msg types.JSONMessageBody) {
	subAcc := ft.cfg.SubAccBTCD
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	ft.handlePosition(client, ft.cfg.FutureBTC)
}

func (ft *FtxTrade) CheckPositionETH(msg types.JSONMessageBody) {
	subAcc := ft.cfg.SubAccETHD
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	ft.handlePosition(client, ft.cfg.FutureETH)
}

func (ft *FtxTrade) TradeLev(msg types.JSONMessageBody) {
	ft.notif.Log("", msg)
	var market string
	var subAcc string
	ticker := strings.ToUpper(msg.Ticker)
	if strings.HasPrefix(ticker, "BTC") || strings.HasPrefix(ticker, "XBT") {
		market = ft.cfg.FutureBTC
		subAcc = ft.cfg.SubAccBTCD
	} else if strings.HasPrefix(ticker, "ETH") {
		market = ft.cfg.FutureETH
		subAcc = ft.cfg.SubAccETHD
	} else {
		ft.notif.Log("ERROR", "unknown ticker. Abort.", msg.Ticker)
		return
	}

	ft.notif.Log("", "market", market, "subaccount", subAcc)

	if msg.Risk == 0 || msg.AtrTP == 0 || msg.AtrSL == 0 {
		ft.notif.Log("ERROR", "incorrect data - missing risk or tp or sl. Abort.", msg)
		return
	}

	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)

	position, _ := ft.CheckFuturePosition(client, market)

	if position.Size == 0 && (msg.Signal == 2 || msg.Signal == -2) {
		ft.notif.Log("INFO", "TradeLev No position nothing to exit. Abort.", msg.Signal)
		return
	}

	if position.Size != 0 && position.Side == "buy" && (msg.Signal == 2 || msg.Signal == -1) { // exit long
		ft.closePosition(client, market, position)
		return
	}

	if position.Size != 0 && position.Side == "sell" && (msg.Signal == -2 || msg.Signal == 1) { // exit short
		ft.closePosition(client, market, position)
		return
	}

	if position.Size != 0 && (msg.Signal == 1 || msg.Signal == -1) {
		ft.notif.Log("ERROR", "TradeLev already in position, no entry. Abort.", msg.Signal)
		return
	}

	// create new position logic, for now equity USD based
	sBalanceUSD, err := ft.CheckSpotBalance(client, subAcc, "USD")

	if err != nil {
		fmt.Println("Error getting balance", err)
		ft.notif.Log("ERROR", "TradeLev CheckSpotBalance. Abort.", err.Error(), sBalanceUSD)
		return
	}

	if sBalanceUSD.Free < 10 {
		ft.notif.Log("ERROR", "TradeLev Free USD less than 10. Abort.", sBalanceUSD)
		return
	}

	price, err := ft.appState.ReadLatestPriceForMarket(market)

	if err != nil {
		ft.notif.Log("ERROR", "TradeLev ReadLatestPriceForMarket. Abort.", err.Error())
		return
	}

	ft.notif.Log("", market, price)

	equity := RoundDown(sBalanceUSD.Free/price, 4)
	ft.notif.Log("", "equity", fmt.Sprintf("%.4f", equity))
	if equity < 0.0002 {
		ft.notif.Log("ERROR", "TradeLev Equity less than 0.0002. Abort.", equity)
		return
	}

	atrRiskPerc := msg.AtrSL * 100 / price
	riskRatio := msg.Risk / atrRiskPerc
	positionSize := RoundDown(equity*riskRatio, 4)
	ft.notif.Log("", "Position size", market, fmt.Sprintf("%.4f", positionSize))
	if positionSize < 0.0001 {
		ft.notif.Log("ERROR", "TradeLev positionSize less than 0.0001. Abort.", fmt.Sprintf("%.6f", positionSize))
		return
	}

	var side string
	var sideOpposite string
	if msg.Signal == 1 {
		side = "buy"
		sideOpposite = "sell"
	} else if msg.Signal == -1 {
		side = "sell"
		sideOpposite = "buy"
	}

	ft.notif.Log("", market, side)

	time.Sleep(time.Second)

	orderMarketFuture, err := client.PlaceMarketOrder(market, side, "market", positionSize)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLev orderMarketFuture. Abort.", err.Error(), orderMarketFuture)
		return
	}
	if !orderMarketFuture.Success {
		ft.notif.Log("ERROR", "TradeLev orderMarketFuture no success. Abort.", orderMarketFuture)
		return
	}

	ft.notif.Log("INFO", "TradeLev orderMarketFuture SUCCESS.", orderMarketFuture.Result)

	time.Sleep(time.Second)

	tpSize := RoundUp(positionSize/3, 4)
	slSize := positionSize

	price, err = ft.appState.ReadLatestPriceForMarket(market)
	var slPrice float64
	var tpPrice float64
	if msg.Signal == 1 {
		slPrice = Round(price-msg.AtrSL, 1)
		tpPrice = Round(price+msg.AtrTP, 1)
	} else if msg.Signal == -1 {
		slPrice = Round(price+msg.AtrSL, 1)
		tpPrice = Round(price-msg.AtrTP, 1)
	}

	ft.notif.Log("", "price", fmt.Sprintf("%.2f", price))
	ft.notif.Log("", sideOpposite, "slSize", fmt.Sprintf("%.4f", slSize), "slPrice", fmt.Sprintf("%.2f", slPrice))
	ft.notif.Log("", sideOpposite, "tpSize", fmt.Sprintf("%.4f", tpSize), "tpPrice", fmt.Sprintf("%.2f", tpPrice))

	slOrder, err := client.PlaceTriggerOrder(market, sideOpposite, slSize, "stop", true, true, slPrice, slPrice, 0.0)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLev SL order. TODO close position. Abort.", err.Error(), slOrder.Result)
		return
	}
	if !slOrder.Success {
		ft.notif.Log("ERROR", "TradeLev SL order UNSUCCESSFUL. TODO close position. Abort.", slOrder.Result)
		return
	}
	ft.notif.Log("INFO", "TradeLev SL SUCCESS.", slOrder.Result)

	ft.notif.Log("", sideOpposite, tpSize, tpPrice)
	tpOrder, err := client.PlaceOrder(market, sideOpposite, tpPrice, "limit", tpSize, true, false, true)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLev TP order. Check manually.", err.Error(), tpOrder.Result)
		return
	}
	if !tpOrder.Success {
		ft.notif.Log("ERROR", "TradeLev TP order UNSUCCESSFUL. Check manually.", tpOrder.Result)
		return
	}
	ft.notif.Log("INFO", "TradeLev TP SUCCESS.", tpOrder.Result)
}
