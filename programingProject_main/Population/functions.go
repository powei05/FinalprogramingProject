package main

import (
	"math"
	"math/rand"
)

func UpdateAcceleration(ecosystem *Ecosystem, i int) OrderedPair {
	// 1. Calculate the separation force to avoid crowding.
	forceX, forceY := CalculateSeparationForce(ecosystem, i)

	// 2. Create a persistent, slowly wandering propulsion force.
	// a. Retrieve the family's current propulsion direction.
	propulsionDir := ecosystem.Families[i].PropulsionDirection

	// b. If the family has nearly stopped, force it to pick a new random direction.
	currentSpeedMag := math.Hypot(ecosystem.Families[i].MovementSpeed.x, ecosystem.Families[i].MovementSpeed.y)
	if currentSpeedMag < 0.1 { // Threshold for being "stuck"
		newAngle := rand.Float64() * 2 * math.Pi
		propulsionDir = OrderedPair{x: math.Cos(newAngle), y: math.Sin(newAngle)}
	} else {
		// c. If moving, apply a small, random turn to the propulsion direction to make it wander smoothly.
		turnStrength := 0.3 // How sharply the propulsion direction can change per step.
		angleChange := (rand.Float64()*2 - 1) * turnStrength
		cos := math.Cos(angleChange)
		sin := math.Sin(angleChange)
		newPropulsionX := propulsionDir.x*cos - propulsionDir.y*sin
		newPropulsionY := propulsionDir.x*sin + propulsionDir.y*cos
		propulsionDir = OrderedPair{x: newPropulsionX, y: newPropulsionY}
	}

	ecosystem.Families[i].PropulsionDirection = propulsionDir

	// c. Calculate the final propulsion force based on the new direction.
	propulsionStrength := 3.0 // The magnitude of the "gas pedal".
	propulsionX := propulsionDir.x * propulsionStrength
	propulsionY := propulsionDir.y * propulsionStrength

	// 3. The final acceleration is the sum of the propulsion force and the separation force.
	// CRITICAL FIX: For neutral species like humans who may not have other families to interact with,
	// we need to ensure their propulsion force is strong enough to guarantee movement.
	if ecosystem.Families[i].species.Type == "neutral" {
		// Give neutral species a stronger base propulsion to ensure they always move.
		return OrderedPair{x: propulsionX * 2.0, y: propulsionY * 2.0}
	}

	separationWeight := 2.0
	return OrderedPair{x: propulsionX + forceX*separationWeight, y: propulsionY + forceY*separationWeight}
}

// UpdatePropulsionDirection calculates the new propulsion direction for the next frame.
func UpdatePropulsionDirection(f Family) OrderedPair {
	propulsionDir := f.PropulsionDirection

	// If the family has nearly stopped, force it to pick a new random direction.
	currentSpeedMag := math.Hypot(f.MovementSpeed.x, f.MovementSpeed.y)
	if currentSpeedMag < 0.1 { // Threshold for being "stuck"
		newAngle := rand.Float64() * 2 * math.Pi
		propulsionDir = OrderedPair{x: math.Cos(newAngle), y: math.Sin(newAngle)}
	} else {
		// If moving, apply a small, random turn to the propulsion direction to make it wander smoothly.
		turnStrength := 0.3 // How sharply the propulsion direction can change per step.
		angleChange := (rand.Float64()*2 - 1) * turnStrength
		cos := math.Cos(angleChange)
		sin := math.Sin(angleChange)
		newPropulsionX := propulsionDir.x*cos - propulsionDir.y*sin
		newPropulsionY := propulsionDir.x*sin + propulsionDir.y*cos
		propulsionDir = OrderedPair{x: newPropulsionX, y: newPropulsionY}
	}
	return propulsionDir
}

