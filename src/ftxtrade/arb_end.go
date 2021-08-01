package ftxtrade

import (
	"fmt"

	"github.com/cqtrade/infobot/src/ftx"
)

func (ft *FtxTrade) ArbEnd(subAcc string, spot string) {
	var spotUSD string
	var future string
	if spot == "BTC" {
		spotUSD = "BTC/USD"
		future = ft.cfg.FutureBTC
	} else if spot == "ETH" {
		spotUSD = "ETH/USD"
		future = ft.cfg.FutureETH
	}
	fmt.Println("spotUSD: ", spotUSD)
	fmt.Println("future: ", future)

	// check if there's future short to close
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)

	sBalance, _ := ft.CheckSpotBalance(client, subAcc, spot)
	position, _ := ft.CheckFuturePosition(client, future)

	// if !positions.Success {
	// 	fmt.Println("FAILURE arb1 positions")
	// 	return
	// }
	fmt.Println(sBalance.Total, position.Size)
	if sBalance.Total > 0 {
		orderMarketSpot, err := client.PlaceMarketOrder(spotUSD, "sell", "market", sBalance.Total)
		if err != nil {
			fmt.Println("Spot sell err", err)
		} else if orderMarketSpot.Success {
			fmt.Println("Spot sell success")
		} else if !orderMarketSpot.Success {
			fmt.Println("Spot sell failure", orderMarketSpot.Result)
		}
	} else {
		fmt.Println("No spot to sell")
	}

	if position.Size > 0 {
		orderMarketFuture, err := client.PlaceMarketOrder(future, "buy", "market", position.Size)
		if err != nil {
			fmt.Println("Future buy err", err)
		} else if orderMarketFuture.Success {
			fmt.Println("Future close success")
		}
	} else {
		fmt.Println("No future to close")
	}

}
