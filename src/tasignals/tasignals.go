package tasignals

import (
	"fmt"
	"math"
	"time"

	"github.com/cqtrade/infobot/src/config"
	tvcontroller "github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/markcheno/go-talib"
)

func NtdfiValues(inClose []float64, inTimePeriod int) []float64 {
	lookback := inTimePeriod

	closes := make([]float64, len(inClose))
	tdfs := make([]float64, len(inClose))
	ntdfs := make([]float64, len(inClose))

	for i, close := range inClose {
		closes[i] = close * 10000000
	}

	mmas := talib.Ema(closes, lookback)
	smmas := talib.Ema(mmas, lookback)

	for i, _ := range mmas {
		tdfs[i] = 0
		ntdfs[i] = 0
		if i == 0 {
			continue
		}
		impetmma := mmas[i] - mmas[i-1]
		impetsmma := smmas[i] - smmas[i-1]
		divma := math.Abs(mmas[i] - smmas[i])
		averimpet := (impetmma + impetsmma) / 2
		number := averimpet
		pow := 3

		var result float64
		for i := 1; i <= pow-1; i++ {
			if i == 1 {
				result = number
			}
			result = result * number
		}

		tdfs[i] = divma * result
		if i < lookback*3 {
			ntdfs[i] = 0
			continue
		}
		highest := math.Abs(tdfs[i])
		for j := 0; j < lookback*3; j++ {
			if highest < math.Abs(tdfs[i-j]) {
				highest = math.Abs(tdfs[i-j])
			}
		}
		ntdfs[i] = tdfs[i] / highest
	}

	return ntdfs
}

type Indicators struct {
	Close        float64
	High         float64
	Low          float64
	Open         float64
	StartTime    time.Time
	Volume       float64
	Rsi          float64
	Roc          float64
	Mfi          float64
	Ntdf         float64
	NtdfSmoothed float64
}

func (ts *TaSignals) ta(candles []structs.HistoricalPrice) []Indicators {
	if len(candles) > 0 { // rm ongoing candle
		candles = candles[:len(candles)-1]
	}
	indicators := make([]Indicators, len(candles))
	closes := make([]float64, len(candles))
	highs := make([]float64, len(candles))
	lows := make([]float64, len(candles))
	volumes := make([]float64, len(candles))
	for i, candle := range candles {
		closes[i] = candle.Close
		highs[i] = candle.High
		lows[i] = candle.Low
		volumes[i] = candle.Volume
		indicators[i].Close = candle.Close
		indicators[i].High = candle.High
		indicators[i].Low = candle.Low
		indicators[i].Open = candle.Open
		indicators[i].Volume = candle.Volume
		indicators[i].StartTime = candle.StartTime
	}

	rsiPeriod := ts.cfg.TaRsiPeriod
	rocPeriod := ts.cfg.TaRocPeriod
	ntdfLookback := ts.cfg.TaNtdfLookback
	ntdfSmoothedPeriod := ts.cfg.TaNtdfSmoothedPeriod
	mfiPeriod := ts.cfg.TaMfiPeriod

	rsis := talib.Rsi(closes, rsiPeriod)
	rocs := talib.Roc(closes, rocPeriod)
	ntdfs := NtdfiValues(closes, ntdfLookback)
	ntdfisSmoothedRaw := talib.Ema(ntdfs, ntdfSmoothedPeriod)
	var ntdfisSmoothed []float64
	for _, x := range ntdfisSmoothedRaw {
		ntdfisSmoothed = append(ntdfisSmoothed, ftxtrade.Round(x, 2))
	}
	mfis := talib.Mfi(highs, lows, closes, volumes, mfiPeriod)

	for i, rsi := range rsis {
		if i < ntdfLookback {
			continue
		}
		indicators[i].Rsi = rsi
		indicators[i].Mfi = mfis[i]

		indicators[i].Ntdf = ntdfs[i]
		indicators[i].NtdfSmoothed = ntdfisSmoothed[i]
		indicators[i].Roc = rocs[i]
	}
	return indicators
}

