// https://github.com/ftexchange/ftx/blob/master/go/ftx/main.go

package ftxtrade

import (
	"fmt"
	"math"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftx"
)

type FtxTrade struct {
	cfg config.Config
}

func New(cfg config.Config) *FtxTrade {
	return &FtxTrade{
		cfg: cfg,
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
