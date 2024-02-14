package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
)

type Person struct {
	Name    string `json:"name"`
	GUID    string `json:"guid"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// go run generate_random_json.go 10000

func main() {
	// Check if the number of objects is provided as a command-line argument
	if len(os.Args) < 2 {
		color.HiRed("Usage: go run script.go <number_of_objects>")
		os.Exit(1)
	}

	// Parse the number of objects from command-line arguments
	numObjectsStr := os.Args[1]
	numObjects, err := strconv.Atoi(numObjectsStr)
	if err != nil {
		color.HiRed("Invalid number of objects provided")
		os.Exit(1)
	}

	// Create an array to hold the generated data
	var people []Person

	// Generate random data and populate the array
	for i := 0; i < numObjects; i++ {
		person := Person{
			Name:    faker.Name(),
			GUID:    faker.UUIDDigit(),
			Company: faker.Gender(),
			Email:   faker.Email(),
			Address: faker.ID,
		}
		people = append(people, person)
	}

	// Convert the array to JSON
	jsonData, err := json.MarshalIndent(people, "", "    ")
	if err != nil {
		color.HiRed("Error marshaling JSON:", err)
		os.Exit(1)
	}

	// Write the JSON data to a file
	file, err := os.Create("random_data.json")
	if err != nil {
		color.HiRed("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(jsonData)

	color.HiGreen("JSON file generated successfully: random_data.json")
}
