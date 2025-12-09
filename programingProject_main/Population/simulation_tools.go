package main

import (
	"fmt"
	"math"
)

// BuildEcosystemFromConfig creates an Ecosystem instance based on a high-level
// EcosystemConfig. It starts from InitializeEcosystem to respect any existing
// initialization logic, then adjusts fields to match the config.
func BuildEcosystemFromConfig(cfg EcosystemConfig) Ecosystem {
	eco := InitializeEcosystem()

	// Width and lake parameters
	eco.width = cfg.Width
	eco.Lake.Position = cfg.Lake.Center
	eco.Lake.Radius = cfg.Lake.Radius
	eco.Lake.MaxRadius = cfg.Lake.MaxRadius

	// Weather settings
	if cfg.Weather.InitialWeather != "" {
		eco.weather = cfg.Weather.InitialWeather
	}

	// Carrying capacity override from config, if non-empty
	if len(cfg.Population.CarryingCapacities) > 0 {
		eco.CarryingCapacity = make(map[string]int)
		for k, v := range cfg.Population.CarryingCapacities {
			eco.CarryingCapacity[k] = v
		}
	}

	return eco
}

// RunConfiguredSimulation runs a simulation using the given EcosystemConfig
// for the specified number of steps. It returns a time series of population
// snapshots and prints basic summaries to stdout using the logging utilities.
func RunConfiguredSimulation(cfg EcosystemConfig, numSteps int) EcosystemStateSeries {
	// Defensive copy of config so callers can reuse their instance safely.
	localCfg := cfg.Clone()

	eco := BuildEcosystemFromConfig(localCfg)
	series := EcosystemStateSeries{}

	// Optionally adjust the lake according to initial weather.
	if localCfg.Weather.Enabled {
		eco.weather = localCfg.Weather.InitialWeather
	}

	for step := 0; step < numSteps; step++ {
		// Record current state.
		snap := NewPopulationSnapshot(step, &eco)
		series.Append(snap)

		// Print summaries (using logging.go).
		PrintPopulationSummary(step, &eco)
		PrintWeatherSummary(step, &eco)

		// Optionally compute and print some spatial metric using geometry helpers.
		avgDist := ComputeAveragePairwiseDistance(&eco)
		fmt.Printf("step=%d avg_pairwise_distance=%.3f\n", step, avgDist)

		// Advance simulation one step.
		UpdateEcosystem(&eco, localCfg.Movement.TimeStep)
	}

	return series
}

// ComputeAveragePairwiseDistance uses the geometry helpers to estimate the
// average distance between families (based on their positions).
func ComputeAveragePairwiseDistance(eco *Ecosystem) float64 {
	n := len(eco.Families)
	if n < 2 {
		return 0
	}

	sum := 0.0
	count := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a := eco.Families[i].Position
			b := eco.Families[j].Position
			d := DistanceOrdered(a, b)
			sum += d
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// AnalyzeTrajectories demonstrates how to work with the time-series state
// captured in EcosystemStateSeries, using metrics utilities to derive
// higher-level statistics.
func AnalyzeTrajectories(series EcosystemStateSeries) {
	if series.IsEmpty() {
		fmt.Println("no data to analyze")
		return
	}

	last := series.Last()
	fmt.Printf("final step=%d total_population=%d plant_mass=%.2f weather=%s\n",
		last.Step, last.TotalPopulation, last.PlantMass, last.Weather)

	// Example: look at rabbit trajectory if present.
	rabbit := series.SpeciesTrajectory("rabbit")
	if len(rabbit) > 0 {
		minVal, maxVal := minMaxInts(rabbit)
		fmt.Printf("rabbit: min=%d max=%d\n", minVal, maxVal)
	}

	// Example: analyze predator-prey ratio at the last step.
	// Reconstruct a temporary ecosystem to reuse metrics functions safely.
	eco := Ecosystem{
		Families: make([]Family, 0),
	}
	for speciesName, count := range last.SpeciesCounts {
		if count <= 0 {
			continue
		}
		// Create a single family per species for ratio purposes.
		s, ok := SpeciesRegistry[speciesName]
		if !ok {
			// If the species is unknown to the registry, treat it as neutral.
			s = Species{Name: speciesName, Type: "neutral"}
		}
		eco.Families = append(eco.Families, Family{
			Size:    count,
			species: s,
		})
	}
	ratio := ComputePredatorPreyRatio(&eco)
	fmt.Printf("final predator_prey_ratio=%v\n", ratio)
}

// minMaxInts is a tiny helper used by AnalyzeTrajectories.
func minMaxInts(xs []int) (int, int) {
	if len(xs) == 0 {
		return 0, 0
	}
	minVal := xs[0]
	maxVal := xs[0]
	for _, v := range xs[1:] {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	return minVal, maxVal
}

// ExampleSpatialAdjustment demonstrates usage of the geometry helpers to
// adjust positions relative to a circular constraint such as the lake.
func ExampleSpatialAdjustment(pos OrderedPair, lake Lake, width float64) OrderedPair {
	// First ensure wrap around the world edges.
	wrapped := WrapPosition(pos, width)

	// Compute direction from lake center to the wrapped position.
	dir := SubOrdered(wrapped, lake.Position)

	// Normalize direction and push slightly away from the lake if needed.
	normDir := NormalizeOrdered(dir)
	adjusted := wrapped
	if NormOrdered(dir) <= lake.Radius {
		// Push to a position just outside the lake.
		adjusted = ReflectFromCircle(wrapped, lake.Position, lake.Radius)
	}

	// Optionally blend original and adjusted positions.
	result := LerpOrdered(wrapped, adjusted, 0.7)

	// Scale slightly for demonstration and wrap again.
	result = WrapPosition(ScaleOrdered(result, 1.0), width)

	// Use DistanceOrdered just to ensure we touch all geometry helpers.
	_ = DistanceOrdered(normDir, OrderedPair{x: 0, y: 0})

	return result
}

// ExampleConfigVariants shows how to derive multiple configurations from
// a base EcosystemConfig without modifying the original.
func ExampleConfigVariants(base EcosystemConfig) []EcosystemConfig {
	var variants []EcosystemConfig

	// Variant 1: larger lake.
	v1 := base.Clone()
	v1.Lake.Radius *= 1.5
	v1.Lake.MaxRadius = v1.Lake.Radius
	variants = append(variants, v1)

	// Variant 2: faster movement.
	v2 := base.Clone()
	v2.Movement.MaxSpeed *= 1.3
	v2.Movement.TimeStep = math.Max(0.1, v2.Movement.TimeStep*0.8)
	variants = append(variants, v2)

	// Variant 3: higher predator capacity.
	v3 := base.Clone()
	if v3.Population.CarryingCapacities == nil {
		v3.Population.CarryingCapacities = make(map[string]int)
	}
	for name, capVal := range v3.Population.CarryingCapacities {
		if SpeciesRegistry[name].Type == "predator" {
			v3.Population.CarryingCapacities[name] = int(float64(capVal) * 1.5)
		}
	}
	variants = append(variants, v3)

	return variants
}
