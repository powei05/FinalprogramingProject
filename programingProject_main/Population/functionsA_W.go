package main

import (
	"math"
	"math/rand"
	"sort"
)

// function to initialize an ecosystem with initial populations and families.
func InitializeEcosystem() Ecosystem {
	width := Ecosystem_Width
	var families []Family

	// CRITICAL FIX: Initialize the lake *before* creating families that need to check its position.
	lake := InitializeLake(250, 250, 75) // Center (250,250), Radius 75

	// Initialize animal families from the initialPopulations map
	for speciesName, totalPopulation := range initialPopulations { // Now includes humans
		speciesData := SpeciesRegistry[speciesName]
		familySizes := randomPartition(totalPopulation, initial_family_number, Smallest_Family_Size)
		for _, size := range familySizes {
			// Increase initial speed to make movement more visible from the start.
			initialSpeedMagnitude := 10.0 // You can adjust this value

			var pos OrderedPair
			// Ensure families do not spawn inside the lake.
			for {
				pos = OrderedPair{x: rand.Float64() * width, y: rand.Float64() * width}
				if !IsInLake(pos, lake) {
					break // Found a valid position, exit the loop.
				}
			}

			angle := rand.Float64() * 2 * math.Pi // Generate a random direction
			speed := OrderedPair{x: initialSpeedMagnitude * math.Cos(angle), y: initialSpeedMagnitude * math.Sin(angle)}

			// Initialize the propulsion direction to a random unit vector
			propulsionAngle := rand.Float64() * 2 * math.Pi
			propulsionDir := OrderedPair{x: math.Cos(propulsionAngle), y: math.Sin(propulsionAngle)}

			families = append(families, Family{
				Size:                size,
				MovementSpeed:       speed,
				Position:            pos,
				MovementDirection:   speed,                   // Initially, direction is same as speed
				Acceleration:        OrderedPair{x: 0, y: 0}, // Initialize acceleration to zero
				PropulsionDirection: propulsionDir,
				species:             speciesData,
			})
		}
	}

	// Initialize plants
	var plants []Plant
	numPlants := 200 // Let's start with 200 plants
	for i := 0; i < numPlants; i++ {
		pos := OrderedPair{x: rand.Float64() * width, y: rand.Float64() * width}
		// Ensure plants do not spawn inside the lake.
		if !IsInLake(pos, lake) {
			plants = append(plants, Plant{position: pos, size: rand.Float64()*10 + 5}) // Random initial size
		}
	}

	return Ecosystem{
		Families: families,
		Lake:     lake,
		Plants:   plants, // Add the initialized plants
		width:    width,
		CarryingCapacity: map[string]int{
			"rabbit": 1200,
			"sheep":  800,
			"deer":   500,
			"wolf":   150,
		},
	}
}

// Help function to initialize family sizes randomly
func randomPartition(total, k, min int) []int {
	if k <= 0 || total < k*min {
		// If the total population is too small to partition,
		// put all individuals into a single family, provided the total is not zero.
		if total > 0 {
			return []int{total}
		}
		return []int{} // Return empty if total is zero.
	}

	// 1. the number to be randomly partitioned after allocating min to each part
	remain := total - k*min

	// 2. Generate k-1 random cut points between 0 and remain
	cuts := make([]int, k-1)
	for i := range cuts {
		cuts[i] = rand.Intn(remain + 1)
	}

	// 3. Sort the cut points and add 0 and remain as the boundaries
	sort.Ints(cuts)
	allCuts := append([]int{0}, cuts...)
	allCuts = append(allCuts, remain)

	// 4. Calculate the length of each segment (difference between adjacent cut points) and add back min.
	parts := make([]int, k)
	for i := 0; i < k; i++ {
		segmentLength := allCuts[i+1] - allCuts[i]
		parts[i] = segmentLength + min
	}
	return parts
}

