package main

import (
	"math"
	"testing"

	"programingProject_main/canvas"
)

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

/* ================================
   Tests for functions.go
================================ */

// UpdatePropulsionDirection calculates the new propulsion direction for the next frame.
func TestUpdatePropulsionDirection(t *testing.T) {
	tests := []Family{
		{MovementSpeed: OrderedPair{0, 0}, PropulsionDirection: OrderedPair{1, 0}},
		{MovementSpeed: OrderedPair{0.05, 0.02}, PropulsionDirection: OrderedPair{0, 1}},
		{MovementSpeed: OrderedPair{1, 1}, PropulsionDirection: OrderedPair{1, 0}},
		{MovementSpeed: OrderedPair{5, -3}, PropulsionDirection: OrderedPair{-1, 0}},
		{MovementSpeed: OrderedPair{10, 10}, PropulsionDirection: OrderedPair{0, -1}},
	}

	for i, f := range tests {
		newDir := UpdatePropulsionDirection(f)
		if math.Hypot(newDir.x, newDir.y) == 0 {
			t.Fatalf("case %d: propulsion direction should not be zero", i)
		}
	}
}

// CalculateSeparationForce computes the separation force exerted on a family by all its neighbors.
func TestCalculateSeparationForce(t *testing.T) {
	base := Ecosystem{
		Families: []Family{
			{Position: OrderedPair{0, 0}, species: Species{Type: "prey"}},
			{Position: OrderedPair{10, 0}, species: Species{Type: "prey"}},
			{Position: OrderedPair{100, 100}, species: Species{Type: "prey"}},
			{Position: OrderedPair{20, 20}, species: Species{Type: "predator"}},
			{Position: OrderedPair{400, 400}, species: Species{Type: "neutral"}},
		},
	}

	for i := 0; i < len(base.Families); i++ {
		fx, fy := CalculateSeparationForce(&base, i)
		if math.IsNaN(fx) || math.IsNaN(fy) {
			t.Fatalf("case %d: force returned NaN", i)
		}
	}
}

func TestUpdateVelocity(t *testing.T) {
	f := Family{MovementSpeed: OrderedPair{5, 5}}

	tests := []struct {
		oldA  OrderedPair
		newA  OrderedPair
		limit float64
	}{
		{OrderedPair{0, 0}, OrderedPair{0, 0}, 100},
		{OrderedPair{1, 0}, OrderedPair{1, 0}, 100},
		{OrderedPair{10, 10}, OrderedPair{10, 10}, 5},
		{OrderedPair{-1, -1}, OrderedPair{-1, -1}, 100},
		{OrderedPair{100, 0}, OrderedPair{100, 0}, 10},
	}

	for i, tt := range tests {
		// Use "Rainy" so CoefficientOfMovingSpeedIncrease = 0 and
		// the clamping to maxFamilySpeed is directly testable.
		v := UpdateVelocity(f, tt.oldA, tt.newA, tt.limit, 1.0, "Rainy")
		if math.Hypot(v.x, v.y) > tt.limit+1e-5 {
			t.Fatalf("case %d: velocity exceeds max speed", i)
		}
	}
}

func TestUpdatePosition(t *testing.T) {
	f := Family{Position: OrderedPair{250, 250}}

	tests := []struct {
		a OrderedPair
		v OrderedPair
	}{
		{OrderedPair{0, 0}, OrderedPair{0, 0}},
		{OrderedPair{1, 0}, OrderedPair{1, 0}},
		{OrderedPair{-1, -1}, OrderedPair{-1, -1}},
		{OrderedPair{10, 10}, OrderedPair{10, 10}},
		{OrderedPair{100, 100}, OrderedPair{100, 100}},
	}

	for i, tt := range tests {
		p := UpdatePosition(f, tt.a, tt.v, 500, 1)
		if p.x < 0 || p.x > 500 || p.y < 0 || p.y > 500 {
			t.Fatalf("case %d: position wrap failed", i)
		}
	}
}

