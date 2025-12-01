package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"programingProject_main/gifhelper" // Import the gifhelper package
	"strconv"
)

func main() {
	fmt.Println("Starting Ecosystem Simulation!")

	// Parse command line arguments.
	// Usage: go run . [numGens] [timeStep] [canvasWidth] [imageFrequency]
	numGens, _ := strconv.Atoi(os.Args[1])
	timeStep, _ := strconv.ParseFloat(os.Args[2], 64)
	canvasWidth, _ := strconv.ParseFloat(os.Args[3], 64)
	imageFrequency, _ := strconv.Atoi(os.Args[4])

	// Initialize ecosystem using the combined approach
	initialEcosystem := InitializeEcosystem()

	// Run simulation using ecosystem dynamics
	timePoints := SimulateEcosystem(initialEcosystem, numGens+1, timeStep)

	// --- Create a CSV file to log population data for R Shiny ---
	logFile, err := os.Create("population_log.csv")
	if err != nil {
		log.Fatalf("failed to create log file: %s", err)
	}
	defer logFile.Close()

	csvWriter := csv.NewWriter(logFile)
	defer csvWriter.Flush()

	// Write header row
	header := []string{"Generation", "rabbit", "sheep", "deer", "wolf", "human", "plant_mass"}
	csvWriter.Write(header)

	// Write data rows and print to console
	for i, ecosystem := range timePoints {
		counts := CountSpecies(&ecosystem)
		plantMass := CountPlantMass(&ecosystem) // Calculate plant mass
		// Print to console (optional, but good for real-time feedback)
		fmt.Printf("t=%d, rabbit=%d, sheep=%d, deer=%d, wolf=%d, human=%d, plants=%.2f\n", i, counts["rabbit"], counts["sheep"], counts["deer"], counts["wolf"], counts["human"], plantMass)

		// Prepare row for CSV
		row := []string{
			strconv.Itoa(i),
			strconv.Itoa(counts["rabbit"]),
			strconv.Itoa(counts["sheep"]),
			strconv.Itoa(counts["deer"]),
			strconv.Itoa(counts["wolf"]),
			strconv.Itoa(counts["human"]),
			strconv.FormatFloat(plantMass, 'f', 2, 64), // Add plant mass to the row
		}
		csvWriter.Write(row)
	}
	fmt.Println("Population data saved to population_log.csv")

	// Defining configuration settings for animation.
	config := Config{
		CanvasWidth:     int(canvasWidth),
		AgentColor:      Color{R: 255, G: 255, B: 255, A: 255}, // Generic color
		BackgroundColor: Color{R: 173, G: 216, B: 230},         // Light blue background
	}

	// Generate images for the animation
	// For the visualization, we are only going to draw every nth board to be more efficient
	imageList := AnimateSystem(timePoints, config, imageFrequency)
	outFile := "Animal_Sim.gif"
	gifhelper.ImagesToGIF(imageList, outFile) // code is given
	fmt.Println("GIF drawn!")

	fmt.Println("Ecosystem simulation completed!")

}

//Victor- Movement/Behavior of Animals

// example arg: go run . 1000 0.1 500 10 (numGens, timeStep, canvasWidth, imageFrequency)
