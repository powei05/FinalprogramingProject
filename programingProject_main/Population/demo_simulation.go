package main

// DemoSimulationRun is an optional demo entry point that shows how to
// execute the configurable simulation pipeline without affecting
// the project's main execution flow.
func DemoSimulationRun() {
	cfg := NewDefaultEcosystemConfig()
	cfg.Weather.InitialWeather = "Sunny"

	RunConfiguredSimulation(cfg, 5)
}
