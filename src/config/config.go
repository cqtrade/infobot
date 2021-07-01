package config

import "github.com/spf13/viper"

type Config interface {
	GetDiscordEnabled() bool
	GetDiscordChByChName(string) string
	GetServerUrl() string
}

type conf struct {
	discordEnabled       bool
	discordChBtcVibe     string
	discordChRandomIdeas string
	discordChAltSignals  string
	discordChFlash       string
	serverUrl            string
}

func New() Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	return &conf{
		discordEnabled:       viper.GetBool("DISCORD_ENABLED"),
		discordChBtcVibe:     viper.GetString("DISCORD_CH_BTC_VIBE"),
		discordChRandomIdeas: viper.GetString("DISCORD_CH_RANDOM_IDEAS"),
		discordChAltSignals:  viper.GetString("DISCORD_CH_ALTSIGNALS"),
		discordChFlash:       viper.GetString("DISCORD_CH_FLASH"),
		serverUrl:            viper.GetString("SERVER_URL"),
	}
}

func (c *conf) GetDiscordEnabled() bool {
	return c.discordEnabled
}

func (c *conf) GetDiscordChByChName(name string) string {
	if name == "btc-vibe" {
		return c.discordChBtcVibe
	}
	if name == "flash" {
		return c.discordChFlash
	}
	if name == "alt-signals" {
		return c.discordChAltSignals
	}
	if name == "random-ideas" {
		return c.discordChRandomIdeas
	}
	return ""
}

func (c *conf) GetServerUrl() string {
	return c.serverUrl
}
