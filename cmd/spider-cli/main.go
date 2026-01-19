package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aldernero/spider"
)

func main() {
	var (
		configFile = flag.String("config", "", "Path to configuration file (JSON or YAML)")
		outputFile = flag.String("output", "", "Output file path (PNG or SVG)")
	)
	flag.Parse()

	if *configFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -config flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -output flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Load chart from configuration file
	chart, err := spider.NewChartFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load chart: %v", err)
	}

	// Save chart to output file
	if err := chart.Save(*outputFile); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}

	fmt.Printf("Chart saved successfully to %s\n", *outputFile)
}
