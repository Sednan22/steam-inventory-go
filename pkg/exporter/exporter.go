package exporter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DataToFile(format, output string, data []byte) error {

	var endFile string

	ext := getOutputExtension(output)
	if ext == "" {
		endFile = fmt.Sprintf("%s.%s", output, format)
	} else {
		endFile = output
	}

	err := os.WriteFile(endFile, data, 0o666)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func getOutputExtension(output string) (extension string) {
	return filepath.Ext(output)
}

func FormatData(allInventories map[string]map[string]int, format string) ([]byte, error) {

	var data []byte
	switch format {
	case "json":
		data, err := json.MarshalIndent(allInventories, "", " ")
		if err != nil {
			return nil, fmt.Errorf("error marshalling data to json: %w", err)
		}
		return data, nil

	case "csv":
		var b strings.Builder
		fmt.Fprintf(&b, "SteamID, Item, Quantity\n")
		for user, inv := range allInventories {
			for key, value := range inv {
				fmt.Fprintf(&b, "%s,%q,%d\n", user, key, value)
			}
		}
		data = []byte(b.String())
		return data, nil

	default:
		var b strings.Builder
		for user, inv := range allInventories {
			fmt.Fprintf(&b, "%s inventory:\n", user)
			for key, value := range inv {
				fmt.Fprintf(&b, "%s -> %d\n", key, value)
			}
		}
		data = []byte(b.String())
		return data, nil
	}

}

func GetAllIDsFromFile() ([]string, error) {

	file, err := os.Open("usersList.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading usersIDs from file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var usersIDs []string
	for scanner.Scan() {

		line := scanner.Text()

		trimed := strings.TrimSpace(line)

		if trimed == "" {
			continue
		}

		usersIDs = append(usersIDs, trimed)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return usersIDs, nil
}
