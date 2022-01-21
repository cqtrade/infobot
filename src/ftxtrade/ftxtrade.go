// https://github.com/ftexchange/ftx/blob/master/go/ftx/main.go

package ftxtrade

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/state"
	"github.com/markcheno/go-talib"
)

type FtxTrade struct {
	cfg        config.Config
	notif      notification.Notification
	appState   state.State
	httpClient *http.Client
}

func New(cfg config.Config, notif notification.Notification, appState state.State) *FtxTrade {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &FtxTrade{
		cfg:        cfg,
		notif:      notif,
		appState:   appState,
		httpClient: httpClient,
	}
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func RoundDown(val float64, precision int) float64 {
	return math.Floor(val*(math.Pow10(precision))) / math.Pow10(precision)
}

func RoundUp(val float64, precision int) float64 {
	return math.Ceil(val*(math.Pow10(precision))) / math.Pow10(precision)
}

func Round(val float64, precision int) float64 {
	return math.Round(val*(math.Pow10(precision))) / math.Pow10(precision)
}

func (ft *FtxTrade) TpCoinBull(subAcc string, market string, coin string) {
	tpPerc := 0.05

	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)

	balanceCoin, err := ft.CheckSpotBalance(client, subAcc, coin)
	if err != nil {
		ft.notif.Log("ERROR", coin, "TpCoinBull CheckSpotBalance. Abort.", err.Error())
		return
	}
	if balanceCoin.Free < 0.0001 {
		ft.notif.Log("", "TpCoinBull too little coin to take profit. Abort", balanceCoin.Free)
		return
	}

	spotPrice, err := ft.appState.ReadLatestPriceForMarket(market)
	if err != nil {
		ft.notif.Log("ERROR", market, "TpCoinBull ReadLatestPriceForMarket. Abort.", err.Error())
		return
	}
	coinUSD := balanceCoin.Free * spotPrice
	balanceUSD, err := ft.CheckSpotBalance(client, subAcc, "USD")
	if err != nil {
		ft.notif.Log("ERROR", "TpCoinBull get USD balance", err.Error())
		return
	}
	equity := balanceUSD.Free + coinUSD
	ft.notif.Log("", "free coin", balanceCoin.Free)
	ft.notif.Log("", "equity", equity)

	tpUSD := RoundDown((coinUSD * tpPerc), 4)
	ft.notif.Log("", "tpUSD", tpUSD)
	tpCoin := RoundDown((balanceCoin.Free * tpPerc), 4)

	size := tpCoin

	if size*spotPrice >= 10 {
		order, err := client.PlaceMarketOrder(market, "sell", "market", size)
		if err != nil {
			ft.notif.Log("ERROR", "TpCoinBull Market BUY order. Abort.", market, err.Error())
			return
		}
		if !order.Success {
			ft.notif.Log("ERROR", "TpCoinBull  UNSUCCESSFUL", market, order.HTTPCode, order.ErrorMessage)
			return
		}
		ft.notif.Log("INFO", "TpCoinBull FLOW SUCCESS", market)
	} else {
		ft.notif.Log("INFO", "TpCoinBull too small capital, Fraction value. Abort", market, size*spotPrice)
	}
}

// https://yourbasic.org/golang/convert-string-to-float/
func (ft *FtxTrade) GetOverview(subAcc string) string {
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	sBalances, err := client.GetSubaccountBalances(subAcc)
	msg := subAcc
	if err != nil {
		return msg + " error receiveing balances: " + err.Error()
	}
	freeUSD := 0.0
	totalETHBULL := 0.0
	totalETHBULLUSD := 0.0
	if sBalances.Success {
		for _, balance := range sBalances.Result {
			if balance.Coin == "USD" {
				freeUSD = balance.Free
			} else if balance.Coin == "ETHBULL" {
				totalETHBULL = balance.Total
			}
		}

		ethBullPrice, err := ft.appState.ReadLatestPriceForMarket("ETHBULL/USD")
		if err != nil {
			ft.notif.Log("ERROR", "ArbStart ReadLatestPriceForMarket. Abort.", err.Error())
			return ""
		}

		totalETHBULLUSD = totalETHBULL * ethBullPrice
		totalEquityUSD := freeUSD + totalETHBULLUSD
		if totalEquityUSD == 0 {
			totalEquityUSD = 0.00001
		}
		return subAcc + " cash: " + fmt.Sprintf("%.2f%%", freeUSD*100/totalEquityUSD)
	} else {
		return "No success getting balances for " + subAcc
	}
}

type Item struct {
	coin  string
	alloc float64
}