func (ts *TaSignals) flash(indicators []Indicators) {
	for i, indicator := range indicators {
		if i < 10 {
			continue
		}

		rsiLow := ts.cfg.TaRsiLow
		lowMFI := ts.cfg.TaMfiLow
		lowRoc := ts.cfg.TaRocLow
		filterLow := ts.cfg.TaFilterLow
		filterHigh := ts.cfg.TaFilterHigh

		rsiBuy := indicator.Rsi < rsiLow || indicators[i-1].Rsi < rsiLow
		mfiBuy := indicators[i-1].Mfi < lowMFI && indicator.Mfi > indicators[i-1].Mfi
		rocBuy := indicators[i-1].Roc < lowRoc || indicators[i-2].Roc < lowRoc

		buy1 := rocBuy && rsiBuy && indicator.NtdfSmoothed < filterLow
		buy2 := (indicator.Rsi < 12 || indicators[i-1].Rsi < 12) && (indicator.Roc < -6 || indicators[i-1].Roc < -6) && indicator.NtdfSmoothed < filterLow
		buy3 := (indicators[i-1].Roc < lowRoc || indicators[i-2].Roc < lowRoc) && indicators[i-1].Rsi < 20 && indicator.Rsi > indicators[i-1].Rsi

		buyEnter := (((buy1 || buy2) && mfiBuy) || buy3)

		if buyEnter {
			println(fmt.Sprintf(
				"BUY %02d.%02d %02d:%02d\tROC %.2f\tRSI %.2f\tMFI %.3f\tNTDF %.3f\tNTDFsmoothed %.3f",
				indicator.StartTime.Day(),
				indicator.StartTime.Month(),
				indicator.StartTime.Hour()+3,
				indicator.StartTime.Minute(),
				indicator.Roc,
				indicator.Rsi,
				indicator.Mfi,
				indicator.Ntdf,
				indicator.NtdfSmoothed,
			))
		}

		exitBuy := (indicators[i-1].Rsi > 85 && indicator.Rsi < indicators[i-1].Rsi && indicator.NtdfSmoothed > 0.01) ||
			(indicators[i-1].NtdfSmoothed > filterHigh && indicators[i-2].NtdfSmoothed < indicators[i-1].NtdfSmoothed && indicators[i-1].NtdfSmoothed > indicator.NtdfSmoothed) ||
			(indicators[i-1].NtdfSmoothed > 0.85 && indicator.NtdfSmoothed < indicators[i-1].NtdfSmoothed)

		if exitBuy {
			println(fmt.Sprintf(
				"SELL %02d.%02d %02d:%02d\tROC %.2f\tRSI %.2f\tMFI %.3f\tNTDF %.3f\tNTDFsmoothed %.3f",
				indicator.StartTime.Day(),
				indicator.StartTime.Month(),
				indicator.StartTime.Hour()+3,
				indicator.StartTime.Minute(),
				indicator.Roc,
				indicator.Rsi,
				indicator.Mfi,
				indicator.Ntdf,
				indicator.NtdfSmoothed,
			))
		}
	}
}

type TaSignals struct {
	cfg    config.Config
	tvCtrl tvcontroller.TvController
}

func New(cfg config.Config, tvCtrl tvcontroller.TvController) *TaSignals {
	return &TaSignals{
		cfg:    cfg,
		tvCtrl: tvCtrl,
	}
}

func (ts *TaSignals) CheckFlashSignals() {
	key := ts.cfg.FTXKey
	secret := ts.cfg.FTXSecret
	client := ftx.New(key, secret, "")
	defer client.Client.CloseIdleConnections()

	var m15 int64
	m15 = 15 * 60
	candles, err := client.GetHistoricalPriceLatest("BTC/USD", m15, 5000)
	if err != nil {
		fmt.Println(err.Error())
	}

	if !candles.Success {
		fmt.Println(fmt.Sprintf("HERE1 %d ", candles.HTTPCode) + candles.ErrorMessage)
	} else {
		cndles := make([]structs.HistoricalPrice, len(candles.Result))
		for i, candle := range candles.Result {
			cndles[i] = candle
			cndles[i].StartTime = cndles[i].StartTime.Add(time.Minute * 15)
		}

		println(len(cndles))
		indicators := ts.ta(cndles)
		ts.flash(indicators)
	}
}
