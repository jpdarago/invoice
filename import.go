package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func importData(path string, structure *Invoice, flags *pflag.FlagSet) error {
	fileText, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %w", path, err)
	}

	if strings.HasSuffix(path, ".json") {
		err = importJson(fileText, structure)
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		err = importYaml(fileText, structure)
	} else {
		return fmt.Errorf("unsupported file type")
	}
	if err != nil {
		return err
	}

	// Explicitly-set flags take precedence over the imported file.
	applyChangedFlags(structure, flags)
	return nil
}

// applyChangedFlags overrides invoice fields with any flag that was explicitly
// set on the command line, so CLI flags win over the imported file. Values are
// read straight from the flag set via typed getters. Flags not handled here
// are either not invoice fields (import, output) or merged separately
// (metadata).
func applyChangedFlags(structure *Invoice, flags *pflag.FlagSet) {
	if flags.Changed("id") {
		structure.Id, _ = flags.GetString("id")
	}
	if flags.Changed("title") {
		structure.Title, _ = flags.GetString("title")
	}
	if flags.Changed("logo") {
		structure.Logo, _ = flags.GetString("logo")
	}
	if flags.Changed("from") {
		structure.From, _ = flags.GetString("from")
	}
	if flags.Changed("to") {
		structure.To, _ = flags.GetString("to")
	}
	if flags.Changed("date") {
		structure.Date, _ = flags.GetString("date")
	}
	if flags.Changed("due") {
		structure.Due, _ = flags.GetString("due")
	}
	if flags.Changed("currency") {
		structure.Currency, _ = flags.GetString("currency")
	}
	if flags.Changed("note") {
		structure.Note, _ = flags.GetString("note")
	}
	if flags.Changed("tax") {
		structure.Tax, _ = flags.GetFloat64("tax")
	}
	if flags.Changed("discount") {
		structure.Discount, _ = flags.GetFloat64("discount")
	}
	if flags.Changed("item") {
		structure.Items, _ = flags.GetStringSlice("item")
	}
	if flags.Changed("quantity") {
		structure.Quantities, _ = flags.GetIntSlice("quantity")
	}
	if flags.Changed("rate") {
		structure.Rates, _ = flags.GetFloat64Slice("rate")
	}
}

func importJson(text []byte, structure *Invoice) error {
	if !json.Valid(text) {
		return fmt.Errorf("invalid json syntax")
	}

	err := json.Unmarshal(text, structure)
	if err != nil {
		return fmt.Errorf("json does not match invoice schema: %w", err)
	}

	return nil
}

func importYaml(text []byte, structure *Invoice) error {
	err := yaml.Unmarshal(text, structure)
	if err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}

	return nil
}
