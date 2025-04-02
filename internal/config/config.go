package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
	Database struct {
		Path string `toml:"path"`
	} `toml:"database"`
	MuchSender struct {
		Host         string `toml:"host"`
		Port         int    `toml:"port"`
		BearerToken  string `toml:"bearertoken"`
		ReplyToName  string `toml:"replytoname"`
		ReplyToEmail string `toml:"replytoemail"`
		Subject      string `toml:"subject"`
	} `toml:"muchsender"`
	GigaWallet struct {
		Host             string `toml:"host"`
		AdminPort        int    `toml:"adminport"`
		PubPort          int    `toml:"pubport"`
		AdminBearerToken string `toml:"adminbearertoken"`
		PubBearerToken   string `toml:"pubbearertoken"`
	} `toml:"gigawallet"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(data, &config)
	return config, err
}
