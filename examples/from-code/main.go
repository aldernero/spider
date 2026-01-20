package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/aldernero/spider"
)

// normalRandom generates a normally distributed random number
// with mean and stddev, clamped to [min, max]
func normalRandom(mean, stddev, min, max float64) float64 {
	// Generate normal distribution using Box-Muller transform
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	value := mean + z*stddev

	// Clamp to [min, max]
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

func main() {
	// Create a new chart
	chart := spider.NewChart()

	// Set chart options

	chart.Options.Title = "Programmatic Spider Chart Example"

	// Add 7 axes (A through G) with max values 10^0, 10^1, ..., 10^6
	axisNames := []string{"A", "B", "C", "D", "E", "F", "G"}
	axisMaxes := make([]float64, 7)
	for i := 0; i < 7; i++ {
		axisMaxes[i] = math.Pow(10, float64(i))
		if err := chart.AddAxis(axisNames[i]); err != nil {
			log.Fatalf("Failed to add axis %s: %v", axisNames[i], err)
		}
		// Set the max value for this axis
		chart.Data.Axes[i].Max = axisMaxes[i]
	}

	// Add 5 series with normally distributed random values
	for seriesNum := 1; seriesNum <= 5; seriesNum++ {
		seriesName := fmt.Sprintf("Series %d", seriesNum)
		seriesData := make(map[string]float64)

		// Generate data for each axis
		for i, axisName := range axisNames {
			max := axisMaxes[i]
			mean := max / 2.0
			stddev := max / 6.0 // Use max/6 to keep most values within [0, max]

			value := normalRandom(mean, stddev, 0, max)
			seriesData[axisName] = value
		}

		if err := chart.AddSeries(seriesName, seriesData); err != nil {
			log.Fatalf("Failed to add series %s: %v", seriesName, err)
		}
	}

	// Save the chart
	if err := chart.Save("output.svg"); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}
	if err := chart.Save("output.png"); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}

	log.Println("Chart saved successfully to output.svg and output.png")
}
