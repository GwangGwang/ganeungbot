package config

import (
	"github.com/GwangGwang/ganeungbot/pkg/util"
)

const TelegramAPIKey string = "telegram"
const TelegramConsoleChatID string = "telegram-console"
const WeatherAPIKey string = "weatherAPIkey"

// Map of config name to location
var configInfos map[string]string {
	TelegramAPIKey: "/secrets/telegram",
	TelegramConsoleChatID: "/secrets/telegram-consoleChatId",
	WeatherAPIKey: "/secrets/weatherAPIKey",
}

type ConfigMap map[string]string

func Get() ConfigMap {
	var configMap ConfigMap = make(map[string]string)

	for name, location := range configInfos {
		value, err := util.FileReadString(location)

		if err != nil {
			log.Printf("Config \"%s\" not found at \"%s\"\n", name, location)
		}

		configMap[name] = value
	}

	return configMap
}
