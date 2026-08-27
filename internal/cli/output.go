package cli

import (
	"encoding/json"
	"flag"
	"fmt"
)

func printJSON(value interface{}) int {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fail(err)
	}
	fmt.Println(string(data))
	return 0
}

func flagSet(name string) *flag.FlagSet {
	return flag.NewFlagSet(name, flag.ContinueOnError)
}

func printText(text string) int {
	fmt.Print(text)
	return 0
}

func printSummary(value interface{}) int {
	return printJSON(value)
}
