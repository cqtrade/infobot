package ftxtrade

import (
	"fmt"

	"github.com/cqtrade/infobot/src/ftx"
)

func (ft *FtxTrade) Arb(subAcc string, name string) {
	var spot string
	var future string
	if name == "BTC" {
		spot = "BTC/USD"
		future = ft.cfg.FutureBTC
	} else if name == "ETH" {
		spot = "ETH/USD"
		future = ft.cfg.FutureETH
	}
	fmt.Println("spot: ", spot)
	fmt.Println("future: ", future)
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)
	equity, _ := ft.CheckBalanceUSD(subAcc, client)
	size := 0.0011 // math.Round(equity * 0.95)
	if equity < 20 {
		fmt.Println("Size too small for arb: ", fmt.Sprintf("%.0f", size))
		return
	}

	orderMarketSpot, err := client.PlaceMarketOrder(spot, "buy", "market", size)
	if err == nil && orderMarketSpot.Success {
		sizeSpot, err := ft.CheckFreeBalanceBTC(subAcc, client)
		fmt.Println("Size spot: ", fmt.Sprintf("%.0f", sizeSpot))
		orderMarketFuture, err := client.PlaceMarketOrder(future, "sell", "market", size)
		if err == nil && orderMarketFuture.Success {
			fmt.Println("Arb success: ", future)
			return
		}
		fmt.Println("ERROR future: ", err, " res ", orderMarketFuture)
		return
	}
	fmt.Println("ERROR spot: ", err, " res ", orderMarketSpot)
}
