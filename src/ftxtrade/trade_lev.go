package ftxtrade

import (
	"fmt"
	"strings"
	"time"

	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
	"github.com/cqtrade/infobot/src/types"
)

// closing position should be easier, by position ID and delete method maybe?
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
		if err != nil {
			ft.notif.Log("ERROR", position.Future, "cancel all open orders", err.Error())
		} else if !co.Success {
			ft.notif.Log("ERROR", position.Future, "cancel all open orders", co)
		} else {
			ft.notif.Log("INFO", position.Future, "SUCCESS cancel all open orders")
		}
	}
}

func (ft *FtxTrade) swapEquity(
	side string,
	balanceUSD structs.SubaccountBalance,
	balanceCoin structs.SubaccountBalance,
	spotPrice float64,
	client *ftx.FtxClient,
	spotMarket string,
	coinUSD float64,
) {
	if side == "buy" && balanceUSD.Free > 10 {
		coinToBuySize := RoundDown(0.97*(balanceUSD.Free/spotPrice), 4)
		ft.notif.Log("", "coinToBuySize", coinToBuySize)
		orderMarketSpot, err := client.PlaceMarketOrder(spotMarket, "buy", "market", coinToBuySize)
		if err != nil {
			ft.notif.Log("ERROR", "TradeLevCrypto flip cash to crypto", err.Error())
			return
		}
		if orderMarketSpot.Success {
			ft.notif.Log("INFO", "TradeLevCrypto flip cash to crypto SUCCESS")
		} else {
			ft.notif.Log("ERROR", "TradeLevCrypto flip cash to crypto unsuccessful", orderMarketSpot.Result)
		}
	} else if side == "buy" && balanceUSD.Free <= 4 {
		ft.notif.Log("INFO", "TradeLevCrypto no flip cash to crypto USD", balanceUSD.Free)
	} else if side == "sell" && coinUSD > 10 {
		orderMarketSpot, err := client.PlaceMarketOrder(spotMarket, "sell", "market", balanceCoin.Free)
		if err != nil {
			ft.notif.Log("ERROR", "TradeLevCrypto flip cash to crypto", err.Error())
			return
		}
		if orderMarketSpot.Success {
			ft.notif.Log("INFO", "TradeLevCrypto flip crypto to cash SUCCESS")
		} else {
			ft.notif.Log("ERROR", "TradeLevCrypto flip crypto to cash unsuccessful", orderMarketSpot.Result)
		}
	} else if side == "sell" && coinUSD <= 4 {
		ft.notif.Log("INFO", "TradeLevCrypto no flip crypto to cash balanceCoin Free", balanceCoin.Free)
	}
}

