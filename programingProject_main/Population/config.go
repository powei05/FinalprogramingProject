package main

type MovementConfig struct {
	MaxSpeed         float64
	TimeStep         float64
	SeparationWeight float64
}

type PopulationConfig struct {
	CarryingCapacities map[string]int
	InitialPopulations map[string]int
	MinFamilySize      int
	MaxFamilySize      int
}

type WeatherConfig struct {
	Enabled        bool
	StepInterval   int
	InitialWeather string
}

type LakeConfig struct {
	Radius    float64
	MaxRadius float64
	Center    OrderedPair
}

type EcosystemConfig struct {
	Width      float64
	Movement   MovementConfig
	Population PopulationConfig
	Weather    WeatherConfig
	Lake       LakeConfig
}

func NewDefaultMovementConfig() MovementConfig {
	return MovementConfig{
		MaxSpeed:         Max_Family_Speed,
		TimeStep:         1.0,
		SeparationWeight: 2.0,
	}
}

func NewDefaultPopulationConfig() PopulationConfig {
	cc := make(map[string]int)
	for k, v := range initialPopulations {
		cc[k] = int(float64(v) * 1.5)
	}
	return PopulationConfig{
		CarryingCapacities: cc,
		InitialPopulations: initialPopulations,
		MinFamilySize:      Smallest_Family_Size,
		MaxFamilySize:      Max_Family_Size,
	}
}

func NewDefaultWeatherConfig() WeatherConfig {
	return WeatherConfig{
		Enabled:        true,
		StepInterval:   Weather_Change_Interval,
		InitialWeather: "Sunny",
	}
}

func NewDefaultLakeConfig() LakeConfig {
	r := Ecosystem_Width * 0.15
	center := OrderedPair{x: Ecosystem_Width / 2.0, y: Ecosystem_Width / 2.0}
	return LakeConfig{
		Radius:    r,
		MaxRadius: r,
		Center:    center,
	}
}

func NewDefaultEcosystemConfig() EcosystemConfig {
	return EcosystemConfig{
		Width:      Ecosystem_Width,
		Movement:   NewDefaultMovementConfig(),
		Population: NewDefaultPopulationConfig(),
		Weather:    NewDefaultWeatherConfig(),
		Lake:       NewDefaultLakeConfig(),
	}
}

func (c EcosystemConfig) Clone() EcosystemConfig {
	cc := make(map[string]int)
	for k, v := range c.Population.CarryingCapacities {
		cc[k] = v
	}
	ip := make(map[string]int)
	for k, v := range c.Population.InitialPopulations {
		ip[k] = v
	}
	c.Population.CarryingCapacities = cc
	c.Population.InitialPopulations = ip
	return c
}
