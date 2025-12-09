package main

import (
	"fmt"
	"sort"
	"strings"
)

func FormatPopulationLine(step int, ecosystem *Ecosystem) string {
	summary := BuildPopulationSummary(ecosystem)
	names := make([]string, 0, len(summary.SpeciesCounts))
	for name := range summary.SpeciesCounts {
		names = append(names, name)
	}
	sort.Strings(names)
	var parts []string
	for _, n := range names {
		parts = append(parts, fmt.Sprintf("%s=%d", n, summary.SpeciesCounts[n]))
	}
	return fmt.Sprintf(
		"step=%d total=%d families=%d plants=%.2f species={%s}",
		step,
		summary.TotalPopulation,
		summary.FamilyCount,
		summary.PlantMass,
		strings.Join(parts, ","),
	)
}

func FormatWeatherLine(step int, ecosystem *Ecosystem) string {
	return fmt.Sprintf(
		"step=%d weather=%s lake_radius=%.2f",
		step,
		ecosystem.weather,
		ecosystem.Lake.Radius,
	)
}

func PrintPopulationSummary(step int, ecosystem *Ecosystem) {
	line := FormatPopulationLine(step, ecosystem)
	fmt.Println(line)
}

func PrintWeatherSummary(step int, ecosystem *Ecosystem) {
	line := FormatWeatherLine(step, ecosystem)
	fmt.Println(line)
}
