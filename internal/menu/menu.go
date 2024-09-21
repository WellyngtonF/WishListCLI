package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/WellyngtonF/WishListCLI/internal/repository"
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
		fmt.Println("Web Scraping is not yet implemented.")
	case 6:
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid option. Please choose again.")
	}
}

// Helper function to add an item to the wishlist
func handleAddItem() {
	var name, category, producer, scrapingSources string
	var maxPrice float64

	fmt.Print("Enter item name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()

	fmt.Print("Enter category: ")
	fmt.Scanln(&category)

	fmt.Print("Enter producer: ")
	fmt.Scanln(&producer)

	fmt.Print("Enter max price (e.g., 99.99): ")
	fmt.Scanln(&maxPrice)

	fmt.Print("Enter scraping sources (comma-separated): ")
	fmt.Scanln(&scrapingSources)

	newItem := item.Item{
		Name:            name,
		Category:        category,
		Producer:        producer,
		MaxPrice:        maxPrice,
		ScrapingSources: parseScrapingSources(scrapingSources),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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

	fmt.Println("===== Wishlist =====")
	for _, itm := range items {
		fmt.Printf("Name: %s | Category: %s | Producer: %s | Max Price: %.2f | Sources: %v\n",
			itm.Name, itm.Category, itm.Producer, itm.MaxPrice, itm.ScrapingSources)
	}
}

// Helper function to update an item
func handleUpdateItem() {
	var name, category, producer, scrapingSources string
	var maxPrice float64

	fmt.Print("Enter the name of the item to update: ")
	fmt.Scanln(&name)

	itemToUpdate, err := repository.ReadItem(name)
	if err != nil {
		fmt.Printf("Error fetching item: %v\n", err)
		return
	}

	fmt.Printf("Updating item: %s\n", itemToUpdate.Name)

	fmt.Print("Enter new category (leave blank to keep unchanged): ")
	fmt.Scanln(&category)
	if category != "" {
		itemToUpdate.Category = category
	}

	fmt.Print("Enter new producer (leave blank to keep unchanged): ")
	fmt.Scanln(&producer)
	if producer != "" {
		itemToUpdate.Producer = producer
	}

	fmt.Print("Enter new max price (leave blank to keep unchanged): ")
	fmt.Scanln(&maxPrice)
	if maxPrice != 0 {
		itemToUpdate.MaxPrice = maxPrice
	}

	fmt.Print("Enter new scraping sources (leave blank to keep unchanged): ")
	fmt.Scanln(&scrapingSources)
	if scrapingSources != "" {
		itemToUpdate.ScrapingSources = parseScrapingSources(scrapingSources)
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

// Helper function to parse scraping sources from a comma-separated string, split and trim the values
func parseScrapingSources(sources string) []string {
	scrapingSources := strings.Split(sources, ",")
	for i := range scrapingSources {
		scrapingSources[i] = strings.TrimSpace(scrapingSources[i])
	}
	return scrapingSources
}