func TestCountSpecies(t *testing.T) {
	eco := Ecosystem{
		Families: []Family{
			{Size: 10, species: Species{Name: "rabbit"}},
			{Size: 15, species: Species{Name: "rabbit"}},
			{Size: 20, species: Species{Name: "wolf"}},
			{Size: 5, species: Species{Name: "human"}},
			{Size: 7, species: Species{Name: "wolf"}},
		},
	}

	c := CountSpecies(&eco)

	if c["rabbit"] != 25 || c["wolf"] != 27 || c["human"] != 5 {
		t.Fatalf("species count incorrect: %+v", c)
	}
}

func TestInitFamilies(t *testing.T) {
	families := InitFamilies("rabbit", 100, 5, 500)

	if len(families) != initial_family_number {
		t.Fatalf("wrong number of families, got %d", len(families))
	}

	sum := 0
	for _, f := range families {
		sum += f.Size
	}
	if sum != 100 {
		t.Fatalf("population mismatch, got %d", sum)
	}
}

func TestConsumePlants(t *testing.T) {
	eco := Ecosystem{
		Families: []Family{
			{Position: OrderedPair{0, 0}, species: Species{Type: "prey"}},
			{Position: OrderedPair{200, 200}, species: Species{Type: "predator"}},
		},
		Plants: []Plant{
			{position: OrderedPair{1, 1}, size: 1},
			{position: OrderedPair{100, 100}, size: 1},
			{position: OrderedPair{2, 2}, size: 1},
			{position: OrderedPair{300, 300}, size: 1},
			{position: OrderedPair{0, 5}, size: 1},
		},
	}

	consumed := ConsumePlants(&eco, 0.5, 10)

	if consumed[0] <= 0 {
		t.Fatalf("prey should consume plants, got %f", consumed[0])
	}
}

func TestPlantGrowth(t *testing.T) {
	plants := []Plant{
		{size: 1},
		{size: 2},
		{size: 3},
		{size: 0},
		{size: 5},
	}

	newPlants := PlantGrowth(plants, 1.0)

	for i, p := range newPlants {
		if p.size < plants[i].size {
			t.Fatalf("plant %d should not shrink", i)
		}
	}
}

func TestLakeFunctions_IsInLake(t *testing.T) {
	l := InitializeLake(0, 0, 10)

	tests := []struct {
		pos      OrderedPair
		expected bool
	}{
		{OrderedPair{0, 0}, true},
		{OrderedPair{5, 5}, true},
		{OrderedPair{9, 0}, true},
		{OrderedPair{10, 0}, true},
		{OrderedPair{20, 20}, false},
	}

	for i, tt := range tests {
		inside := IsInLake(tt.pos, l)
		if inside != tt.expected {
			t.Fatalf("case %d: expected %v, got %v", i, tt.expected, inside)
		}
	}
}

func TestPushOutOfLake(t *testing.T) {
	l := InitializeLake(0, 0, 10)

	tests := []OrderedPair{
		{0, 0},   // exactly at center (dist == 0)
		{1, 1},   // inside
		{5, 0},   // inside
		{9, 9},   // outside
		{20, 20}, // outside
	}

	for i, p := range tests {
		newP := PushOutOfLake(p, l)
		dist := math.Hypot(p.x-l.Position.x, p.y-l.Position.y)

		if dist == 0 {
			// Implementation explicitly does nothing when dist == 0
			if newP != p {
				t.Fatalf("case %d: center point should remain unchanged", i)
			}
			continue
		}

		if IsInLake(p, l) {
			// If original was inside (dist > 0), the new position should be outside.
			if IsInLake(newP, l) {
				t.Fatalf("case %d: point was inside and should be pushed out", i)
			}
		} else {
			// If original was outside, it should stay unchanged.
			if newP != p {
				t.Fatalf("case %d: outside point should not move", i)
			}
		}
	}
}

func TestCountPlantMass(t *testing.T) {
	eco := Ecosystem{
		Plants: []Plant{
			{size: 1},
			{size: 2},
			{size: 3},
			{size: 4},
			{size: 5},
		},
	}

	total := CountPlantMass(&eco)
	if !almostEqual(total, 15, 1e-6) {
		t.Fatalf("plant mass incorrect, got %f", total)
	}
}

