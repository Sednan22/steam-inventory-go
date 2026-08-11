package exporter

import (
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

func FormatData(userInv map[string]int, format string) ([]byte, error) {

	var data []byte
	switch format {
	case "json":
		data, err := json.MarshalIndent(userInv, "", " ")
		if err != nil {
			return nil, fmt.Errorf("error marshalling data to json: %w", err)
		}
		return data, nil

	case "csv":
		var b strings.Builder
		fmt.Fprintf(&b, "Item, Quantity\n")
		for key, value := range userInv {
			fmt.Fprintf(&b, "%q,%d\n", key, value)
		}
		data = []byte(b.String())
		return data, nil

	default:
		var b strings.Builder
		for key, value := range userInv {
			fmt.Fprintf(&b, "%s -> %d\n", key, value)
		}
		data = []byte(b.String())
		return data, nil
	}

}