func (ft *FtxTrade) TradeLevCrypto(
	msg types.JSONMessageBody,
	risk float64,
	side string,
	sideOpposite string,
	subAccType string,
) {
	ft.notif.Log("", msg)
	if risk <= 0.1 || msg.AtrTP <= 0 || msg.AtrSL <= 0 {
		ft.notif.Log("ERROR", "invalid risk|tp|sl. Abort.", msg)
		return
	}
	var shouldSwapEquity bool
	var market string
	var subAcc string
	var spotMarket string
	var coin string

	if subAccType == "dc" {
		shouldSwapEquity = true
	} else if subAccType == "d" {
		shouldSwapEquity = false
	} else {
		ft.notif.Log("ERROR", "subAccType not in dc|d. Abort.", subAccType)
		return
	}
	ticker := strings.ToUpper(msg.Ticker)
	if strings.HasPrefix(ticker, "BTC") || strings.HasPrefix(ticker, "XBT") {
		market = ft.cfg.FutureBTC
		if subAccType == "dc" {
			subAcc = ft.cfg.SubAccBTCDC
		} else if subAccType == "d" {
			subAcc = ft.cfg.SubAccBTCD
		}
		spotMarket = "BTC/USD"
		coin = "BTC"
	} else if strings.HasPrefix(ticker, "ETH") {
		market = ft.cfg.FutureETH
		if subAccType == "dc" {
			subAcc = ft.cfg.SubAccETHDC
		} else if subAccType == "d" {
			subAcc = ft.cfg.SubAccETHD
		}
		spotMarket = "ETH/USD"
		coin = "ETH"
	} else {
		ft.notif.Log("ERROR", "unknown ticker. Abort.", msg.Ticker)
		return
	}

	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)

	balanceUSD, err := ft.CheckSpotBalance(client, subAcc, "USD")
	if err != nil {
		ft.notif.Log("ERROR", "TradeLevCrypto get USD balance", err.Error())
		return
	}

	balanceCoin, err := ft.CheckSpotBalance(client, subAcc, coin)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLevCrypto get Coin balance", err.Error())
		return
	}

	spotPrice, err := ft.appState.ReadLatestPriceForMarket(spotMarket)
	coinUSD := balanceCoin.Free * spotPrice
	equity := balanceUSD.Free + coinUSD
	ft.notif.Log("", "equity", equity)
	position, _ := ft.CheckFuturePosition(client, market)

	if position.Size == 0 && shouldSwapEquity {
		ft.swapEquity(
			side,
			balanceUSD,
			balanceCoin,
			spotPrice,
			client,
			spotMarket,
			coinUSD,
		)
	}

	if position.Size == 0 && (side == "exitBuy" || side == "exitSell") {
		ft.notif.Log("INFO", "TradeLev No position nothing to exit. Abort.", msg.Signal)
		return
	}

	if position.Size != 0 && position.Side == "buy" && (side == "exitBuy" || side == "sell") { // exit long
		ft.closePosition(client, market, position)
		return
	}

	if position.Size != 0 && position.Side == "sell" && (side == "exitSell" || side == "buy") { // exit short
		ft.closePosition(client, market, position)
		return
	}

	if position.Size != 0 && (side == "buy" || side == "sell") {
		ft.notif.Log("INFO", "TradeLev already in position, no entry. Abort.", side, msg.Signal)
		return
	}

	time.Sleep(time.Second)
	spotPrice, err = ft.appState.ReadLatestPriceForMarket(spotMarket)

	atrRiskPerc := msg.AtrSL * 100 / spotPrice
	riskRatio := risk / atrRiskPerc
	positionSize := RoundDown((equity*riskRatio)/spotPrice, 4)
	ft.notif.Log("INFO", "positionSize", positionSize)
	if positionSize < 0.0001 {
		ft.notif.Log("INFO", "TradeLev positionSize less than 0.0001. Abort.", fmt.Sprintf("%.6f", positionSize))
		return
	}
	ft.notif.Log("", market, side)

	orderMarketFuture, err := client.PlaceMarketOrder(market, side, "market", positionSize)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLev orderMarketFuture. Abort.", err.Error())
		return
	}
	if !orderMarketFuture.Success {
		ft.notif.Log("ERROR", "TradeLev orderMarketFuture no success. Abort.", orderMarketFuture.HTTPCode, orderMarketFuture.ErrorMessage)
		return
	}

	ft.notif.Log("INFO", "TradeLev orderMarketFuture SUCCESS.")
	tpSize := RoundUp(positionSize/3, 4)
	slSize := positionSize

	spotPrice, err = ft.appState.ReadLatestPriceForMarket(spotMarket)
	var slPrice float64
	var tpPrice float64
	if side == "buy" {
		slPrice = Round(spotPrice-msg.AtrSL, 1)
		tpPrice = Round(spotPrice+msg.AtrTP, 1)
	} else if side == "sell" {
		slPrice = Round(spotPrice+msg.AtrSL, 1)
		tpPrice = Round(spotPrice-msg.AtrTP, 1)
	}
	ft.notif.Log("", "price", fmt.Sprintf("%.2f", spotPrice))
	ft.notif.Log("", sideOpposite, "slSize", fmt.Sprintf("%.4f", slSize), "slPrice", fmt.Sprintf("%.2f", slPrice))
	ft.notif.Log("", sideOpposite, "tpSize", fmt.Sprintf("%.4f", tpSize), "tpPrice", fmt.Sprintf("%.2f", tpPrice))
	slOrder, err := client.PlaceTriggerOrder(market, sideOpposite, slSize, "stop", true, true, slPrice, 0.0, 0.0)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLev SL order. TODO close position. Abort.", err.Error())
		return
	}
	if !slOrder.Success {
		ft.notif.Log("ERROR", "TradeLev SL order UNSUCCESSFUL. TODO close position. Abort.", slOrder.HTTPCode, slOrder.ErrorMessage)
		return
	}
	ft.notif.Log("INFO", "TradeLev SL SUCCESS.")
	tpOrder, err := client.PlaceTriggerOrder(market, sideOpposite, tpSize, "takeProfit", true, true, tpPrice, 0.0, 0.0)
	if err != nil {
		ft.notif.Log("ERROR", "TradeLev TP order. Check manually.", err.Error())
		return
	}
	if !tpOrder.Success {
		ft.notif.Log("ERROR", "TradeLev TP order UNSUCCESSFUL. Check manually.", tpOrder.HTTPCode, tpOrder.ErrorMessage)
		return
	}
	ft.notif.Log("INFO", "TradeLev TP SUCCESS.")
}