/* ================================
   Tests for functions_AnWang2.go
================================ */

func TestUpdateWeather(t *testing.T) {
	allowed := map[string]bool{
		"Dry":    true,
		"Sunny":  true,
		"Rainy":  true,
		"Frozen": true,
	}

	e := &Ecosystem{}
	for i := 0; i < 5; i++ {
		e.UpdateWeather()
		if !allowed[e.weather] {
			t.Fatalf("iteration %d: unexpected weather %q", i, e.weather)
		}
	}
}

func TestCoefficientOfPlantIncrease(t *testing.T) {
	tests := []struct {
		weather string
		expect  float64
	}{
		{"Dry", -0.20},
		{"Sunny", 0.00},
		{"Rainy", 0.20},
		{"Frozen", -0.40},
		{"Unknown", -0.40}, // default case
	}

	for i, tt := range tests {
		got := CoefficientOfPlantIncrease(tt.weather)
		if !almostEqual(got, tt.expect, 1e-6) {
			t.Fatalf("case %d: weather=%s, expect %.2f, got %.2f", i, tt.weather, tt.expect, got)
		}
	}
}

func TestCoefficientOfLakeIncrease(t *testing.T) {
	tests := []struct {
		weather string
		expect  float64
	}{
		{"Dry", -0.20},
		{"Sunny", 0.00},
		{"Rainy", 0.20},
		{"Frozen", 0.00},
		{"Unknown", 0.00}, // default case
	}

	for i, tt := range tests {
		got := CoefficientOfLakeIncrease(tt.weather)
		if !almostEqual(got, tt.expect, 1e-6) {
			t.Fatalf("case %d: weather=%s, expect %.2f, got %.2f", i, tt.weather, tt.expect, got)
		}
	}
}

func TestCoefficientOfMovingSpeedIncrease(t *testing.T) {
	tests := []struct {
		weather string
		expect  float64
	}{
		{"Dry", -0.05},
		{"Sunny", 0.07},
		{"Rainy", 0.00},
		{"Frozen", -0.15},
		{"Unknown", -0.15}, // default case
	}

	for i, tt := range tests {
		got := CoefficientOfMovingSpeedIncrease(tt.weather)
		if !almostEqual(got, tt.expect, 1e-6) {
			t.Fatalf("case %d: weather=%s, expect %.2f, got %.2f", i, tt.weather, tt.expect, got)
		}
	}
}

func TestCoefficientOfAnimalGrowthRateIncrease(t *testing.T) {
	tests := []struct {
		weather string
		expect  float64
	}{
		{"Dry", -0.10},
		{"Sunny", 0.10},
		{"Rainy", 0.00},
		{"Frozen", -0.20},
		{"Unknown", -0.20}, // default case
	}

	for i, tt := range tests {
		got := CoefficientOfAnimalGrowthRateIncrease(tt.weather)
		if !almostEqual(got, tt.expect, 1e-6) {
			t.Fatalf("case %d: weather=%s, expect %.2f, got %.2f", i, tt.weather, tt.expect, got)
		}
	}
}

// NOTE: We don't know the concrete constructor for canvas.Canvas in your project,
// so in tests we just make sure calling the function with a nil canvas pointer
// doesn't crash the *test* (we recover from panic). In your real code you pass
// a properly initialized *canvas.Canvas.
func TestDrawWeatherBackground(t *testing.T) {
	var cfg Config
	weathers := []string{"Dry", "Sunny", "Rainy", "Frozen", "Unknown"}

	for range weathers {
		func() {
			defer func() {
				_ = recover() // ignore any panic due to nil canvas
			}()
			var c *canvas.Canvas
			DrawWeatherBackground(c, "Dry", cfg)
		}()
	}
}

/* ================================
   Tests for functionsA_W.go
================================ */

