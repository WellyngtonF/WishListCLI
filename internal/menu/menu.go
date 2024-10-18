package menu

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/WellyngtonF/WishListCLI/internal/formComponents"
	"github.com/WellyngtonF/WishListCLI/internal/item"
	"github.com/WellyngtonF/WishListCLI/internal/repository"
	"github.com/awesome-gocui/gocui"
)

// GetMenuText returns the menu text as a string
func GetMenuText() string {
	return `1. Add Item to Wishlist
2. View Wishlist
3. Update Item in Wishlist
4. Delete Item from Wishlist
5. Run Web Scraping
6. Exit
Choose an option:`
}

func GetMenuOptions() []string {
	return []string{
		"Add Item to Wishlist",
		"View Wishlist",
		"Update Item in Wishlist",
		"Delete Item from Wishlist",
		"Run Web Scraping",
		"Exit",
	}
}

type InputField struct {
	name      string
	x, y      int
	w         int
	maxLength int
	label     string
}

func NewInputField(name string, x, y, w, maxLength int, label string) *InputField {
	return &InputField{name: name, x: x, y: y, w: w, maxLength: maxLength, label: label}
}

func (i *InputField) Layout(g *gocui.Gui) error {
	labelView, err := g.SetView(i.name+"Label", i.x, i.y, i.x+len(i.label)+1, i.y+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelView.Frame = false
		fmt.Fprint(labelView, i.label)
	}

	inputView, err := g.SetView(i.name, i.x+len(i.label)+1, i.y, i.x+i.w, i.y+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		inputView.Editable = true
		inputView.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			if v.BufferLines()[0] == "" || len(v.BufferLines()[0]) < i.maxLength {
				gocui.DefaultEditor.Edit(v, key, ch, mod)
			}
		})
	}
	return nil
}

func HandleAddItem(g *gocui.Gui, v *gocui.View) error {
	maxX := 30
	maxY := 0

	// Create input fields
	nameInput := formComponents.NewInputField(g, "Name", maxX, maxY, 10, 30).
		AddValidate("Name is required", func(value string) bool {
			return len(strings.TrimSpace(value)) > 0
		})

	categoryInput := formComponents.NewInputField(g, "Category", maxX, maxY+2, 10, 30)

	producerInput := formComponents.NewInputField(g, "Producer", maxX, maxY+4, 10, 30)

	maxPriceInput := formComponents.NewInputField(g, "Max Price", maxX, maxY+6, 10, 15).
		AddValidate("Invalid price format", func(value string) bool {
			_, err := strconv.ParseFloat(value, 64)
			return err == nil
		})

	minPriceInput := formComponents.NewInputField(g, "Min Price", maxX, maxY+8, 10, 15).
		AddValidate("Invalid price format", func(value string) bool {
			_, err := strconv.ParseFloat(value, 64)
			return err == nil
		})

	sourcesInput := formComponents.NewInputField(g, "Sources", maxX, maxY+10, 10, 40)

	inputs := []*formComponents.InputField{
		nameInput,
		categoryInput,
		producerInput,
		maxPriceInput,
		minPriceInput,
		sourcesInput,
	}

	// Draw input fields
	for _, input := range inputs {
		input.Draw()
	}

	// Set initial focus
	g.SetCurrentView("Name")

	// Add handler for navigating between fields
	nextField := func(g *gocui.Gui, v *gocui.View) error {
		return nextView(g, inputs)
	}

	prevField := func(g *gocui.Gui, v *gocui.View) error {
		return prevView(g, inputs)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextField); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, nextField); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, prevField); err != nil {
		return err
	}

	// Add handler for submitting the form
	submitHandler := func(g *gocui.Gui, v *gocui.View) error {
		for _, input := range inputs {
			if !input.Validate() {
				return nil
			}
		}

		maxPrice, _ := strconv.ParseFloat(maxPriceInput.GetFieldText(), 64)
		minPrice, _ := strconv.ParseFloat(minPriceInput.GetFieldText(), 64)

		newItem := item.Item{
			Name:            nameInput.GetFieldText(),
			Category:        categoryInput.GetFieldText(),
			Producer:        producerInput.GetFieldText(),
			MaxPrice:        maxPrice,
			MinPrice:        minPrice,
			ScrapingSources: strings.Split(sourcesInput.GetFieldText(), ","),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		err := repository.CreateItem(newItem)
		if err != nil {
			return fmt.Errorf("error adding item: %v", err)
		}

		// Close input fields
		for _, input := range inputs {
			input.Close()
		}

		g.SetCurrentView("menu")
		return nil
	}

	// Set keybinding for form submission
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, submitHandler); err != nil {
		return err
	}

	// Set keybinding for canceling the form
	cancelHandler := func(g *gocui.Gui, v *gocui.View) error {
		for _, input := range inputs {
			input.Close()
		}

		g.SetCurrentView("menu")
		return nil
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, cancelHandler); err != nil {
		return err
	}

	return nil
}

func nextView(g *gocui.Gui, views []*formComponents.InputField) error {
	currentView := g.CurrentView().Name()
	for i, view := range views {
		if view.GetLabel() == currentView {
			nextIndex := (i + 1) % len(views)
			g.SetCurrentView(views[nextIndex].GetLabel())
			return nil
		}
	}
	return nil
}

func prevView(g *gocui.Gui, views []*formComponents.InputField) error {
	currentView := g.CurrentView().Name()
	for i, view := range views {
		if view.GetLabel() == currentView {
			prevIndex := (i - 1 + len(views)) % len(views)
			g.SetCurrentView(views[prevIndex].GetLabel())
			return nil
		}
	}
	return nil
}

func createItem(g *gocui.Gui, v *gocui.View) error {
	var name, category, producer, maxPriceStr, minPriceStr, sources string
	var maxPrice, minPrice float64
	var err error

	if v, err := g.View("name"); err == nil {
		name = strings.TrimSpace(v.Buffer())
	}
	if v, err := g.View("category"); err == nil {
		category = strings.TrimSpace(v.Buffer())
	}
	if v, err := g.View("producer"); err == nil {
		producer = strings.TrimSpace(v.Buffer())
	}
	if v, err := g.View("maxPrice"); err == nil {
		maxPriceStr = strings.TrimSpace(v.Buffer())
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return fmt.Errorf("invalid max price: %v", err)
		}
	}
	if v, err := g.View("minPrice"); err == nil {
		minPriceStr = strings.TrimSpace(v.Buffer())
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return fmt.Errorf("invalid min price: %v", err)
		}
	}
	if v, err := g.View("sources"); err == nil {
		sources = strings.TrimSpace(v.Buffer())
	}

	newItem := item.Item{
		Name:            name,
		Category:        category,
		Producer:        producer,
		MaxPrice:        maxPrice,
		MinPrice:        minPrice,
		ScrapingSources: strings.Split(sources, ","),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = repository.CreateItem(newItem)
	if err != nil {
		return fmt.Errorf("error adding item: %v", err)
	}

	g.DeleteView("addItem")
	g.SetCurrentView("menu")
	return nil
}
