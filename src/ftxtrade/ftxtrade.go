package ftxtrade

import (
	"fmt"
	"log"
	"math"

	"github.com/cqtrade/infobot/src/ftx"
)

// https://github.com/ftexchange/ftx/blob/master/go/ftx/main.go
func StartStuff() {
	key := "key"
	secret := "secret"
	subAcc := "s"
	client := ftx.New(key, secret, subAcc)
	// positions, err := client.GetPositions(true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(positions)
	sBalances, err := client.GetSubaccountBalances(subAcc)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(sBalances)
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
	if equity > 25 {
		candles, _ := client.GetHistoricalPricesLatest(market, 60, 1)
		if len(candles.Result) > 0 {
			candle := candles.Result[0]
			fmt.Println(candle.Close)
			marketPrice := candle.Close
			size := math.Round((10/candle.Close)*10000) / 10000
			fmt.Println("Buy market")
			fmt.Println(market)
			fmt.Println(size)
			orderMarket, _ := client.PlaceMarketOrder(market, "buy", "market", size)
			fmt.Println(orderMarket)
			if orderMarket.Success == true {
				sizeTP := math.Round((size/2)*10000) / 10000
				priceTP := math.Round((marketPrice+marketPrice*0.05)*10) / 10
				fmt.Println("TP")
				fmt.Println(sizeTP)
				fmt.Println(priceTP)
				orderTP, _ := client.PlaceOrder(market, "sell", priceTP, "limit", sizeTP, false, false, false)
				fmt.Println(orderTP)
			}
		}
	}
}
