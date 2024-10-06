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

func ScrapeAmazon(item item.Item) (float64, string, error) {
	// https://www.zoom.com.br/search?q=ps5&hitsPerPage=24&refinements%5B0%5D%5Bid%5D=bestSellingMerchantName&refinements%5B0%5D%5Bvalues%5D%5B0%5D=Amazon&sortBy=default&enableRefinementsSuggestions=true&isDealsPage=false
	searchURL := fmt.Sprintf("https://www.zoom.com.br/search?q=%s&refinements%%5B0%%5D%%5Bid%%5D=bestSellingMerchantName&refinements%%5B0%%5D%%5Bvalues%%5D%%5B0%%5D=Amazon&sortBy=default&enableRefinementsSuggestions=true&isDealsPage=false", item.Name)

	c := colly.NewCollector(
		colly.IgnoreRobotsTxt(),
		colly.AllowedDomains("www.zoom.com.br"),
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

	c.OnHTML("div[data-testid='product-card']", func(e *colly.HTMLElement) {

		if len(products) >= 10 {
			return
		}

		price := e.ChildText("p[data-testid='product-card::price']")
		price = strings.ReplaceAll(price, "R$", "")
		price = strings.TrimSpace(price)
		price = strings.ReplaceAll(price, ".", "")
		price = strings.ReplaceAll(price, ",", ".")
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("Error parsing price:", err)
		}

		if priceFloat < item.MinPrice {
			return
		}

		url := e.ChildAttr("a.ProductCard_ProductCard_Inner__gapsh", "href")
		url = "https://www.zoom.com.br" + url
		product := Product{Price: priceFloat, URL: url}
		products = append(products, product)

		if lowestPriceProduct.Price == 0 || priceFloat < lowestPriceProduct.Price {
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
