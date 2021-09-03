package ftx

import (
	"log"
	"strconv"

	"github.com/cqtrade/infobot/src/ftx/structs"
)

type HistoricalPrices structs.HistoricalPrices
type Trades structs.Trades

func (client *FtxClient) GetHistoricalPrices(market string, resolutionSeconds int64,
	limit int64, startTime int64, endTime int64) (HistoricalPrices, error) {
	var historicalPrices HistoricalPrices
	resp, err := client._get(
		"markets/"+market+
			"/candles?resolution="+strconv.FormatInt(resolutionSeconds, 10)+
			"&limit="+strconv.FormatInt(limit, 10)+
			"&start_time="+strconv.FormatInt(startTime, 10)+
			"&end_time="+strconv.FormatInt(endTime, 10),
		[]byte(""))
	if err != nil {
		log.Printf("Error GetHistoricalPrices", err)
		return historicalPrices, err
	}
	err = _processResponse(resp, &historicalPrices)
	historicalPrices.HTTPCode = resp.StatusCode
	return historicalPrices, err
}

/*
Name	Type	Value	Description
market_name	string	BTC-0628	name of the market
resolution	number	300	window length in seconds. options: 15, 60, 300, 900, 3600, 14400, 86400, or any multiple of 86400 up to 30*86400
start_time	number	1559881511	optional
end_time	number	1559881711	optional
*/
func (client *FtxClient) GetHistoricalPriceLatestNoRetry(
	market string,
	resolution int64,
	limit int64,
) (HistoricalPrices, error) {
	var historicalPrices HistoricalPrices
	resp, err := client._get(
		"markets/"+market+
			"/candles?resolution="+strconv.FormatInt(resolution, 10)+
			"&limit="+strconv.FormatInt(limit, 10),
		[]byte(""))
	if err != nil {
		log.Printf("Error GetHistoricalPrices", err)
		return historicalPrices, err
	}
	err = _processResponse(resp, &historicalPrices)
	historicalPrices.HTTPCode = resp.StatusCode
	return historicalPrices, err
}

func (client *FtxClient) GetHistoricalPriceLatest(
	market string,
	resolution int64,
	limit int64,
) (HistoricalPrices, error) {
	var historicalPrices HistoricalPrices
	resp, err := client._getRetry(
		2,
		"markets/"+market+
			"/candles?resolution="+strconv.FormatInt(resolution, 10)+
			"&limit="+strconv.FormatInt(limit, 10),
		[]byte(""),
	)
	if err != nil {
		log.Printf("Error GetHistoricalPrices", err)
		return historicalPrices, err
	}
	err = _processResponse(resp, &historicalPrices)
	historicalPrices.HTTPCode = resp.StatusCode
	return historicalPrices, err
}

func (client *FtxClient) GetTrades(market string, limit int64, startTime int64, endTime int64) (Trades, error) {
	var trades Trades
	resp, err := client._get(
		"markets/"+market+"/trades?"+
			"&limit="+strconv.FormatInt(limit, 10)+
			"&start_time="+strconv.FormatInt(startTime, 10)+
			"&end_time="+strconv.FormatInt(endTime, 10),
		[]byte(""))
	if err != nil {
		log.Printf("Error GetTrades", err)
		return trades, err
	}
	err = _processResponse(resp, &trades)
	return trades, err
}
