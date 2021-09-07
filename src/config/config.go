package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DiscordEnabled       bool
	HealthLogEnabled     bool
	DiscordChRandomIdeas string
	DiscordChAltSignals  string
	DiscordChFlash       string
	DiscordChHealth      string
	DiscordChLogs        string
	ServerUrl            string
	FTXKey               string
	FTXSecret            string
	PositionSize         float64
	ProfitPercentage     float64
	FutureBTC            string
	FutureETH            string
	SubAccBTCDC          string
	SubAccETHDC          string
	RiskD                float64
	RiskDC               float64
}

func New() *Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	return &Config{
		DiscordEnabled:       viper.GetBool("DISCORD_ENABLED"),
		HealthLogEnabled:     viper.GetBool("HEALTH_LOG_ENABLED"),
		DiscordChRandomIdeas: viper.GetString("DISCORD_CH_RANDOM_IDEAS"),
		DiscordChAltSignals:  viper.GetString("DISCORD_CH_ALTSIGNALS"),
		DiscordChFlash:       viper.GetString("DISCORD_CH_FLASH"),
		DiscordChHealth:      viper.GetString("DISCORD_CH_HEALTH"),
		DiscordChLogs:        viper.GetString("DISCORD_CH_LOGS"),
		ServerUrl:            viper.GetString("SERVER_URL"),
		FTXKey:               viper.GetString("FTX_KEY"),
		FTXSecret:            viper.GetString("FTX_SECRET"),
		PositionSize:         viper.GetFloat64("POSITION_SIZE"),
		ProfitPercentage:     viper.GetFloat64("PROFIT_PERCENTAGE"),
		FutureBTC:            viper.GetString("FUTURE_BTC"),
		FutureETH:            viper.GetString("FUTURE_ETH"),
		SubAccBTCDC:          viper.GetString("SUBA_BTCDC"),
		SubAccETHDC:          viper.GetString("SUBA_ETHDC"),
		RiskD:                viper.GetFloat64("RISKD"),
		RiskDC:               viper.GetFloat64("RISKDC"),
	}
}
