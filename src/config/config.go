package config

import "github.com/spf13/viper"

type Config interface {
	GetDiscordEnabled() bool
	GetDiscordChBtcVibe() string
	GetDiscordChRandomIdeas() string
	GetDiscordChAltSignals() string
	GetDiscordChByChName(string) string
	GetServerUrl() string
}

type conf struct {
	discordEnabled       bool
	discordChBtcVibe     string
	discordChRandomIdeas string
	discordChAltSignals  string
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
		serverUrl:            viper.GetString("SERVER_URL"),
	}
}

func (c *conf) GetDiscordEnabled() bool {
	return c.discordEnabled
}

func (c *conf) GetDiscordChBtcVibe() string {
	return c.discordChBtcVibe
}

func (c *conf) GetDiscordChRandomIdeas() string {
	return c.discordChRandomIdeas
}

func (c *conf) GetDiscordChAltSignals() string {
	return c.discordChAltSignals
}

func (c *conf) GetDiscordChByChName(name string) string {
	if name == "btc-vibe" {
		return c.discordChBtcVibe
	}
	return ""
}

func (c *conf) GetServerUrl() string {
	return c.serverUrl
}