// function to update populations based on growth rates
func UpdatePopulations(ecosystem *Ecosystem) {
	if len(ecosystem.Families) == 0 {
		return
	}

	growthRates := make([]float64, len(ecosystem.Families))

	// Step 1: Calculate base growth rates (natural growth, plant consumption)
	for i, f := range ecosystem.Families {
		gr := f.species.GrowthRate
		if ecosystem.Families[i].species.Type == "prey" {
			gr += PlantCoefficient
		}
		growthRates[i] = gr
	}

	// Step 2: Add growth rates from pairwise interactions (predation)
	for i := 0; i < len(ecosystem.Families); i++ {
		for j := i + 1; j < len(ecosystem.Families); j++ {
			contactGR_A, contactGR_B := Check(ecosystem.Families[i], ecosystem.Families[j])
			growthRates[i] += contactGR_A
			growthRates[j] += contactGR_B
		}
	}

	// Step 3: Apply the final calculated growth rates to update family sizes
	for i := range ecosystem.Families {
		size := float64(ecosystem.Families[i].Size)
		newSize := int(math.Round(size * (1.0 + growthRates[i])))
		ecosystem.Families[i].Size = newSize
	}

	// Step 4: Remove extinct families (size <= 0)
	compacted := ecosystem.Families[:0]
	for _, f := range ecosystem.Families {
		if f.Size > 0 {
			compacted = append(compacted, f)
		}
	}
	ecosystem.Families = compacted
}

// function to merge small family with someone nearby
func MergeFamilies(ecosystem *Ecosystem) {
	f := ecosystem.Families
	for i := 0; i < len(f); {
		if f[i].Size < Smallest_Family_Size {
			merged := false
			for j := 0; j < len(f); j++ {
				if i == j || f[i].species.Name != f[j].species.Name {
					continue
				}
				if distance(f[i].Position, f[j].Position) <= Merging_Threshold {
					f[j].Size += f[i].Size
					f[i] = f[len(f)-1]
					f = f[:len(f)-1]
					merged = true
					break
				}
			}
			if merged {
				continue
			}
		}
		i++
	}
	ecosystem.Families = f
}

// SplitLargeFamilies checks for families that have grown larger than Max_Family_Size and splits them.
func SplitLargeFamilies(ecosystem *Ecosystem) {
	// We must build a new slice because we are both modifying existing families (size) and adding new ones.
	var nextGenerationFamilies []Family

	for i := range ecosystem.Families {
		f := &ecosystem.Families[i] // Use a pointer to modify the original family

		if f.Size > Max_Family_Size {
			// This family needs to be split.
			originalNewSize := f.Size / 2
			splitNewSize := f.Size - originalNewSize

			// Update the original family's size.
			f.Size = originalNewSize

			// --- Give both families a large, opposing VELOCITY to push them apart ---
			// This is more effective than acceleration as it's an immediate change in speed,
			// and won't be overwritten by the next frame's acceleration calculation.
			splitSpeedBoost := 30.0 // A large speed boost.
			angle := rand.Float64() * 2 * math.Pi
			pushVx := splitSpeedBoost * math.Cos(angle)
			pushVy := splitSpeedBoost * math.Sin(angle)

			// Add the push velocity to the parent family's current speed.
			f.MovementSpeed.x += pushVx
			f.MovementSpeed.y += pushVy

			// Create the new family, inheriting properties from the parent.
			newFamily := Family{
				Size: splitNewSize,
				// The new family gets pushed in the opposite direction.
				MovementSpeed:     OrderedPair{x: f.MovementSpeed.x - 2*pushVx, y: f.MovementSpeed.y - 2*pushVy},
				Position:          OrderedPair{x: f.Position.x + (rand.Float64()*2 - 1), y: f.Position.y + (rand.Float64()*2 - 1)}, // Slight offset
				MovementDirection: f.MovementDirection,
				Acceleration:      f.Acceleration, // Inherit acceleration
				species:           f.species,
			}
			nextGenerationFamilies = append(nextGenerationFamilies, newFamily)
		}
	}
	// Add the newly created families from splits to the main ecosystem slice.
	ecosystem.Families = append(ecosystem.Families, nextGenerationFamilies...)
}

func distance(a, b OrderedPair) float64 {
	return math.Hypot(a.x-b.x, a.y-b.y)
}

func Check(A, B Family) (float64, float64) {
	if distance(A.Position, B.Position) < Eating_Threshold {
		// Case 1: A is predator, B is prey
		if A.species.Type == "predator" && B.species.Type == "prey" {
			return A.species.ContactGrowthRate, B.species.ContactGrowthRate
		}
		// Case 2: B is predator, A is prey
		if B.species.Type == "predator" && A.species.Type == "prey" {
			return A.species.ContactGrowthRate, B.species.ContactGrowthRate
		}
	}

	return 0.0, 0.0
}
