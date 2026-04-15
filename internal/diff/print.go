package diff

import (
	"fmt"
	"sort"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

// Print writes a human-readable diff report to stdout.
func Print(entries []DriftEntry, labelA, labelB string) {
	if len(entries) == 0 {
		fmt.Printf("%s✔ No drift detected between %s and %s%s\n",
			colorGreen, labelA, labelB, colorReset)
		return
	}

	// Sort entries for deterministic output.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	fmt.Printf("%s%sDrift report: %s ↔ %s%s\n",
		colorBold, colorYellow, labelA, labelB, colorReset)
	fmt.Printf("%d difference(s) found:\n\n", len(entries))

	for _, e := range entries {
		switch e.Kind {
		case Added:
			fmt.Printf("%s+ %-40s = %s%s\n", colorGreen, e.Key, e.ValueB, colorReset)
		case Removed:
			fmt.Printf("%s- %-40s = %s%s\n", colorRed, e.Key, e.ValueA, colorReset)
		case Changed:
			fmt.Printf("%s~ %-40s%s\n", colorYellow, e.Key, colorReset)
			fmt.Printf("  %s- %s%s\n", colorRed, e.ValueA, colorReset)
			fmt.Printf("  %s+ %s%s\n", colorGreen, e.ValueB, colorReset)
		}
	}
}
