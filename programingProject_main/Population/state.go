package main

type PopulationSnapshot struct {
	Step            int
	SpeciesCounts   map[string]int
	TotalPopulation int
	PlantMass       float64
	Weather         string
}

type EcosystemStateSeries struct {
	Snapshots []PopulationSnapshot
}

func NewPopulationSnapshot(step int, ecosystem *Ecosystem) PopulationSnapshot {
	counts := CountSpecies(ecosystem)
	total := ComputeTotalPopulation(ecosystem)
	plants := CountPlantMass(ecosystem)
	return PopulationSnapshot{
		Step:            step,
		SpeciesCounts:   counts,
		TotalPopulation: total,
		PlantMass:       plants,
		Weather:         ecosystem.weather,
	}
}

func (s *EcosystemStateSeries) Append(snapshot PopulationSnapshot) {
	s.Snapshots = append(s.Snapshots, snapshot)
}

func (s *EcosystemStateSeries) Length() int {
	return len(s.Snapshots)
}

func (s *EcosystemStateSeries) IsEmpty() bool {
	return len(s.Snapshots) == 0
}

func (s *EcosystemStateSeries) Last() PopulationSnapshot {
	if len(s.Snapshots) == 0 {
		return PopulationSnapshot{}
	}
	return s.Snapshots[len(s.Snapshots)-1]
}

func (s *EcosystemStateSeries) SpeciesTrajectory(name string) []int {
	values := make([]int, len(s.Snapshots))
	for i, snap := range s.Snapshots {
		values[i] = snap.SpeciesCounts[name]
	}
	return values
}

func (s *EcosystemStateSeries) TotalPopulationTrajectory() []int {
	values := make([]int, len(s.Snapshots))
	for i, snap := range s.Snapshots {
		values[i] = snap.TotalPopulation
	}
	return values
}
