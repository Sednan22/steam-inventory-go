package main

import (
	"fmt"
	"os"
	"strconv"

	inventory "github.com/Sednan22/steam-inventory-go/pkg/steaminventory"
)

func main() {

	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("usage: ./getUserInventory <steamID64> <appID> <contextID>")
		return
	}

	// steamID64 for user you want to fetch
	// number id for game you want to fetch (example 730 is cs2 and 440 is tf2, etc...)
	// contextid usually 2 (cs2 items 2 is tradable 16 is untradable)

	steamID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error converting string to int: %v\n", err)
		return
	}
	game, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Error converting string to int: %v\n", err)
		return
	}
	contextID, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("Error converting string to int: %v\n", err)
		return
	}

	userInv, err := inventory.GetUserInventory(steamID, game, contextID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, value := range userInv {

		// You can check for repeated items with an if
		if value >= 3 {
			fmt.Printf("%s -> %d\n", key, value)
		}

		// Or print all items
		fmt.Printf("%s -> %d\n", key, value)
	}

	fmt.Println("Done!")
}
