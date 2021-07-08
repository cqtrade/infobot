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
	sBalances, err := client.GetSubaccountBalances(subAcc)
	if err != nil {
		fmt.Println(err)
	}
	equity := 0.0
	if sBalances.Success {
		for _, balance := range sBalances.Result {
			if balance.Coin == "USDT" {
				equity = balance.Free
			}
		}
	}
	fmt.Println("Available in test1 USDT", equity)
	market := "ETHBULL/USDT"
	positionSize := ft.cfg.PositionSize
	profitPercentage := ft.cfg.ProfitPercentage
	if equity > positionSize {
		candles, _ := client.GetHistoricalPricesLatest(market, 60, 1)
		if len(candles.Result) > 0 {
			candle := candles.Result[0]
			fmt.Println(candle.Close)
			marketPrice := candle.Close
			size := math.Round((positionSize/candle.Close)*10000) / 10000
			fmt.Println("Buy market")
			fmt.Println(market)
			fmt.Println(size)
			orderMarket, _ := client.PlaceMarketOrder(market, "buy", "market", size)
			fmt.Println(orderMarket)
			if orderMarket.Success == true {
				sizeTP := math.Round((size/2)*10000) / 10000 // sell 50%
				priceTP := math.Round((marketPrice+marketPrice*profitPercentage)*10) / 10
				fmt.Println("TP")
				fmt.Println(sizeTP)
				fmt.Println(priceTP)
				orderTP, _ := client.PlaceOrder(market, "sell", priceTP, "limit", sizeTP, false, false, false)
				fmt.Println(orderTP)
			}
		}
	}
}
