package scraper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/gocolly/colly"
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
		return scrapeMercadoLivre(item)
	// Add more cases for other sources here
	default:
		return 0, "", fmt.Errorf("unsupported source: %s", source)
	}
}

// ScrapeMercadoLivre scrapes the price of an item from Mercado Livre, returning the price, the URL of the item, and an error if any.
func scrapeMercadoLivre(item item.Item) (float64, string, error) {
	searchURL := fmt.Sprintf("https://lista.mercadolivre.com.br/%s", strings.ReplaceAll(item.Name, " ", "-"))

	c := colly.NewCollector(
		colly.AllowedDomains("www.mercadolivre.com.br", "lista.mercadolivre.com.br"),
	)

	proxyURL, username, password, err := getRandomProxy()
	if err != nil {
		fmt.Printf("Error getting proxy: %v\n", err)
	} else {
		err = c.SetProxy(fmt.Sprintf("http://%s:%s@%s", username, password, proxyURL))
		if err != nil {
			fmt.Printf("Error setting proxy: %v\n", err)
		}
	}

	c.UserAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0"

	type Product struct {
		Price float64
		URL   string
	}

	var products []Product
	var lowestPriceProduct Product

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	c.OnHTML("li.ui-search-layout__item", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		if len(products) >= 10 {
			return
		}

		url := e.ChildAttr("a.ui-search-link", "href")
		if url == "" {
			url = e.ChildAttr("a.ui-search-link__title-card", "href")
		}

		priceStr := e.ChildText("div.ui-search-price__second-line span.ui-search-price__part--medium span.andes-money-amount__fraction")
		priceStr = strings.TrimSpace(priceStr)
		priceStr = strings.ReplaceAll(priceStr, ".", "")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Printf("Error parsing price: %v\n", err)
			return
		}

		product := Product{Price: price, URL: url}
		products = append(products, product)

		if lowestPriceProduct.Price == 0 || price < lowestPriceProduct.Price {
			lowestPriceProduct = product
		}
	})

	err = c.Visit(searchURL)
	if err != nil {
		return 0, "", fmt.Errorf("error visiting URL: %v", err)
	}

	if len(products) == 0 {
		return 0, "", errors.New("no products found")
	}

	return lowestPriceProduct.Price, lowestPriceProduct.URL, nil
}

func getRandomProxy() (string, string, string, error) {

	proxyURLs := viper.GetStringSlice("proxy_urls")
	if len(proxyURLs) == 0 {
		return "", "", "", fmt.Errorf("no proxy URLs found in config")
	}

	return proxyURLs[rand.Intn(len(proxyURLs))], viper.GetString("proxy_username"), viper.GetString("proxy_password"), nil
}