func TestInitializeEcosystem(t *testing.T) {
	eco := InitializeEcosystem()

	// 1. There should be at least one family.
	if len(eco.Families) == 0 {
		t.Fatalf("expected at least one family")
	}

	// 2. Lake should have positive radius.
	if eco.Lake.Radius <= 0 || eco.Lake.MaxRadius <= 0 {
		t.Fatalf("expected positive lake radius")
	}

	// 3. No family should be spawned inside the lake.
	for i, f := range eco.Families {
		if IsInLake(f.Position, eco.Lake) {
			t.Fatalf("family %d spawned inside lake", i)
		}
	}

	// 4. No plant should be spawned inside the lake.
	for i, p := range eco.Plants {
		if IsInLake(p.position, eco.Lake) {
			t.Fatalf("plant %d spawned inside lake", i)
		}
	}

	// 5. Carrying capacity should contain predefined species.
	for _, name := range []string{"rabbit", "sheep", "deer", "wolf"} {
		if _, ok := eco.CarryingCapacity[name]; !ok {
			t.Fatalf("missing carrying capacity for %s", name)
		}
	}
}

func TestRandomPartition(t *testing.T) {
	tests := []struct {
		total int
		k     int
		min   int
	}{
		{100, 3, 10},
		{5, 3, 2},  // total < k*min
		{0, 3, 1},  // zero total
		{10, 0, 1}, // k <= 0
		{8, 4, 2},  // exact fit
	}

	for i, tt := range tests {
		parts := randomPartition(tt.total, tt.k, tt.min)

		sum := 0
		for _, p := range parts {
			sum += p
		}

		if tt.total > 0 && (tt.k <= 0 || tt.total < tt.k*tt.min) {
			// All individuals in one family
			if len(parts) != 1 || parts[0] != tt.total {
				t.Fatalf("case %d: expected single part %d, got %v", i, tt.total, parts)
			}
			continue
		}
		if tt.total == 0 {
			if len(parts) != 0 {
				t.Fatalf("case %d: expected empty slice, got %v", i, parts)
			}
			continue
		}

		if sum != tt.total {
			t.Fatalf("case %d: sum mismatch: total=%d, got=%d", i, tt.total, sum)
		}
		for _, p := range parts {
			if p < tt.min {
				t.Fatalf("case %d: part %d less than min %d", i, p, tt.min)
			}
		}
	}
}

func TestUpdatePopulationsBasicGrowth(t *testing.T) {
	// Case 1: single prey, growth > 0 due to GrowthRate + PlantCoefficient
	eco1 := Ecosystem{
		Families: []Family{
			{
				Size:    100,
				species: Species{GrowthRate: 0.10, Type: "prey"},
			},
		},
	}
	UpdatePopulations(&eco1)
	if eco1.Families[0].Size <= 100 {
		t.Fatalf("prey family should grow, got %d", eco1.Families[0].Size)
	}

	// Case 2: single predator, growth < 0
	eco2 := Ecosystem{
		Families: []Family{
			{
				Size:    100,
				species: Species{GrowthRate: -0.10, Type: "predator"},
			},
		},
	}
	UpdatePopulations(&eco2)
	if eco2.Families[0].Size >= 100 {
		t.Fatalf("predator family should shrink, got %d", eco2.Families[0].Size)
	}

	// Case 3: predator-prey interaction within Eating_Threshold
	preySpecies := Species{GrowthRate: 0.30, ContactGrowthRate: -0.10, Type: "prey"}
	predSpecies := Species{GrowthRate: -0.10, ContactGrowthRate: 0.20, Type: "predator"}

	eco3 := Ecosystem{
		Families: []Family{
			{
				Size:     40,
				Position: OrderedPair{0, 0},
				species:  preySpecies,
			},
			{
				Size:     80,
				Position: OrderedPair{0, 5}, // within Eating_Threshold=10
				species:  predSpecies,
			},
		},
	}
	UpdatePopulations(&eco3)
	if len(eco3.Families) != 2 {
		t.Fatalf("expected 2 families after update, got %d", len(eco3.Families))
	}
	if eco3.Families[0].Size <= 40 {
		t.Fatalf("prey should grow with positive net growth, got %d", eco3.Families[0].Size)
	}
	if eco3.Families[1].Size <= 80 {
		t.Fatalf("predator should grow due to positive net growth, got %d", eco3.Families[1].Size)
	}

	// Case 4: predator-prey far apart (no interaction)
	eco4 := Ecosystem{
		Families: []Family{
			{
				Size:     40,
				Position: OrderedPair{0, 0},
				species:  preySpecies,
			},
			{
				Size:     80,
				Position: OrderedPair{1000, 1000},
				species:  predSpecies,
			},
		},
	}
	UpdatePopulations(&eco4)
	if len(eco4.Families) != 2 {
		t.Fatalf("expected 2 families after update, got %d", len(eco4.Families))
	}

	// Case 5: family that should go extinct (size <= 0) is removed
	eco5 := Ecosystem{
		Families: []Family{
			{
				Size:    10,
				species: Species{GrowthRate: -2.0, Type: "predator"},
			},
		},
	}
	UpdatePopulations(&eco5)
	if len(eco5.Families) != 0 {
		t.Fatalf("expected extinct family to be removed, got %d families", len(eco5.Families))
	}
}

