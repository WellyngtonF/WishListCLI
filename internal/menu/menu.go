package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/WellyngtonF/WishListCLI/internal/repository"
	"github.com/WellyngtonF/WishListCLI/internal/scraper"
)

// ShowMenu displays the CLI menu for the application
func ShowMenu() {
	fmt.Println("===== Wishlist Manager =====")
	fmt.Println("1. Add Item to Wishlist")
	fmt.Println("2. View Wishlist")
	fmt.Println("3. Update Item in Wishlist")
	fmt.Println("4. Delete Item from Wishlist")
	fmt.Println("5. Run Web Scraping")
	fmt.Println("6. Exit")
	fmt.Print("Choose an option: ")
}

// HandleMenuInput handles the menu options selected by the user
func HandleMenuInput() {
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		handleAddItem()
	case 2:
		handleViewWishlist()
	case 3:
		handleUpdateItem()
	case 4:
		handleDeleteItem()
	case 5:
		handleWebScraping()
	case 6:
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid option. Please choose again.")
	}
}

func handleWebScraping() {
	ListItemsNames()
	fmt.Print("Enter the name of the item to scrape: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	itemName := scanner.Text()

	item, err := repository.ReadItem(itemName)
	if err != nil {
		fmt.Printf("Error fetching item: %v\n", err)
		return
	}

	fmt.Printf("Scraping prices for %s...\n", item.Name)
	for _, source := range item.ScrapingSources {
		price, url, err := scraper.ScrapePrice(*item, source)
		if err != nil {
			fmt.Printf("Error scraping %s: %v\n", source, err)
			continue
		}
		fmt.Printf("%s: $%.2f - %s\n", source, price, url)
	}
}

// Helper function to add an item to the wishlist
func handleAddItem() {
	var name, category, producer, scrapingSources string
	var maxPrice, minPrice float64

	fmt.Print("Enter item name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()

	fmt.Print("Enter category: ")
	scanner.Scan()
	category = scanner.Text()

	fmt.Print("Enter producer: ")
	scanner.Scan()
	producer = scanner.Text()

	fmt.Print("Enter max price (e.g., 99.99): ")
	fmt.Scanln(&maxPrice)

	fmt.Print("Enter scraping sources (comma-separated): ")
	scanner.Scan()
	scrapingSources = scanner.Text()

	// Get minimum price (optional)
	fmt.Print("Enter minimum price (e.g., 99.99): ")
	fmt.Scanln(&minPrice)
	newItem := item.Item{
		Name:            name,
		Category:        category,
		Producer:        producer,
		MaxPrice:        maxPrice,
		ScrapingSources: parseScrapingSources(scrapingSources),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		MinPrice:        minPrice,
	}

	err := repository.CreateItem(newItem)
	if err != nil {
		fmt.Printf("Error adding item: %v\n", err)
		return
	}
	fmt.Println("Item added successfully!")
}

// Helper function to view the wishlist
func handleViewWishlist() {
	items, err := repository.ListItems()
	if err != nil {
		fmt.Printf("Error loading wishlist: %v\n", err)
		return
	}

	// Print table header
	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "Name", "Category", "Producer", "Max Price", "Scraping Sources")
	fmt.Println(strings.Repeat("-", 85))

	// Print table rows
	for _, item := range items {
		fmt.Printf("%-20s %-15s %-15s $%-9.2f %-20s\n",
			truncateString(item.Name, 20),
			truncateString(item.Category, 15),
			truncateString(item.Producer, 15),
			item.MaxPrice,
			truncateString(strings.Join(item.ScrapingSources, ", "), 20))
	}

	fmt.Println(strings.Repeat("-", 85))
	fmt.Printf("Total items: %d\n", len(items))
}

// Helper function to truncate strings that are too long
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Helper function to update an item
func handleUpdateItem() {
	var name, category, producer, scrapingSources string
	var maxPrice, minPrice float64

	scanner := bufio.NewScanner(os.Stdin)
	ListItemsNames()
	fmt.Print("Enter the name of the item to update: ")
	scanner.Scan()
	name = scanner.Text()

	itemToUpdate, err := repository.ReadItem(name)

	if err != nil {
		fmt.Printf("Error fetching item: %v\n", err)
		return
	}

	fmt.Println("\nCurrent item details:")
	fmt.Printf("Name: %s\n", itemToUpdate.Name)
	fmt.Printf("Category: %s\n", itemToUpdate.Category)
	fmt.Printf("Producer: %s\n", itemToUpdate.Producer)
	fmt.Printf("Max Price: $%.2f\n", itemToUpdate.MaxPrice)
	fmt.Printf("Scraping Sources: %s\n", strings.Join(itemToUpdate.ScrapingSources, ", "))
	fmt.Printf("Min Price: $%.2f\n", itemToUpdate.MinPrice)
	fmt.Println()

	fmt.Print("Enter new category (leave blank to keep unchanged): ")
	scanner.Scan()
	category = scanner.Text()
	if category != "" {
		itemToUpdate.Category = category
	}

	fmt.Print("Enter new producer (leave blank to keep unchanged): ")
	scanner.Scan()
	producer = scanner.Text()
	if producer != "" {
		itemToUpdate.Producer = producer
	}

	fmt.Print("Enter new max price (leave blank to keep unchanged): ")
	fmt.Scanln(&maxPrice)
	if maxPrice != 0 {
		itemToUpdate.MaxPrice = maxPrice
	}

	fmt.Print("Enter new scraping sources (leave blank to keep unchanged): ")
	scanner.Scan()
	scrapingSources = scanner.Text()
	if scrapingSources != "" {
		itemToUpdate.ScrapingSources = parseScrapingSources(scrapingSources)
	}

	fmt.Print("Enter new minimum price (leave blank to keep unchanged): ")
	fmt.Scanln(&minPrice)
	if minPrice != 0 {
		itemToUpdate.MinPrice = minPrice
	}

	itemToUpdate.UpdatedAt = time.Now()

	err = repository.UpdateItem(*itemToUpdate)
	if err != nil {
		fmt.Printf("Error updating item: %v\n", err)
		return
	}
	fmt.Println("Item updated successfully!")
}

// Helper function to delete an item from the wishlist
func handleDeleteItem() {
	var name string
	ListItemsNames()

	fmt.Print("Enter the name of the item to delete: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()

	err := repository.DeleteItem(name)
	if err != nil {
		fmt.Printf("Error deleting item: %v\n", err)
		return
	}
	fmt.Println("Item deleted successfully!")
}

func ListItemsNames() {
	names, err := repository.GetItemsNames()
	if err != nil {
		fmt.Printf("Error fetching items: %v\n", err)
		return
	}

	fmt.Println("Available items:")
	for _, name := range names {
		fmt.Println(name)
	}

}

// Helper function to parse scraping sources from a comma-separated string, split and trim the values
func parseScrapingSources(sources string) []string {
	scrapingSources := strings.Split(sources, ",")
	for i := range scrapingSources {
		scrapingSources[i] = strings.TrimSpace(scrapingSources[i])
	}
	return scrapingSources
}
