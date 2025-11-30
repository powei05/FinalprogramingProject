package main

type Lake struct {
	Radius   float64
	Position OrderedPair // Represents the center of the circle
}

type Ecosystem struct {
	Families             []Family
	Plants               []Plant
	width                float64
	weather              string // new added
	Lake                 Lake   // Add the lake to the ecosystem
	CarryingCapacity     map[string]int
	weatherChangeCounter int
}

type Species struct {
	Name              string
	Type              string
	Class             string
	GrowthRate        float64
	ContactGrowthRate float64
}

type Family struct {
	Size                int
	MovementSpeed       OrderedPair
	Position            OrderedPair
	MovementDirection   OrderedPair
	Acceleration        OrderedPair
	PropulsionDirection OrderedPair // The family's internal "will to move" direction
	species             Species
}

type OrderedPair struct {
	x float64
	y float64
}

type Plant struct {
	position OrderedPair
	size     float64
}

var SpeciesRegistry = map[string]Species{
	"rabbit": {Name: "rabbit", Class: "prey", Type: "prey", GrowthRate: 0.3, ContactGrowthRate: -0.1},
	"sheep":  {Name: "sheep", Class: "prey", Type: "prey", GrowthRate: 0.2, ContactGrowthRate: -0.1},
	"deer":   {Name: "deer", Class: "prey", Type: "prey", GrowthRate: 0.15, ContactGrowthRate: -0.05},
	"wolf":   {Name: "wolf", Class: "predator", Type: "predator", GrowthRate: -0.1, ContactGrowthRate: 0.2},
	"human":  {Name: "human", Class: "neutral", Type: "neutral", GrowthRate: 0.0, ContactGrowthRate: 0.0},
}

var initial_family_number = 3
var initialPopulations = map[string]int{
	"rabbit": 200,
	"sheep":  150,
	"deer":   100,
	"wolf":   70,
	"human":  2,
}

var Eating_Threshold = 10.0 // when distance is less than this value, predation can occur

const Merging_Threshold = 20.0  // when distance is less than this value, families of the same species can merge
const Smallest_Family_Size = 10 // if the family size is smaller than this value, it would merge with other families
const Max_Family_Size = 80      // if the family size is larger than this value, it would split into two
const Ecosystem_Width = 500.0   // the width of the ecosystem

const Separation_Threshold = 50.0 // Proximity threshold for separation force
const Max_Family_Speed = 40.0     // Default max family speed

const PreyPlantCoefficient = 0.05 //New
const PlantCoefficient = 0.05     //New
const consumptionRate = 0.05

const Weather_Change_Interval = 100 // Weather changes every 100 steps
