package main

import (
	"image"
	"math"
	"programingProject_main/canvas"
)

// Config contains customizable parameters for the animation
type Config struct {
	CanvasWidth     int
	AgentColor      Color // Generic color parameter
	BackgroundColor Color
}

// Color represents an RGB color with an optional alpha component
type Color struct {
	R, G, B, A uint8
}

// AnimateSystem takes a collection of Ecosystem states and generates an image for each specified time point.
func AnimateSystem(timePoints []Ecosystem, config Config, drawingFrequency int) []image.Image {
	var images []image.Image

	for i, eco := range timePoints {
		if i%drawingFrequency == 0 {
			img := DrawToCanvas(eco, config)
			images = append(images, img)
		}
	}

	return images
}

// DrawToCanvas renders a single state of the ecosystem onto a canvas.
func DrawToCanvas(ecosystem Ecosystem, config Config) image.Image {
	c := canvas.CreateNewCanvas(config.CanvasWidth, config.CanvasWidth)

	// Set background color and draw weather label based on the current ecosystem's weather.
	DrawWeatherBackground(&c, ecosystem.weather, config)

	// Draw the lake
	lakeColor := canvas.MakeColor(64, 164, 223) // A nice blue color for the lake
	c.SetFillColor(lakeColor)
	canvasX := (ecosystem.Lake.Position.x / ecosystem.width) * float64(config.CanvasWidth)
	canvasY := (ecosystem.Lake.Position.y / ecosystem.width) * float64(config.CanvasWidth)
	canvasRadius := (ecosystem.Lake.Radius / ecosystem.width) * float64(config.CanvasWidth)
	c.Circle(canvasX, canvasY, canvasRadius)
	c.Fill()

	// --- 植物繪製已根據需求停用 ---
	// // Draw the plants
	// for _, p := range ecosystem.Plants {
	// 	DrawPlant(&c, p, config, ecosystem.width)
	// }

	for _, f := range ecosystem.Families {
		// Draw each family
		DrawFamily(&c, f, config, ecosystem.width)
	}
	// DrawLegend(c, config) // Temporarily disable the legend to remove the square at the top-left.
	return c.GetImage()
}

// DrawPlant draws a single plant on the canvas.
func DrawPlant(c *canvas.Canvas, p Plant, config Config, ecosystemWidth float64) {
	// We can represent plants as small green circles.
	// The radius can be based on the plant's size.
	if p.size <= 0 {
		return // Don't draw dead plants
	}
	radius := math.Sqrt(p.size) * 0.5           // Use a smaller multiplier for plants
	plantColor := canvas.MakeColor(34, 139, 34) // ForestGreen
	c.SetFillColor(plantColor)
	canvasX := (p.position.x / ecosystemWidth) * float64(config.CanvasWidth)
	canvasY := (p.position.y / ecosystemWidth) * float64(config.CanvasWidth)
	c.Circle(canvasX, canvasY, radius)
	c.Fill()
}

// DrawFamily draws the family on the canvas as a circle.
func DrawFamily(c *canvas.Canvas, f Family, config Config, ecosystemWidth float64) {
	// The size of the circle will represent the family size
	// Further reduce the multiplier to make the circles even smaller, creating a sense of a larger space.
	radius := math.Sqrt(float64(f.Size)) * 0.8 // Use square root to make size differences less extreme
	// Get color based on species
	speciesColor := GetColorForSpecies(f.species.Name)
	c.SetFillColor(canvas.MakeColor(speciesColor.R, speciesColor.G, speciesColor.B))
	// Calculate position on canvas
	canvasX := (f.Position.x / ecosystemWidth) * float64(config.CanvasWidth)
	canvasY := (f.Position.y / ecosystemWidth) * float64(config.CanvasWidth)

	// Draw a different shape based on the species type.
	switch f.species.Name {
	case "wolf":
		// Draw a triangle for predators
		c.MoveTo(canvasX, canvasY-radius) // Top point
		c.LineTo(canvasX-radius, canvasY+radius)
		c.LineTo(canvasX+radius, canvasY+radius)
		c.LineTo(canvasX, canvasY-radius) // Close the path

	case "deer":
		// Draw a diamond (rhombus) for deer
		c.MoveTo(canvasX, canvasY-radius)     // Top point
		c.LineTo(canvasX-radius*0.7, canvasY) // Left point
		c.LineTo(canvasX, canvasY+radius)     // Bottom point
		c.LineTo(canvasX+radius*0.7, canvasY) // Right point
		c.LineTo(canvasX, canvasY-radius)     // Close the path

	case "human":
		// Draw a five-pointed star for humans
		// Intentionally make the human icon larger for better visibility,
		// regardless of their small population size.
		humanSizeMultiplier := 3.0
		radius *= humanSizeMultiplier

		outerRadius := radius * 1.5
		innerRadius := radius * 0.6
		c.MoveTo(canvasX, canvasY-outerRadius) // Start at the top point
		for i := 1; i <= 5; i++ {
			// Inner point
			angle := float64(i*2-1)*math.Pi/5.0 - math.Pi/2
			c.LineTo(canvasX+innerRadius*math.Cos(angle), canvasY+innerRadius*math.Sin(angle))
			// Outer point
			angle = float64(i*2)*math.Pi/5.0 - math.Pi/2
			c.LineTo(canvasX+outerRadius*math.Cos(angle), canvasY+outerRadius*math.Sin(angle))
		}

	case "sheep":
		// Draw a hexagon for sheep
		c.MoveTo(canvasX+radius, canvasY) // Start at the right-most point
		for i := 1; i <= 6; i++ {
			angle := float64(i) * 2.0 * math.Pi / 6.0
			x := canvasX + radius*math.Cos(angle)
			y := canvasY + radius*math.Sin(angle)
			c.LineTo(x, y)
		}

	default:
		// Draw a circle for rabbits and any other default species
		c.Circle(canvasX, canvasY, radius)
	}

	c.Fill()
}

func GetColorForSpecies(name string) Color {
	switch name {
	case "rabbit":
		return Color{R: 255, G: 200, B: 200} // Light red
	case "sheep":
		return Color{R: 200, G: 255, B: 200} // Light green
	case "deer":
		return Color{R: 200, G: 200, B: 255} // Light blue
	case "wolf":
		return Color{R: 100, G: 100, B: 100} // Dark gray
	case "human":
		return Color{R: 255, G: 255, B: 0} // Yellow
	default:
		return Color{R: 255, G: 255, B: 255} // White
	}
}

func DrawLegend(c canvas.Canvas, config Config) {
	speciesList := []string{"rabbit", "sheep", "deer", "wolf", "human"}
	legendX := 10.0
	legendY := 10.0

	for _, name := range speciesList {
		color := GetColorForSpecies(name)
		// Draw color box
		c.SetFillColor(canvas.MakeColor(color.R, color.G, color.B))
		c.MoveTo(legendX, legendY)
		c.LineTo(legendX+10, legendY)
		c.LineTo(legendX+10, legendY+10)
		c.LineTo(legendX, legendY+10)
		c.LineTo(legendX, legendY)
		c.Fill()

		// Draw label next to the box
		c.SetStrokeColor(canvas.MakeColor(0, 0, 0))
		//c.DrawText(startX+15, y+10, name)
	}
}