// CalculateSeparationForce computes the separation force exerted on a family by all its neighbors.
// The strength of the force depends on the family's own species type.
func CalculateSeparationForce(ecosystem *Ecosystem, i int) (float64, float64) {
	// Fx = C_separation * (x1-x2) / d^2
	// Fy = C_separation * (y1-y2) / d^2

	var separationCoefficient float64
	familyType := ecosystem.Families[i].species.Type

	// Determine the separation coefficient based on the current family's type.
	switch familyType {
	case "predator":
		separationCoefficient = 1.5 // CsepPred
	case "prey":
		separationCoefficient = 1.0 // CsepPrey
	case "neutral":
		separationCoefficient = 1.0 // CsepNeutral
	default:
		separationCoefficient = 1.0
	}

	dthreshold := Separation_Threshold // proximity threshold
	currentFamily := ecosystem.Families[i]
	x1 := currentFamily.Position.x
	y1 := currentFamily.Position.y

	SepSumX := 0.0
	SepSumY := 0.0
	n := 0.0

	// Iterate over all other families to calculate the total force
	for j, otherFamily := range ecosystem.Families {
		if i == j {
			continue
		}
		x2 := otherFamily.Position.x
		y2 := otherFamily.Position.y
		dx := x1 - x2
		dy := y1 - y2
		d := math.Sqrt(dx*dx + dy*dy)
		if d < dthreshold && d > 0 { // d > 0 to avoid division by zero
			forceMagnitude := separationCoefficient / (d * d)
			SepSumX += dx * forceMagnitude
			SepSumY += dy * forceMagnitude
			n++
		}
	}

	// Average the force over the number of influencing neighbors.
	// This is a common technique in Boids-like simulations to prevent extreme forces.
	if n > 0 {
		return SepSumX / n, SepSumY / n
	}
	return 0.0, 0.0
}

func UpdateVelocity(f Family, oldAcceleration OrderedPair, newAcceleration OrderedPair, maxFamilySpeed, timeStep float64, weather string) OrderedPair {
	//vx(n+1)=(1/2)(ax(n)+ax(n+1))*t+vx(n)
	//vy(n+1)=(1/2)(ay(n)+ay(n+1))*t+vy(n)
	oldAx := oldAcceleration.x
	oldAy := oldAcceleration.y
	oldVx := f.MovementSpeed.x
	oldVy := f.MovementSpeed.y
	newAx := newAcceleration.x
	newAy := newAcceleration.y
	vx := (0.5)*(oldAx+newAx)*timeStep + oldVx
	vy := (0.5)*(oldAy+newAy)*timeStep + oldVy
	s := math.Sqrt(vx*vx + vy*vy)
	if s > maxFamilySpeed {
		vx = vx * (maxFamilySpeed / s)
		vy = vy * (maxFamilySpeed / s)
	}

	// Apply weather effect to speed
	weatherCoeff := 1.0 + CoefficientOfMovingSpeedIncrease(weather)

	return OrderedPair{x: vx * weatherCoeff, y: vy * weatherCoeff}

}

func UpdatePosition(f Family, oldAcceleration OrderedPair, oldVelocity OrderedPair, ecosystemWidth, timeStep float64) OrderedPair {
	//px(n+1)=(1/2)ax(n)*t^2+vx(n)*t+px(n)
	//py(n+1)=(1/2)ay(n)*t^2+vy(n)*t+py(n)
	oldAx := oldAcceleration.x
	oldAy := oldAcceleration.y
	oldVx := oldVelocity.x
	oldVy := oldVelocity.y
	Px := (0.5)*((oldAx)*timeStep*timeStep) + (oldVx * timeStep) + f.Position.x
	Py := (0.5)*((oldAy)*timeStep*timeStep) + (oldVy * timeStep) + f.Position.y
	for Px > ecosystemWidth {
		Px = Px - ecosystemWidth
	}
	for Px < 0 {
		Px = ecosystemWidth + Px
	}
	for Py > ecosystemWidth {
		Py = Py - ecosystemWidth
	}
	for Py < 0 {
		Py = ecosystemWidth + Py
	}

	return OrderedPair{x: Px, y: Py}
}