func (ft *FtxTrade) BuyCoin(
	subAcc string,
	coin string,
	buyQty float64,
	sellQty float64,
	tp float64,
) {
	market := coin + "/USD"
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)

	// positionSize := positionUSD
	// marketPrice, err := ft.appState.ReadLatestPriceForMarket(market)

	// if err != nil {
	// 	ft.notif.Log("ERROR", "BuyCoin ReadLatestPriceForMarket. Abort.", market, err.Error())
	// 	return
	// }

	// if marketPrice <= 0.0 {
	// 	ft.notif.Log("ERROR", "BuyCoin marketPrice price not greater than zero: ", marketPrice, market)
	// 	return
	// }

	// TODO use different rounding here?
	// size := math.Round((positionSize/marketPrice)*10000) / 10000
	size := buyQty

	orderMarket, err := client.PlaceMarketOrder(market, "buy", "market", size)
	if err != nil {
		ft.notif.Log("ERROR", "BuyCoin Market BUY order. Abort.", market, err.Error())
		return
	}

	if !orderMarket.Success {
		ft.notif.Log("ERROR", "BuyCoin market order.", orderMarket, market)
	}

	sizeTP := sellQty

	priceTP := tp

	orderTP, err := client.PlaceOrder(market, "sell", priceTP, "limit", sizeTP, false, false, false)

	if err != nil {
		ft.notif.Log("ERROR", "BuyCoin TP order. Abort.", market, err.Error())
	} else if orderTP.Success {
		ft.notif.Log("INFO", "BuyCoin FLOW SUCCESS", market)
	}
}

func (ft *FtxTrade) Portfolio(subAcc string) {
	var portfolio = []Item{
		// {
		// 	coin:  "BULL",
		// 	alloc: 1,
		// },
		{
			coin:  "ATOMBULL",
			alloc: 0.5,
		},
		{
			coin:  "MATICBULL",
			alloc: 0.5,
		},
	}

	fmt.Println("portfolio", portfolio, len(portfolio))
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	sBalances, err := client.GetSubaccountBalances(subAcc)

	if err != nil {
		ft.notif.Log("ERROR", subAcc, " Abort. error receiveing balances", err.Error())
		return
	}

	freeUSD := 0.0
	if !sBalances.Success {
		ft.notif.Log("ERROR", "Abort. No success getting balances for ", subAcc)
		return

	}

	for _, balance := range sBalances.Result {
		if balance.Coin == "USD" {
			freeUSD = balance.Free
		}
	}

	usd := 0.99 * freeUSD / 2

	ft.notif.Log("", "USD", usd)
	for _, item := range portfolio {
		fmt.Println(item)
		// GET ATR

		atr, _, err := ft.GetAtr(item.coin)
		if err != nil {
			ft.notif.Log("ERROR", "Abort. No success getting atr ", err.Error())
			return
		}
		market := item.coin + "/USD"
		close, err := ft.appState.ReadLatestPriceForMarket(market)

		if err != nil {
			ft.notif.Log("ERROR", "BuyCoin portfolio. Abort.", market, err.Error())
			return
		}

		if close <= 0.0 {
			ft.notif.Log("ERROR", "BuyCoin portfolio price not greater than zero: ", close, market)
			return
		}

		tp := close + 3.0*atr
		profitPercentage := (tp / close) - 1.0
		fmt.Println("%", profitPercentage)
		// p = 0.05
		t := 0.02

		if profitPercentage < t {
			ft.notif.Log("ERROR", "Abort. Coin Buy p less than ", t, profitPercentage, item.coin)
			continue
		}
		equity := usd * item.alloc
		// equity = 1000
		// fmt.Println(usd)
		// fmt.Println(equity)
		s := equity / (equity + equity*((profitPercentage*100)/100.0))
		buyQty := Round(equity/close, 8)
		sellQty := Round(buyQty*s*1.02, 8)
		remainingQty := buyQty - sellQty
		fmt.Println("Buy Qty", buyQty, "@", close, "sell Qty", sellQty,
			"TP", tp)
		// fmt.Println("Buy Qty", buyQtyUsd, "@", close, "sell Qty", sellQtyUsd)
		fmt.Println("remainingQty", remainingQty,
			"remaining money", remainingQty*(close+close*profitPercentage))
		ft.BuyCoin(subAcc, item.coin, buyQty, sellQty, tp)
		time.Sleep(time.Second)

	}
	/*
		ATOMBULL/USD
		MATICBULL/USD
		ETHBULL/USD
		BULL/USD
		ADABULL/USD
		BNBBULL/USD
		SOL/USD
		BTC/USD
	*/

}

func (ft *FtxTrade) GetAtr(coin string) (atr float64, close float64, err error) {
	market := coin + "/USD"
	fmt.Println(market)
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, "")
	defer client.Client.CloseIdleConnections()

	var m15 int64
	m15 = 15 * 60
	ohclv, err := client.GetHistoricalPriceLatest(market, m15, 200)

	if err != nil {
		return atr, close, err
	}

	if !ohclv.Success {
		fmt.Println(fmt.Sprintf("HERE1 %d ", ohclv.HTTPCode) + ohclv.ErrorMessage)
		return atr, close, err
	}

	candles := make([]structs.HistoricalPrice, len(ohclv.Result))
	for i, candle := range ohclv.Result {
		candles[i] = candle
		candles[i].StartTime = candles[i].StartTime.Add(time.Minute * 15)
	}

	// if len(candles) > 0 { // rm ongoing candle
	// 	candles = candles[:len(candles)-1]
	// }

	closes := make([]float64, len(candles))
	highs := make([]float64, len(candles))
	lows := make([]float64, len(candles))
	for i, candle := range candles {
		closes[i] = candle.Close
		highs[i] = candle.High
		lows[i] = candle.Low
	}

	atrs := talib.Atr(highs, lows, closes, 14)

	if len(atrs) > 0 {
		atr = atrs[len(atrs)-1]
		close = closes[len(closes)-1]
		return atr, close, err
	}

	return atr, close, errors.New("no atr")
}