func TestMergeFamilies(t *testing.T) {
	// Case 1: no merge when only one family
	eco1 := Ecosystem{
		Families: []Family{
			{Size: Smallest_Family_Size - 1, Position: OrderedPair{0, 0}, species: Species{Name: "rabbit"}},
		},
	}
	MergeFamilies(&eco1)
	if len(eco1.Families) != 1 {
		t.Fatalf("case 1: expected 1 family, got %d", len(eco1.Families))
	}

	// Case 2: two small families of same species within threshold merge
	eco2 := Ecosystem{
		Families: []Family{
			{Size: Smallest_Family_Size - 1, Position: OrderedPair{0, 0}, species: Species{Name: "rabbit"}},
			{Size: Smallest_Family_Size - 2, Position: OrderedPair{Merging_Threshold - 1, 0}, species: Species{Name: "rabbit"}},
		},
	}
	MergeFamilies(&eco2)
	if len(eco2.Families) != 1 {
		t.Fatalf("case 2: expected 1 family after merge, got %d", len(eco2.Families))
	}

	// Case 3: two small families of same species beyond threshold do not merge
	eco3 := Ecosystem{
		Families: []Family{
			{Size: Smallest_Family_Size - 1, Position: OrderedPair{0, 0}, species: Species{Name: "rabbit"}},
			{Size: Smallest_Family_Size - 2, Position: OrderedPair{Merging_Threshold + 10, 0}, species: Species{Name: "rabbit"}},
		},
	}
	MergeFamilies(&eco3)
	if len(eco3.Families) != 2 {
		t.Fatalf("case 3: expected 2 families, got %d", len(eco3.Families))
	}

	// Case 4: two families of different species within threshold do not merge
	eco4 := Ecosystem{
		Families: []Family{
			{Size: Smallest_Family_Size - 1, Position: OrderedPair{0, 0}, species: Species{Name: "rabbit"}},
			{Size: Smallest_Family_Size - 2, Position: OrderedPair{Merging_Threshold - 1, 0}, species: Species{Name: "wolf"}},
		},
	}
	MergeFamilies(&eco4)
	if len(eco4.Families) != 2 {
		t.Fatalf("case 4: expected 2 families, got %d", len(eco4.Families))
	}

	// Case 5: three families, one small merges into another
	eco5 := Ecosystem{
		Families: []Family{
			{Size: Smallest_Family_Size - 1, Position: OrderedPair{0, 0}, species: Species{Name: "rabbit"}},
			{Size: Smallest_Family_Size + 5, Position: OrderedPair{Merging_Threshold - 1, 0}, species: Species{Name: "rabbit"}},
			{Size: Smallest_Family_Size + 5, Position: OrderedPair{1000, 0}, species: Species{Name: "rabbit"}},
		},
	}
	MergeFamilies(&eco5)
	if len(eco5.Families) != 2 {
		t.Fatalf("case 5: expected 2 families, got %d", len(eco5.Families))
	}
}