func updateFamilyPopulations(eco *Ecosystem) {
	growthRates := make([]float64, len(eco.Families))
	currentCounts := CountSpecies(eco)

	// Step 1: Calculate base growth rates (natural growth/decline, weather, carrying capacity)
	for i := range eco.Families {
		f := eco.Families[i]
		gr := f.species.GrowthRate * (1.0 + CoefficientOfAnimalGrowthRateIncrease(eco.weather))

		if capacity, ok := eco.CarryingCapacity[f.species.Name]; ok && capacity > 0 { // Assumes CarryingCapacity is part of Ecosystem struct
			gr *= (1.0 - float64(currentCounts[f.species.Name])/float64(capacity))
		}

		if f.species.Type == "prey" {
			gr += PlantCoefficient // This could be more complex, e.g., based on local plant density
		}

		// Add a growth bonus if the family is inside the lake.
		if IsInLake(f.Position, eco.Lake) {
			growthBonus := 0.1 // Define the bonus for being in the lake.
			gr += growthBonus
		}
		growthRates[i] = gr
	}

	// Step 2: Add growth rates from pairwise interactions
	for i := 0; i < len(eco.Families); i++ {
		for j := i + 1; j < len(eco.Families); j++ {
			contactGR_A, contactGR_B := Check(eco.Families[i], eco.Families[j])
			// Modify growth rate based on interaction, scaled by the other family's size
			// This makes predation more impactful.
			// A (i) interacting with B (j). We scale the growth rate by a factor of the other family's size.
			// Adding a small constant prevents the effect from being zero if a family size is 1.
			growthRates[i] += contactGR_A * (1.0 + float64(eco.Families[j].Size))
			growthRates[j] += contactGR_B * (1.0 + float64(eco.Families[i].Size))
		}
	}

	// Step 3: Apply the final calculated growth rates to update family sizes
	for i := range eco.Families {
		size := float64(eco.Families[i].Size)
		// Change from multiplicative to additive model for continuous time simulation.
		// The change in size is the current size multiplied by the growth rate.
		newSize := size + size*growthRates[i]
		eco.Families[i].Size = int(math.Round(newSize))
	}

	// Step 3.5: Enforce hard carrying capacity limits.
	// This ensures the total population of a species never exceeds its defined limit.
	finalCounts := CountSpecies(eco)
	for speciesName, totalCount := range finalCounts {
		if capacity, ok := eco.CarryingCapacity[speciesName]; ok && totalCount > capacity {
			// The population has exceeded the carrying capacity.
			excess := totalCount - capacity

			// Reduce the population proportionally from each family of that species.
			for i := range eco.Families {
				if eco.Families[i].species.Name == speciesName {
					// Calculate how much this family should contribute to the reduction.
					reduction := int(math.Round(float64(excess) * (float64(eco.Families[i].Size) / float64(totalCount))))
					if reduction > 0 {
						eco.Families[i].Size -= reduction
						// Ensure size doesn't drop below zero from rounding errors.
						if eco.Families[i].Size < 0 {
							eco.Families[i].Size = 0
						}
					}
				}
			}
		}
	}

	// Remove extinct families
	compacted := eco.Families[:0]
	for _, f := range eco.Families {
		if f.Size > 0 {
			compacted = append(compacted, f)
		}
	}
	eco.Families = compacted
}

func UpdateEcosystem(ecosystem *Ecosystem, timeStep float64) {
	// 1. Update Weather periodically.
	ecosystem.weatherChangeCounter++
	if ecosystem.weatherChangeCounter >= Weather_Change_Interval {
		ecosystem.UpdateWeather()
		ecosystem.weatherChangeCounter = 0 // Reset the counter
	}

	// First, update family movement and physics
	updatedFamilies := make([]Family, len(ecosystem.Families))

	for i, f := range ecosystem.Families {
		oldAcceleration := f.Acceleration // Use the stored acceleration from the previous step
		// The acceleration calculation now only reads state, it doesn't change it.
		newAcceleration := UpdateAcceleration(ecosystem, i)
		newVelocity := UpdateVelocity(f, oldAcceleration, newAcceleration, Max_Family_Speed, timeStep, ecosystem.weather)
		newPosition := UpdatePosition(f, newAcceleration, newVelocity, ecosystem.width, timeStep) // Calculate potential new position

		// Check if the next position is inside the lake. If so, treat it as a collision.
		if IsInLake(newPosition, ecosystem.Lake) {
			// Reverse the velocity to "bounce" off the lake.
			newVelocity.x *= -1
			newVelocity.y *= -1
			// Recalculate the position based on the bounced velocity to prevent entering the lake.
			newPosition = UpdatePosition(f, newAcceleration, newVelocity, ecosystem.width, timeStep)
		}

		// Decide the *next* frame's propulsion direction based on the *current* state.
		nextPropulsionDirection := UpdatePropulsionDirection(f)

		updatedFamilies[i] = Family{
			Size:                f.Size,
			MovementSpeed:       newVelocity,
			Position:            newPosition,
			MovementDirection:   newVelocity,
			Acceleration:        newAcceleration,         // Store the new acceleration for the next step
			PropulsionDirection: nextPropulsionDirection, // Store the NEWLY decided direction for the next frame.
			species:             f.species,
		}
	}
	ecosystem.Families = updatedFamilies

	// 3. Update Plants (Growth and Consumption)
	// Plant growth with weather effect
	plantGrowthCoeff := 1.0 + CoefficientOfPlantIncrease(ecosystem.weather)
	ecosystem.Plants = PlantGrowth(ecosystem.Plants, plantGrowthCoeff)

	// Prey consume plants
	ConsumePlants(ecosystem, consumptionRate, Eating_Threshold)

	// 4. Update Animal Populations based on interactions and environment
	updateFamilyPopulations(ecosystem)

	// 5. Split large families
	SplitLargeFamilies(ecosystem)

	// 6. Merge small families
	MergeFamilies(ecosystem)
}

