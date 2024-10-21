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
		interfaceInputs := make([]interface{}, len(inputs))
		for i, v := range inputs {
			interfaceInputs[i] = v
		}
		return nextView(g, interfaceInputs)
	}

	prevField := func(g *gocui.Gui, v *gocui.View) error {
		interfaceInputs := make([]interface{}, len(inputs))
		for i, v := range inputs {
			interfaceInputs[i] = v
		}
		return prevView(g, interfaceInputs)
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

	// Create submit button
	submitButton := formComponents.NewButton(g, "Submit", maxX, maxY+12, 10)
	submitButton.Draw()

	// Add handler for submitting the form
	submitHandler := func(g *gocui.Gui) error {
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

	if err := g.SetKeybinding("Submit", gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return submitHandler(g)
	}); err != nil {
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
