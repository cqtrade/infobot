package ftxtrade

import (
	"fmt"

	"github.com/cqtrade/infobot/src/ftx"
)

func (ft *FtxTrade) ArbStart(subAcc string, name string) {
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

	sBalance, err := ft.CheckSpotBalance(client, subAcc, "USD")
	if err != nil {
		fmt.Println("Error getting balance", err)
		return
	}
	marketPrice, err := ft.GetLastPriceForMarket(spot, client)
	if err != nil {
		fmt.Println("Error getting market price", err)
		return
	}
	// 0.00110000
	size := RoundDown(0.9*(sBalance.Free/marketPrice), 4)
	fmt.Println(size)
	fmt.Println("price ", marketPrice, " sBalance: ", sBalance.Free, " size ", fmt.Sprintf("%.4f", size))
	if size < 0.0011 {
		fmt.Println("Size too small for arb: ", fmt.Sprintf("%.4f", size))
		return
	}

	orderMarketSpot, err := client.PlaceMarketOrder(spot, "buy", "market", size)
	if err == nil && orderMarketSpot.Success {
		orderMarketFuture, err := client.PlaceMarketOrder(future, "sell", "market", size)
		if err == nil && orderMarketFuture.Success {
			fmt.Println("Arb success")
			sBalance, _ := ft.CheckSpotBalance(client, subAcc, name)
			position, _ := ft.CheckFuturePosition(client, future)
			fmt.Println(sBalance.Coin, "\t", sBalance.Total, "\t", position.Future, "\t", position.Side, "\t", position.Size)
			return
		}
		fmt.Println("ERROR future: ", err, " res ", orderMarketFuture)
		return
	}
	fmt.Println("ERROR spot: ", err, " res ", orderMarketSpot)
}
