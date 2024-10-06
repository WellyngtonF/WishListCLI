package sources

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/WellyngtonF/WishListCLI/internal/scraper/utils"
	"github.com/gocolly/colly"
)

func ScrapeMercadoLivre(item item.Item) (float64, string, error) {
	searchURL := fmt.Sprintf("https://lista.mercadolivre.com.br/%s", strings.ReplaceAll(item.Name, " ", "-"))

	c := colly.NewCollector(
		colly.AllowedDomains("www.mercadolivre.com.br", "lista.mercadolivre.com.br"),
	)

	proxyURL, username, password, err := utils.GetRandomProxy()
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

		if price < item.MinPrice {
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