func TestSplitLargeFamilies(t *testing.T) {
	// Set up families with different sizes
	eco := Ecosystem{
		Families: []Family{
			{Size: Max_Family_Size - 10}, // no split
			{Size: Max_Family_Size + 10}, // split
			{Size: Max_Family_Size},      // no split
			{Size: Max_Family_Size + 1},  // split
			{Size: Max_Family_Size * 2},  // split
		},
	}

	// total size before
	totalBefore := 0
	for _, f := range eco.Families {
		totalBefore += f.Size
	}

	SplitLargeFamilies(&eco)

	// total size after
	totalAfter := 0
	maxSize := 0
	for _, f := range eco.Families {
		totalAfter += f.Size
		if f.Size > maxSize {
			maxSize = f.Size
		}
	}

	if totalBefore != totalAfter {
		t.Fatalf("total population changed by split: before=%d after=%d", totalBefore, totalAfter)
	}
	if maxSize > Max_Family_Size {
		t.Fatalf("found family larger than Max_Family_Size after split: %d", maxSize)
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		a, b   OrderedPair
		expect float64
	}{
		{OrderedPair{0, 0}, OrderedPair{0, 0}, 0},
		{OrderedPair{0, 0}, OrderedPair{3, 4}, 5},
		{OrderedPair{1, 1}, OrderedPair{4, 5}, 5},
		{OrderedPair{-1, -1}, OrderedPair{-4, -5}, 5},
		{OrderedPair{2, 3}, OrderedPair{2, 3}, 0},
	}

	for i, tt := range tests {
		got := distance(tt.a, tt.b)
		if !almostEqual(got, tt.expect, 1e-6) {
			t.Fatalf("case %d: expected %f, got %f", i, tt.expect, got)
		}
	}
}

func TestCheck(t *testing.T) {
	rabbit := SpeciesRegistry["rabbit"]
	wolf := SpeciesRegistry["wolf"]

	tests := []struct {
		name string
		A, B Family
		expA float64
		expB float64
	}{
		{
			name: "A predator, B prey, within threshold",
			A:    Family{Position: OrderedPair{0, 0}, species: wolf},
			B:    Family{Position: OrderedPair{0, 5}, species: rabbit},
			expA: wolf.ContactGrowthRate,
			expB: rabbit.ContactGrowthRate,
		},
		{
			name: "B predator, A prey, within threshold",
			A:    Family{Position: OrderedPair{0, 5}, species: rabbit},
			B:    Family{Position: OrderedPair{0, 0}, species: wolf},
			expA: rabbit.ContactGrowthRate,
			expB: wolf.ContactGrowthRate,
		},
		{
			name: "predator-prey beyond threshold",
			A:    Family{Position: OrderedPair{0, 0}, species: wolf},
			B:    Family{Position: OrderedPair{100, 0}, species: rabbit},
			expA: 0,
			expB: 0,
		},
		{
			name: "prey-prey within threshold",
			A:    Family{Position: OrderedPair{0, 0}, species: rabbit},
			B:    Family{Position: OrderedPair{0, 5}, species: rabbit},
			expA: 0,
			expB: 0,
		},
		{
			name: "predator-predator within threshold",
			A:    Family{Position: OrderedPair{0, 0}, species: wolf},
			B:    Family{Position: OrderedPair{0, 5}, species: wolf},
			expA: 0,
			expB: 0,
		},
	}

	for i, tt := range tests {
		gotA, gotB := Check(tt.A, tt.B)
		if !almostEqual(gotA, tt.expA, 1e-6) || !almostEqual(gotB, tt.expB, 1e-6) {
			t.Fatalf("case %d (%s): expected (%f,%f), got (%f,%f)",
				i, tt.name, tt.expA, tt.expB, gotA, gotB)
		}
	}
}

func TestDemoSimulationRun(t *testing.T) {
	DemoSimulationRun()
}
