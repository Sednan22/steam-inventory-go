package exporter

import (
	"testing"
)

func TestFormatData(t *testing.T) {

	tests := map[string]struct {
		inventories    map[string]map[string]int
		format         string
		expectedOutput string
		expectError    bool
	}{
		"valid json format": {
			inventories: map[string]map[string]int{
				"76561198000000000": {
					"AK-47": 2,
				},
			},
			format:         "json",
			expectedOutput: "{\n \"76561198000000000\": {\n  \"AK-47\": 2\n }\n}",
			expectError:    false,
		},
		"unsupported_format_error": {
			inventories: map[string]map[string]int{
				"76561198000000000": {"AK-47": 2},
			},
			format:         "xml",
			expectedOutput: "76561198000000000 inventory:\nAK-47 -> 2\n",
			expectError:    false,
		},
		"valid_csv_format": {
			inventories: map[string]map[string]int{
				"76561198000000000": {"AK-47": 2},
			},
			format:         "csv",
			expectedOutput: "SteamID, Item, Quantity\n76561198000000000,\"AK-47\",2\n",
			expectError:    false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := FormatData(test.inventories, test.format)

			if (err != nil) != test.expectError {
				t.Errorf("FormatData() error unexpected = %v, expected error? %v", err, test.expectError)
				return
			}

			if !test.expectError && string(output) != test.expectedOutput {
				t.Errorf("FormatData()\nGOT:      %q\nEXPECTED: %q", string(output), test.expectedOutput)
			}
		})
	}
}
