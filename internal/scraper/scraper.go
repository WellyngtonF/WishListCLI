package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/WellyngtonF/WishListCLI/internal/scraper/sources"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
	}

	rand.Seed(uint64(time.Now().UnixNano()))
}

func ScrapePrice(item item.Item, source string) (float64, string, error) {
	switch strings.TrimSpace(strings.ToLower(source)) {
	case "mercado livre":
		return sources.ScrapeMercadoLivre(item)
	case "amazon":
		return sources.ScrapeAmazon(item)
	default:
		return 0, "", fmt.Errorf("unsupported source: %s", source)
	}
}
