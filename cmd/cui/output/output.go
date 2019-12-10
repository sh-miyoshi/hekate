package output

import (
	"fmt"
	"encoding/json"
)

var outputFormat string

// Init ...
func Init(format string) error {
	if format == "json" || format == "text" {
		outputFormat = format
		return nil
	}
	return fmt.Errorf("Unknown output format: %s", format)
}

// Print ...
func Print(data interface{}) error {
	switch outputFormat {
	case "text":
		fmt.Printf("%v\n", data)
	case "json":
		out, err := json.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(out))
	}
	return nil
}