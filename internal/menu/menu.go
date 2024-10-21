package menu

import (
	"github.com/WellyngtonF/WishListCLI/internal/formComponents"
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

func nextView(g *gocui.Gui, components []interface{}) error {
	currentView := g.CurrentView().Name()
	for i, comp := range components {
		var label string
		switch v := comp.(type) {
		case *formComponents.InputField:
			label = v.GetLabel()
		case *formComponents.Button:
			label = v.GetLabel()
		}
		if label == currentView {
			nextIndex := (i + 1) % len(components)
			g.SetCurrentView(getLabel(components[nextIndex]))
			return nil
		}
	}
	return nil
}

func getLabel(component interface{}) string {
	switch v := component.(type) {
	case *formComponents.InputField:
		return v.GetLabel()
	case *formComponents.Button:
		return v.GetLabel()
	}
	return ""
}

func prevView(g *gocui.Gui, components []interface{}) error {
	currentView := g.CurrentView().Name()
	for i, comp := range components {
		var label string
		switch v := comp.(type) {
		case *formComponents.InputField:
			label = v.GetLabel()
		case *formComponents.Button:
			label = v.GetLabel()
		}
		if label == currentView {
			prevIndex := (i - 1 + len(components)) % len(components)
			g.SetCurrentView(getLabel(components[prevIndex]))
			return nil
		}
	}
	return nil
}
