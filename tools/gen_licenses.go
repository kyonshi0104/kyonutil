//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command(
		"go", "run", "github.com/google/go-licenses@latest", "report",
		"--ignore", "github.com/kyonshi0104/kyonutil",
		"./cmd/bot",
	)

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run go-licenses: %v", err)
	}

	err = os.WriteFile("internal/discord/licenses.csv", output, 0644)
	if err != nil {
		log.Fatalf("Failed to save licenses.csv: %v", err)
	}

	log.Println("✅ licenses.csv generated successfully!")
}
