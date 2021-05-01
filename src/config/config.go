package config

import "github.com/spf13/viper"

type Config interface {
	GetDiscordChBtcVibe() string
	GetDiscordChRandomIdeas() string
	GetDiscordChByChName(string) string
	GetServerUrl() string
}

type conf struct {
	discordChBtcVibe     string
	discordChRandomIdeas string
	serverUrl            string
}

func New() Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	return &conf{
		discordChBtcVibe:     viper.GetString("DISCORD_CH_BTC_VIBE"),
		discordChRandomIdeas: viper.GetString("DISCORD_CH_RANDOM_IDEAS"),
		serverUrl:            viper.GetString("SERVER_URL"),
	}
}

func (c *conf) GetDiscordChBtcVibe() string {
	return c.discordChBtcVibe
}

func (c *conf) GetDiscordChRandomIdeas() string {
	return c.discordChRandomIdeas
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
