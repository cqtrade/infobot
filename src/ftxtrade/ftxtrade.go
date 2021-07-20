// https://github.com/ftexchange/ftx/blob/master/go/ftx/main.go

package ftxtrade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/types"
)

type FtxTrade struct {
	cfg        config.Config
	httpClient *http.Client
}

func New(cfg config.Config) *FtxTrade {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &FtxTrade{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (ft *FtxTrade) BuyEthBull(subAcc string) {
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	equity, _ := ft.CheckBalanceUSD("test1", client)
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
			fmt.Println("ERROR with Market order ", err)
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
				fmt.Println("FLOW SUCCESS")
			}
		} else {
			fmt.Println("FAILED market order, ", orderMarket)
		}
	}
}

func (ft *FtxTrade) CheckBalanceUSD(subAcc string, client *ftx.FtxClient) (float64, error) {
	sBalances, err := client.GetSubaccountBalances(subAcc)
	equity := 0.0
	if err != nil {
		return equity, err
	}
	fmt.Println(sBalances)
	if sBalances.Success {
		for _, balance := range sBalances.Result {
			if balance.Coin == "USD" {
				equity = balance.Free
			}
			fmt.Println("Coin\t", balance.Coin, "\tfree:\t", balance.Free, "\ttotal:\t", balance.Total)
		}
	}

	return equity, nil
}

func (ft *FtxTrade) GetLastPriceForMarket(market string, client *ftx.FtxClient) (float64, error) {
	marketPrice := 0.0
	candles, err := client.GetHistoricalPricesLatest(market, 60, 1)
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
		fmt.Println(totalETHBULL, ethBullPrice)
		totalETHBULLUSD = totalETHBULL * ethBullPrice
		totalEquityUSD := freeUSD + totalETHBULLUSD
		if totalEquityUSD == 0 {
			totalEquityUSD = 0.001
		}
		// cashPercentage := math.Round(freeUSD * 100 / totalEquityUSD)
		return " total: " + fmt.Sprintf("%.2f", totalEquityUSD) + " free USD: " + fmt.Sprintf("%.2f", freeUSD) +
			" cash: " + fmt.Sprintf("%.2f%%", freeUSD*100/totalEquityUSD)
	} else {
		return "No success getting balances for " + subAcc
	}

}

func (ft *FtxTrade) StartHealthPing() {
	for t := range time.Tick(time.Second) {
		if t.Second() == 33 {
			reqBody := types.NotificationBody{Content: t.Format("Jan 02 15:04") + ft.GetOverview("test1")}

			reqBodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				println("ERROR json.Marshal(reqBody)" + err.Error())
				return
			}

			req, err := http.NewRequest(http.MethodPost, ft.cfg.DiscordChHealth, bytes.NewBuffer(reqBodyBytes))

			if err != nil {
				println("ERROR preparing discord payload" + err.Error())
				return
			}

			req.Header.Add("Content-Type", "application/json")

			resp, err := ft.httpClient.Do(req)
			defer ft.httpClient.CloseIdleConnections()

			if err != nil {
				println("ERROR logger http " + err.Error())
				return
			}

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				println("Discord resp.StatusCode: " + fmt.Sprintf("%d", resp.StatusCode))
				return
			}
		}
	}
}
