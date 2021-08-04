package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DiscordEnabled       bool
	DiscordChBtcVibe     string
	DiscordChRandomIdeas string
	DiscordChAltSignals  string
	DiscordChFlash       string
	DiscordChHealth      string
	serverUrl            string
	FTXKey               string
	FTXSecret            string
	PositionSize         float64
	ProfitPercentage     float64
	FutureBTC            string
	FutureETH            string
	SubAccBTCD           string
	SubAccETHD           string
}

func New() *Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	return &Config{
		DiscordEnabled:       viper.GetBool("DISCORD_ENABLED"),
		DiscordChBtcVibe:     viper.GetString("DISCORD_CH_BTC_VIBE"),
		DiscordChRandomIdeas: viper.GetString("DISCORD_CH_RANDOM_IDEAS"),
		DiscordChAltSignals:  viper.GetString("DISCORD_CH_ALTSIGNALS"),
		DiscordChFlash:       viper.GetString("DISCORD_CH_FLASH"),
		DiscordChHealth:      viper.GetString("DISCORD_CH_HEALTH"),
		serverUrl:            viper.GetString("SERVER_URL"),
		FTXKey:               viper.GetString("FTX_KEY"),
		FTXSecret:            viper.GetString("FTX_SECRET"),
		PositionSize:         viper.GetFloat64("POSITION_SIZE"),
		ProfitPercentage:     viper.GetFloat64("PROFIT_PERCENTAGE"),
		FutureBTC:            viper.GetString("FUTURE_BTC"),
		FutureETH:            viper.GetString("FUTURE_ETH"),
		SubAccBTCD:           viper.GetString("SUBA_BTCD"),
		SubAccETHD:           viper.GetString("SUBA_ETHD"),
	}
}

func (c *Config) GetDiscordEnabled() bool {
	return c.DiscordEnabled
}

func (c *Config) GetDiscordChByChName(name string) string {
	if name == "btc-vibe" {
		return c.DiscordChBtcVibe
	}
	if name == "flash" {
		return c.DiscordChFlash
	}
	if name == "alt-signals" {
		return c.DiscordChAltSignals
	}
	if name == "random-ideas" {
		return c.DiscordChRandomIdeas
	}
	return ""
}

func (c *Config) GetServerUrl() string {
	return c.serverUrl
}
