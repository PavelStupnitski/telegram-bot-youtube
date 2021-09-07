package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const pathConfig = "/home/pavel/go/src/telegram-bot-youtube/configs/config.yml"

type UniqueKeys struct {
	TelegramToken string `yaml:"telegram_token"`
	ConsumerKey   string `yaml:"consumer_key"`
}

func ReadConfigFile() (UniqueKeys, error) {
	var config UniqueKeys
	readByte, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		return UniqueKeys{}, err
	}
	err = yaml.Unmarshal(readByte, &config)
	if err != nil {
		return UniqueKeys{}, err
	}
	return config, nil
}
