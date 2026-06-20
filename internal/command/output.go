package command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CommandResult struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Title     string `json:"title,omitempty"`
	Updated   bool   `json:"updated,omitempty"`
	Structure string `json:"structure,omitempty"`
	Card      string `json:"card,omitempty"`
	Relation  string `json:"relation,omitempty"`
	Warning   string `json:"warning,omitempty"`
}

func isJSONOutput(cmd *cobra.Command) bool {
	return viper.GetString("output") == "json"
}

func printResult(cmd *cobra.Command, out io.Writer, result CommandResult) {
	if isJSONOutput(cmd) {
		data, _ := json.Marshal(result)
		fmt.Fprintln(out, string(data))
	} else {
		if result.ID != "" {
			fmt.Fprintf(out, "✓ %s %s\n", result.Type, result.ID)
			if result.Title != "" {
				fmt.Fprintf(out, "  Title: %s\n", result.Title)
			}
		} else if result.Updated {
			fmt.Fprintf(out, "✓ Updated card %s\n", result.ID)
		} else if result.Structure != "" {
			fmt.Fprintf(out, "✓ Added %s to %s\n", result.Card, result.Structure)
			fmt.Fprintf(out, "  relation: %s\n", result.Relation)
			if result.Warning != "" {
				fmt.Fprintf(out, "  warning: %s\n", result.Warning)
			}
		}
	}
}