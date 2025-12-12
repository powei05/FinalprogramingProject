# Ecosystem Simulation Project

## Overview
This project is a high-performance ecosystem simulation written in **Go**, with data visualization and analysis scripts provided in **R**. 

The simulation models population dynamics and interactions within a defined environment, tracking the lifecycle and behavior of different agents over distinct time steps. The results are exported to csv filewhich are then processed by R to generate visual trends of the population growth and decline.

## Prerequisites
Before running the project, ensure you have the following installed: Go, R, R Packages(ggplot2 or tidyverse).

## How to run
1. Navigate to the main project directory and run the Go program. This will execute the simulation logic and generate the data logs.

```bash
cd programingProject_main/Population
.\population 1000 0.1 500 10
```

where .\population numGens timestep canvasWidth imageFrequency

2. Visualize the Results

  2.1 Population curves by Rshiny
  ```r
  setwd("programingProject_main/Population")
  shiny::runApp()
  ```
  2.2 auto created animation
  proogramingProject_main/population/Animal_Sim.gif.out.gif
  
## Contributors
Po-Wei Chang, An Wang, Victor Chiu, Luke Smallwood


