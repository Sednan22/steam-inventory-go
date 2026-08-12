package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Sednan22/steam-inventory-go/pkg/exporter"
	inventory "github.com/Sednan22/steam-inventory-go/pkg/steaminventory"
	"github.com/joho/godotenv"
)

// DefaultRateLimit defines the delay between consecutive requests to the Steam API.
// It can be adjusted depending on application requirements or rate-limiting rules.
var DefaultRateLimit = 5 * time.Second

func main() {

	_ = godotenv.Load()

	var allInventories map[string]map[string]int

	var format, output string
	var multi bool
	flag.StringVar(&format, "format", "", "specify the format in which the data will be output")
	flag.StringVar(&output, "output", "", "specify the path where the data will be output")
	flag.BoolVar(&multi, "multi", false, "specify if you want multi fetch")
	flag.Parse()

	args := flag.Args()

	if multi {

		if len(args) != 2 {
			fmt.Println("usage: ./getUserInventory -multi=true [flags] <appID> <contextID>")
			return
		}

		game, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("error converting appID to int: %v\n", err)
			return
		}
		contextID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("error converting contextID to int: %v\n", err)
			return
		}

		usersIDs, err := exporter.GetAllIDsFromFile()
		if err != nil {
			fmt.Printf("error getting usersIDs: %v\n", err)
			return
		}

		type inventoryResult struct {
			UserID string
			Data   map[string]int
			Err    error
		}

		resultsChan := make(chan inventoryResult, len(usersIDs))

		limiter := time.NewTicker(DefaultRateLimit)
		defer limiter.Stop()

		for i, userId := range usersIDs {

			go func(id string) {

				singleInv, err := inventory.GetUserInventory(id, game, contextID)

				resultsChan <- inventoryResult{
					UserID: id,
					Data:   singleInv,
					Err:    err,
				}
			}(userId)

			if i < len(usersIDs)-1 {
				<-limiter.C
			}

		}

		allInventories = make(map[string]map[string]int)
		for i := 0; i < len(usersIDs); i++ {

			res := <-resultsChan

			if res.Err != nil {
				fmt.Printf("Error getting %s inventory: %v\n", res.UserID, res.Err)
				continue
			}

			allInventories[res.UserID] = res.Data

		}

	} else {

		if len(args) != 3 {
			fmt.Println("usage: ./getUserInventory [flags] <steamID64> <appID> <contextID>")
			return
		}

		steamID := args[0]

		game, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("error converting appID to int: %v\n", err)
			return
		}
		contextID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("error converting contextID to int: %v\n", err)
			return
		}

		singleInv, err := inventory.GetUserInventory(steamID, game, contextID)
		if err != nil {
			fmt.Println(err)
			return
		}

		allInventories = map[string]map[string]int{
			steamID: singleInv,
		}

	}

	if format == "" {
		format = getExtension(output)
	}

	data, err := exporter.FormatData(allInventories, format)
	if err != nil {
		fmt.Printf("error formatting data to format %s: %v\n", format, err)
		return
	}

	if output != "" {
		err = exporter.DataToFile(format, output, data)
		if err != nil {
			fmt.Printf("error writing data to file: %v\n", err)
			return
		}
	} else {
		fmt.Println(string(data))
	}

	fmt.Println("done")

}

func getExtension(output string) string {

	if output == "" {
		return "txt"
	}

	switch filepath.Ext(output) {
	case ".json":
		return "json"
	case ".csv":
		return "csv"
	default:
		return "txt"
	}

}
