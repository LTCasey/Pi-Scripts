package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	soundSpeed          = 0.0343                // Sound speed in cm/µs
	pulseDuration       = 10 * time.Microsecond // Pulse duration for trigger
	measurementInterval = 1 * time.Second       // Interval between measurements
)

var (
	triggerPin gpio.PinOut // GPIO pin connected to the trigger
	echoPin    gpio.PinIn  // GPIO pin connected to the echo
)

func initPins() error {

	// Initialize periph library
	if _, err := host.Init(); err != nil {
		return fmt.Errorf("failed to initialize periph: %w", err)
	}

	// Configure trigger as output and echo as input
	triggerPin, ok := rpi.P1_11.(gpio.PinOut) // GPIO 17
	if !ok {
		return fmt.Errorf("failed to cast P1_11 to gpio.PinOut")
	}
	echoPin, ok = rpi.P1_13.(gpio.PinIn) // GPIO 27
	if !ok {
		return fmt.Errorf("failed to cast P1_13 to gpio.PinIn")
	}

	triggerPin.Out(gpio.Low)
	return nil
}

func main() {
	if err := initPins(); err != nil {
		log.Fatalf("%v", err)
	}

	for {
		distance, err := measureDistance()
		if err != nil {
			fmt.Printf("Error measuring distance: %v\n", err)
		} else {
			fmt.Printf("Distance to water surface: %.2f cm\n", distance)
		}

		time.Sleep(measurementInterval)
	}
}

// measureDistance triggers a pulse and measures the echo time
func measureDistance() (float64, error) {
	// Send 10µs pulse to trigger
	triggerPin.Out(gpio.High)
	time.Sleep(pulseDuration)
	triggerPin.Out(gpio.Low)

	start := time.Now()
	var duration int64

	// Wait for echo to go high
	for echoPin.Read() == gpio.Low {
		if time.Since(start) > measurementInterval {
			return 0, fmt.Errorf("timeout waiting for echo")
		}
	}

	// Measure time until echo goes low again
	start = time.Now()
	for echoPin.Read() == gpio.High {
		if time.Since(start) > measurementInterval {
			return 0, fmt.Errorf("timeout waiting for echo to go low")
		}
	}
	duration = time.Since(start).Microseconds()

	if duration == 0 {
		return 0, fmt.Errorf("no pulse detected")
	}

	// Calculate distance based on duration
	distance := float64(duration) * soundSpeed / 2 // Divide by 2 for round-trip

	return distance, nil
}
