package main

import (
	"math/rand"
	"programingProject_main/canvas"
)

// add weather field to Ecosystem struct

// function to update weather randomly
func (e *Ecosystem) UpdateWeather() {
	choices := []string{"Dry", "Sunny", "Rainy", "Frozen"}
	e.weather = choices[rand.Intn(len(choices))]
}

// functions to get coefficients of plant increasing based on weather, when using, multiply the base rate with (1 + coefficient)
func CoefficientOfPlantIncrease(weather string) float64 {
	switch weather {
	case "Dry":
		return -0.20
	case "Sunny":
		return 0.00
	case "Rainy":
		return 0.20
	default: // "Frozen"
		return -0.40
	}
}

// functions to get coefficients of lake increasing based on weather, when using, multiply the base rate with (1 + coefficient)
func CoefficientOfLakeIncrease(weather string) float64 {
	switch weather {
	case "Dry":
		return -0.20
	case "Sunny":
		return 0.00
	case "Rainy":
		return 0.20
	default: // "Frozen"
		return 0.00
	}
}

// functions to get coefficients of moving speed increase based on weather, when using, multiply the base rate with (1 + coefficient)
func CoefficientOfMovingSpeedIncrease(weather string) float64 {
	// The values have been reduced to make the effect more subtle and natural.
	switch weather {
	case "Dry":
		return -0.05 // Slight slowdown in dry conditions.
	case "Sunny":
		return 0.07 // A gentle speed boost in good weather.
	case "Rainy":
		return 0.00
	default: // "Frozen"
		return -0.15 // A noticeable, but not crippling, slowdown in frozen weather.
	}
}

// functions to get coefficients of animal growth rate increase based on weather, when using, multiply the base rate with (1 + coefficient)
func CoefficientOfAnimalGrowthRateIncrease(weather string) float64 {
	switch weather {
	case "Dry":
		return -0.10
	case "Sunny":
		return 0.10
	case "Rainy":
		return 0.00
	default: // "Frozen"
		return -0.20
	}
}

// function to draw weather background and label on the canvas, can be called inside the function DrawToCanvas in drawing.go
func DrawWeatherBackground(c *canvas.Canvas, weather string, config Config) {
	var col Color
	switch weather {
	case "Frozen":
		col = Color{0, 0, 164, 255} // darker blue
	case "Sunny":
		col = Color{253, 112, 43, 158} // orange
	case "Rainy":
		col = Color{74, 106, 125, 158} // grayish blue
	default: // Dry
		col = Color{159, 0, 0, 181} // dark red
	}

	c.SetFillColor(canvas.MakeColor(col.R, col.G, col.B))
	c.ClearRect(0, 0, config.CanvasWidth, config.CanvasWidth)
	c.Fill()

	// c.SetFillColor(canvas.MakeColor(0, 0, 0))
	// c.SetFont("Arial", 20)
	// c.FillText(10, 25, weather)
}
