package utils

import (
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

func GetRandomProxy() (string, string, string, error) {

	proxyURLs := viper.GetStringSlice("proxy_urls")
	if len(proxyURLs) == 0 {
		return "", "", "", fmt.Errorf("no proxy URLs found in config")
	}

	return proxyURLs[rand.Intn(len(proxyURLs))], viper.GetString("proxy_username"), viper.GetString("proxy_password"), nil
}