func SimulateEcosystem(initialEcosystem Ecosystem, numGens int, timeStep float64) []Ecosystem {
	updatedEcosystem := make([]Ecosystem, numGens)
	updatedEcosystem[0] = initialEcosystem
	for i := 1; i < len(updatedEcosystem); i++ {
		nextState := updatedEcosystem[i-1]
		UpdateEcosystem(&nextState, timeStep)
		updatedEcosystem[i] = nextState
	}
	return updatedEcosystem
}

// CountSpecies returns a map of species name to total population count for a given Ecosystem snapshot.
func CountSpecies(ecosystem *Ecosystem) map[string]int {
	counts := make(map[string]int)
	for _, family := range ecosystem.Families {
		counts[family.species.Name] += family.Size
	}
	return counts
}

func InitFamilies(speciesName string, totalPopulation int, initialSpeed, ecosystemWidth float64) []Family {
	numFamilies := initial_family_number
	if totalPopulation < numFamilies {
		numFamilies = totalPopulation
	}

	families := make([]Family, numFamilies)
	s := SpeciesRegistry[speciesName]

	// Distribute population among families
	baseSize := totalPopulation / numFamilies
	remainder := totalPopulation % numFamilies

	for i := 0; i < numFamilies; i++ {
		familySize := baseSize
		if i < remainder {
			familySize++
		}

		x := rand.Float64() * ecosystemWidth
		y := rand.Float64() * ecosystemWidth
		angle := rand.Float64() * 2 * math.Pi
		vx := initialSpeed * math.Cos(angle)
		vy := initialSpeed * math.Sin(angle)

		families[i] = Family{
			Size:              familySize,
			MovementSpeed:     OrderedPair{x: vx, y: vy},
			Position:          OrderedPair{x: x, y: y},
			MovementDirection: OrderedPair{x: vx, y: vy},
			Acceleration:      OrderedPair{x: 0, y: 0}, // Initialize acceleration to zero
			species:           s,
		}
	}
	return families
}

func ConsumePlants(ecosystem *Ecosystem, consumptionRate float64, threshold float64) {
	for fi := range ecosystem.Families {
		f := &ecosystem.Families[fi]
		if f.species.Type == "prey" { // only prey eat plants
			for pi := range ecosystem.Plants {
				dx := ecosystem.Plants[pi].position.x - f.Position.x
				dy := ecosystem.Plants[pi].position.y - f.Position.y
				d := math.Sqrt(dx*dx + dy*dy)
				if d < threshold {
					ecosystem.Plants[pi].size -= consumptionRate
					if ecosystem.Plants[pi].size < 0 {
						ecosystem.Plants[pi].size = 0
					}
				}
			}
		}
	}
}

func PlantGrowth(plants []Plant, growthCoeff float64) []Plant {
	for i := range plants {
		if plants[i].size > 0 {
			plants[i].size += PlantCoefficient * plants[i].size * growthCoeff
		}
	}
	return plants
}

// --- Lake Functions ---

func InitializeLake(x, y, radius float64) Lake {
	return Lake{
		Position: OrderedPair{x: x, y: y},
		Radius:   radius,
	}
}

// check if a given position is within the lake
func IsInLake(position OrderedPair, lake Lake) bool {
	// Calculate the distance from the position to the center of the lake.
	dist := math.Hypot(position.x-lake.Position.x, position.y-lake.Position.y)
	// If the distance is less than the radius, the position is inside the lake.
	return dist <= lake.Radius
}
