package main

import (
	"fmt"

	"github.com/WellyngtonF/WishListCLI/internal/menu"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	clearScreen()
	for {
		menu.ShowMenu()        // Show the menu
		menu.HandleMenuInput() // Handle the user's input

		// Wait for user input before clearing the screen
		fmt.Print("\nPress Enter to continue...")
		fmt.Scanln()

		clearScreen() // Clear the screen before showing the menu again
	}
}
