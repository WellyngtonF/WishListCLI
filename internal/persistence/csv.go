package persistence

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/item"
)

// Helper function to parse float from string
func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// Helper function to format float to string
func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

// Helper function to parse time from string
func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

// Helper function to format time to string
func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func CreateFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

// LoadItems loads the CSV file and returns all items
func LoadItems(filePath string) ([]item.Item, error) {
	CreateFile(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Use semicolon as column separator
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var items []item.Item
	for _, record := range records {
		maxPrice, err := parseFloat(record[3])
		if err != nil {
			return nil, err
		}

		createdAt, err := parseTime(record[5])
		if err != nil {
			return nil, err
		}

		updatedAt, err := parseTime(record[6])
		if err != nil {
			return nil, err
		}

		minPrice, err := parseFloat(record[7])
		if err != nil {
			return nil, err
		}

		items = append(items, item.Item{
			Name:            record[0],
			Category:        record[1],
			Producer:        record[2],
			MaxPrice:        maxPrice,
			ScrapingSources: strings.Split(record[4], ","),
			CreatedAt:       createdAt,
			UpdatedAt:       updatedAt,
			MinPrice:        minPrice,
		})
	}

	return items, nil
}

// AddItem appends a new item to the CSV file
func AddItem(filePath string, newItem item.Item) error {
	CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	err = writer.Write([]string{
		newItem.Name,
		newItem.Category,
		newItem.Producer,
		formatFloat(newItem.MaxPrice),
		strings.Join(newItem.ScrapingSources, ","),
		formatTime(newItem.CreatedAt),
		formatTime(newItem.UpdatedAt),
		formatFloat(newItem.MinPrice),
	})
	if err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// UpdateItem updates an existing item in the CSV file
func UpdateItem(filePath string, updatedItem item.Item) error {
	items, err := LoadItems(filePath)
	if err != nil {
		return err
	}

	found := false
	for i, itm := range items {
		if itm.Name == updatedItem.Name {
			items[i] = updatedItem
			found = true
			break
		}
	}

	if !found {
		return errors.New("item not found")
	}

	return saveItems(filePath, items)
}

// DeleteItem deletes an item from the CSV file
func DeleteItem(filePath string, name string) error {
	items, err := LoadItems(filePath)
	if err != nil {
		return err
	}

	var updatedItems []item.Item
	for _, itm := range items {
		if itm.Name != name {
			updatedItems = append(updatedItems, itm)
		}
	}

	return saveItems(filePath, updatedItems)
}

// Helper function to save items to the CSV file
func saveItems(filePath string, items []item.Item) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	for _, itm := range items {
		err = writer.Write([]string{
			itm.Name,
			itm.Category,
			itm.Producer,
			formatFloat(itm.MaxPrice),
			strings.Join(itm.ScrapingSources, ","),
			formatTime(itm.CreatedAt),
			formatTime(itm.UpdatedAt),
		})
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}
