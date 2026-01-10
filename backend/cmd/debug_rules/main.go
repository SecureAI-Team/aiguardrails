package main

import (
	"encoding/json"
	"fmt"

	"aiguardrails/internal/rules"
)

func main() {
	store := rules.NewMemoryStore()

	// Test loading vendor_siemens.json
	path := "policies/vendor_siemens.json"
	fmt.Printf("Loading %s...\n", path)
	if err := rules.LoadFromJSON(path, store); err != nil {
		fmt.Printf("Error loading vendor_siemens: %v\n", err)
	} else {
		list, _ := store.List()
		fmt.Printf("Loaded vendor_siemens rules: %d\n", len(list))
	}

	// Test loading seed_rules.json
	path2 := "policies/seed_rules.json"
	fmt.Printf("Loading %s...\n", path2)
	if err := rules.LoadFromJSON(path2, store); err != nil {
		fmt.Printf("Error loading seed_rules: %v\n", err)
	} else {
		list, _ := store.List()
		fmt.Printf("Total rules in store: %d\n", len(list))
		if len(list) > 0 {
			b, _ := json.MarshalIndent(list[0], "", "  ")
			fmt.Printf("First rule: %s\n", string(b))
		}
	}
}
