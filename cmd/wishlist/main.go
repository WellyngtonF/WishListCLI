package main

import (
	"fmt"
	"log"

	"github.com/WellyngtonF/WishListCLI/internal/menu"
	"github.com/awesome-gocui/gocui"
)

const (
	menuViewName = "menu"
	mainViewName = "main"
)

var currentSelection = 0

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Mouse = true
	g.SetManagerFunc(layout)

	if err := setKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	menuWidth := 30

	if v, err := g.SetView(menuViewName, 0, 0, menuWidth-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Menu"
		updateMenuView(g)
	}

	if v, err := g.SetView(mainViewName, menuWidth, 0, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Wishlist Manager"
		fmt.Fprintln(v, "Welcome to Wishlist Manager!")
		fmt.Fprintln(v, "Click on a menu option to get started.")
	}

	return nil
}

func updateMenuView(g *gocui.Gui) error {
	v, err := g.View(menuViewName)
	if err != nil {
		return err
	}
	v.Clear()
	for _, option := range menu.GetMenuOptions() {
		if option == menu.GetMenuOptions()[currentSelection] {
			fmt.Fprintf(v, "> %s\n", option)
		} else {
			fmt.Fprintf(v, "%s\n", option)
		}
	}
	return nil
}

func setKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, onClick); err != nil {
		return err
	}

	return nil
}

func onClick(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() != menuViewName {
		return nil
	}

	_, cy := v.Cursor()
	if cy < len(menu.GetMenuOptions()) {
		currentSelection = cy
		updateMenuView(g)
		return handleMenuOption(g, cy)
	}
	return nil
}

func handleMenuOption(g *gocui.Gui, choice int) error {
	mainView, err := g.View(mainViewName)
	if err != nil {
		return err
	}

	mainView.Clear()

	switch choice {
	case 0:
		return menu.HandleAddItem(g, mainView)
	case 1:
		mainView.Title = "View Wishlist"
		fmt.Fprintln(mainView, "View Wishlist")
		// Implement view wishlist functionality
	case 2:
		mainView.Title = "Update Item in Wishlist"
		fmt.Fprintln(mainView, "Update Item in Wishlist")
		// Implement update item functionality
	case 3:
		mainView.Title = "Delete Item from Wishlist"
		fmt.Fprintln(mainView, "Delete Item from Wishlist")
		// Implement delete item functionality
	case 4:
		mainView.Title = "Run Web Scraping"
		fmt.Fprintln(mainView, "Run Web Scraping")
		// Implement web scraping functionality
	case 5:
		return gocui.ErrQuit
	default:
		fmt.Fprintln(mainView, "Invalid option. Please choose again.")
	}

	return nil
}
