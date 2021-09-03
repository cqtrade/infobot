// https://github.com/ftexchange/ftx/blob/master/go/ftx/main.go

package ftxtrade

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/state"
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
	tpPerc := 0.2

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

	tpUSD := RoundDown((equity * tpPerc), 4)
	ft.notif.Log("", "tpUSD", tpUSD)
	tpCoin := RoundDown((tpUSD / spotPrice), 4)
	ft.notif.Log("", "tpCoin", tpCoin)
	if tpCoin > balanceCoin.Free {
		tpCoin = RoundDown((balanceCoin.Free / 2), 4)
		ft.notif.Log("", "Less tpCoin", tpCoin)
	}

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

func (ft *FtxTrade) BuyCoinBull(subAcc string, market string) {
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	balanceCoinUSD, err := ft.CheckSpotBalance(client, subAcc, "USD")
	if err != nil {
		ft.notif.Log("ERROR", "USD", "BuyCoinBull CheckSpotBalance. Abort.", market, err.Error())
		return
	}
	equityUSD := balanceCoinUSD.Free
	positionSize := ft.cfg.PositionSize
	profitPercentage := ft.cfg.ProfitPercentage
	if equityUSD < positionSize {
		ft.notif.Log("INFO", "BuyCoinBull not enough cash. Abort.", market, positionSize, ">", equityUSD)
	}

	marketPrice, err := ft.appState.ReadLatestPriceForMarket(market)
	if err != nil {
		ft.notif.Log("ERROR", "BuyCoinBull ReadLatestPriceForMarket. Abort.", market, err.Error())
		return
	}

	if marketPrice <= 0.0 {
		ft.notif.Log("ERROR", "BuyCoinBull marketPrice price not greater than zero: ", marketPrice, market)
		return
	}

	// TODO use different rounding here?
	size := math.Round((positionSize/marketPrice)*10000) / 10000

	orderMarket, err := client.PlaceMarketOrder(market, "buy", "market", size)
	if err != nil {
		ft.notif.Log("ERROR", "BuyCoinBull Market BUY order. Abort.", market, err.Error())
		return
	}

	if orderMarket.Success == true {
		sizeTP := math.Round((size/2)*10000) / 10000 // sell 50%
		marketPrice, err := ft.appState.ReadLatestPriceForMarket(market)
		if err != nil {
			ft.notif.Log("ERROR", "BuyCoinBull ReadLatestPriceForMarket. Abort.", market, err.Error())
			return
		}

		priceTP := math.Round((marketPrice+marketPrice*profitPercentage)*10) / 10

		orderTP, err := client.PlaceOrder(market, "sell", priceTP, "limit", sizeTP, false, false, false)

		if err != nil {
			ft.notif.Log("ERROR", "BuyCoinBull TP order. Abort.", market, err.Error())
		} else if orderTP.Success {
			ft.notif.Log("INFO", "BuyCoinBull FLOW SUCCESS", market)
		}
	} else {
		ft.notif.Log("ERROR", "BuyCoinBull market order.", orderMarket, market)
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
