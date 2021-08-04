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

func (ft *FtxTrade) TpEthBull(subAcc string) {
	fractionPerc := 0.2
	// get ETHBULL balance
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	balanceETHBULL, _ := ft.CheckFreeBalanceETHBULL(subAcc, client)
	market := "ETHBULL/USD"
	marketPrice, _ := ft.GetLastPriceForMarket(market, client)
	// calculate value
	valueUSD := balanceETHBULL * marketPrice
	size := math.Round((balanceETHBULL*fractionPerc)*10000) / 10000
	// calculate 10% or 20% of value
	fractionUSD := valueUSD * fractionPerc
	if fractionUSD > 4 {
		_, err := client.PlaceMarketOrder(market, "sell", "market", size)
		if err != nil {
			fmt.Println("ERROR with Market BUY order ", err)
		} else {
			fmt.Println("TP SELL FLOW SUCCESS ")
		}
	} else {
		fmt.Println("Fraction value ", fractionUSD)
	}
	// if it is > 10 sell market or add trailing stop
}

func (ft *FtxTrade) BuyEthBull(subAcc string) {
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	equity, _ := ft.CheckBalanceUSD(subAcc, client)
	market := "ETHBULL/USD"
	positionSize := ft.cfg.PositionSize
	profitPercentage := ft.cfg.ProfitPercentage
	if equity > positionSize {
		marketPrice, _ := ft.GetLastPriceForMarket(market, client)
		if marketPrice <= 0.0 {
			fmt.Println(market, " marketPrice price not greater than zero: ", marketPrice)
			return
		}

		size := math.Round((positionSize/marketPrice)*10000) / 10000
		fmt.Println("Buy market\t", market, "\t", size)
		orderMarket, err := client.PlaceMarketOrder(market, "buy", "market", size)
		if err != nil {
			fmt.Println("ERROR with Market BUY order ", err)
		}
		fmt.Println(orderMarket)
		if orderMarket.Success == true {
			sizeTP := math.Round((size/2)*10000) / 10000 // sell 50%
			marketPrice, _ := ft.GetLastPriceForMarket(market, client)
			priceTP := math.Round((marketPrice+marketPrice*profitPercentage)*10) / 10
			fmt.Println("TP size: ", sizeTP, " price: ", priceTP)
			orderTP, err := client.PlaceOrder(market, "sell", priceTP, "limit", sizeTP, false, false, false)
			fmt.Println(orderTP)
			if err != nil {
				fmt.Println("ERROR with TP order ", err)
			} else if orderTP.Success {
				fmt.Println("BUY FLOW SUCCESS")
			}
		} else {
			fmt.Println("FAILED market order, ", orderMarket)
		}
	}
}

func (ft *FtxTrade) GetLastPriceForMarket(market string, client *ftx.FtxClient) (float64, error) {
	marketPrice := 0.0
	candles, err := client.GetHistoricalPriceLatest(market, 15, 1)
	if err != nil {
		return marketPrice, err
	}

	if len(candles.Result) > 0 {
		candle := candles.Result[0]
		marketPrice = candle.Close
	}

	return marketPrice, nil
}

// https://yourbasic.org/golang/convert-string-to-float/
func (ft *FtxTrade) GetOverview(subAcc string) string {
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	sBalances, err := client.GetSubaccountBalances(subAcc)
	msg := subAcc
	if err != nil {
		return msg + " error receiveing balnces: " + err.Error()
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
		ethBullPrice, _ := ft.GetLastPriceForMarket("ETHBULL/USD", client)
		totalETHBULLUSD = totalETHBULL * ethBullPrice
		totalEquityUSD := freeUSD + totalETHBULLUSD
		if totalEquityUSD == 0 {
			totalEquityUSD = 0.001
		}
		return " total: " + fmt.Sprintf("%.2f", totalEquityUSD) + " free USD: " + fmt.Sprintf("%.2f", freeUSD) +
			" cash: " + fmt.Sprintf("%.2f%%", freeUSD*100/totalEquityUSD)
	} else {
		return "No success getting balances for " + subAcc
	}

}
