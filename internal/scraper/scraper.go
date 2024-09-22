package scraper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/gocolly/colly"
)

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

	type Product struct {
		Price float64
		URL   string
	}

	var products []Product
	var lowestPriceProduct Product

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

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

		product := Product{Price: price, URL: url}
		products = append(products, product)

		if lowestPriceProduct.Price == 0 || price < lowestPriceProduct.Price {
			lowestPriceProduct = product
		}
	})

	err := c.Visit(searchURL)
	if err != nil {
		return 0, "", fmt.Errorf("error visiting URL: %v", err)
	}

	if len(products) == 0 {
		return 0, "", errors.New("no products found")
	}

	return lowestPriceProduct.Price, lowestPriceProduct.URL, nil
}
