package output

import (
	"fmt"
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
func Print(data Format) error {
	switch outputFormat {
	case "text":
		out, err := data.ToText()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out)
	case "json":
		out, err := data.ToJSON()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out)
	}
	return nil
}