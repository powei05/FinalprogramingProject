package main

import "math"

type PopulationSummary struct {
	TotalPopulation int
	SpeciesCounts   map[string]int
	AverageFamily   float64
	FamilyCount     int
	PlantMass       float64
	DiversityIndex  float64
}

func ComputeTotalPopulation(ecosystem *Ecosystem) int {
	sum := 0
	for _, f := range ecosystem.Families {
		sum += f.Size
	}
	return sum
}

func ComputeSpeciesCounts(ecosystem *Ecosystem) map[string]int {
	return CountSpecies(ecosystem)
}

func ComputeAverageFamilySize(ecosystem *Ecosystem) float64 {
	if len(ecosystem.Families) == 0 {
		return 0
	}
	sum := ComputeTotalPopulation(ecosystem)
	return float64(sum) / float64(len(ecosystem.Families))
}

func ComputeDiversityIndex(counts map[string]int) float64 {
	total := 0
	for _, v := range counts {
		total += v
	}
	if total == 0 {
		return 0
	}
	h := 0.0
	for _, v := range counts {
		if v == 0 {
			continue
		}
		p := float64(v) / float64(total)
		h -= p * math.Log(p)
	}
	return h
}

func BuildPopulationSummary(ecosystem *Ecosystem) PopulationSummary {
	counts := ComputeSpeciesCounts(ecosystem)
	total := ComputeTotalPopulation(ecosystem)
	avg := ComputeAverageFamilySize(ecosystem)
	plantMass := CountPlantMass(ecosystem)
	div := ComputeDiversityIndex(counts)
	return PopulationSummary{
		TotalPopulation: total,
		SpeciesCounts:   counts,
		AverageFamily:   avg,
		FamilyCount:     len(ecosystem.Families),
		PlantMass:       plantMass,
		DiversityIndex:  div,
	}
}

func ComputePredatorPreyRatio(ecosystem *Ecosystem) float64 {
	pred := 0
	prey := 0
	for _, f := range ecosystem.Families {
		switch f.species.Type {
		case "predator":
			pred += f.Size
		case "prey":
			prey += f.Size
		}
	}
	if prey == 0 {
		if pred == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return float64(pred) / float64(prey)
}
