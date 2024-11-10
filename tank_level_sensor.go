package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

// Define pins
var (
	triggerPin = rpi.P1_11 // GPIO 17
	echoPin    = rpi.P1_13 // GPIO 27
)

// Sound speed in cm/µs
const soundSpeed = 0.0343

func main() {
	// Initialize periph library
	if _, err := host.Init(); err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}

	// Configure trigger as output and echo as input
	trigger := triggerPin.(gpio.PinOut) // Cast to correct type
	echo := echoPin.(gpio.PinIn)        // Cast to correct type
	trigger.Out(gpio.Low)

	// Infinite loop to measure distance
	for {
		distance, err := measureDistance(trigger, echo)
		if err != nil {
			fmt.Printf("Error measuring distance: %v\n", err)
		} else {
			fmt.Printf("Distance to water surface: %.2f cm\n", distance)
		}

		// Sleep for a second between measurements
		time.Sleep(1 * time.Second)
	}
}

// measureDistance triggers a pulse and measures the echo time
func measureDistance(trigger gpio.PinOut, echo gpio.PinIn) (float64, error) {
	// Send 10µs pulse to trigger
	trigger.Out(gpio.High)
	time.Sleep(10 * time.Microsecond)
	trigger.Out(gpio.Low)

	// Wait for echo to go high
	start := time.Now()
	for echo.Read() == gpio.Low {
		start = time.Now()
	}

	// Measure time until echo goes low again
	for echo.Read() == gpio.High {
	}

	duration := time.Since(start).Microseconds()
	if duration == 0 {
		return 0, fmt.Errorf("no pulse detected")
	}

	// Calculate distance based on duration
	distance := float64(duration) * soundSpeed / 2 // Divide by 2 for round-trip

	return distance, nil
}
