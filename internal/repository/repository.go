package repository

import (
	"errors"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/WellyngtonF/WishListCLI/internal/persistence"
)

const filePath = "wishlist.csv" // Define the file path for CSV

// CreateItem adds a new item to the wishlist
func CreateItem(newItem item.Item) error {
	items, err := persistence.LoadItems(filePath)
	if err != nil {
		return err
	}

	// Check if item already exists
	for _, itm := range items {
		if itm.Name == newItem.Name {
			return errors.New("item already exists")
		}
	}

	// Set timestamps
	newItem.CreatedAt = time.Now()
	newItem.UpdatedAt = time.Now()

	return persistence.AddItem(filePath, newItem)
}

// ReadItem fetches an item by name
func ReadItem(name string) (*item.Item, error) {
	items, err := persistence.LoadItems(filePath)
	if err != nil {
		return nil, err
	}

	for _, itm := range items {
		if itm.Name == name {
			return &itm, nil
		}
	}

	return nil, errors.New("item not found")
}

// UpdateItem modifies an existing item
func UpdateItem(updatedItem item.Item) error {
	updatedItem.UpdatedAt = time.Now()
	return persistence.UpdateItem(filePath, updatedItem)
}

// DeleteItem removes an item from the wishlist
func DeleteItem(name string) error {
	return persistence.DeleteItem(filePath, name)
}

// ListItems returns all items in the wishlist
func ListItems() ([]item.Item, error) {
	return persistence.LoadItems(filePath)
}

func GetItemsNames() ([]string, error) {
	items, err := persistence.LoadItems(filePath)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(items))
	for i, itm := range items {
		names[i] = itm.Name
	}
	return names, nil
}
