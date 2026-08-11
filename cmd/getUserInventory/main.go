package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Sednan22/steam-inventory-go/pkg/exporter"
	inventory "github.com/Sednan22/steam-inventory-go/pkg/steaminventory"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	var format, output string
	flag.StringVar(&format, "format", "", "specify the format in which the data will be output")
	flag.StringVar(&output, "output", "", "specify the path where the data will be output")
	flag.Parse()

	args := flag.Args()

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

	userInv, err := inventory.GetUserInventory(steamID, game, contextID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if format == "" {
		if output != "" {
			ext := filepath.Ext(output)
			if ext == ".json" {
				format = "json"

			} else if ext == ".csv" {
				format = "csv"

			} else {
				format = "txt"
			}
		} else {
			fmt.Printf("Unsupported format '%s'. Using 'txt' instead.\n", format)
			format = "txt"
		}
	}

	data, err := exporter.FormatData(userInv, format)
	if err != nil {
		fmt.Printf("error formatting data to format %s: %v", format, err)
		return
	}

	if output != "" {
		err = exporter.DataToFile(format, output, data)
		if err != nil {
			fmt.Printf("error writing data to file: %v", err)
			return
		}

		fmt.Println("done!")
		return

	} else {
		fmt.Println(string(data))
		fmt.Println("done!")
		return
	}

}
