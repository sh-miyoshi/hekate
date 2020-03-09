package output

import (
	"fmt"

	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
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
		print.Print(out)
	case "json":
		out, err := data.ToJSON()
		if err != nil {
			return err
		}
		print.Print(out)
	}
	return nil
}
