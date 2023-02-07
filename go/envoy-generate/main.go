package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	snapshot := GenerateSnapshot()
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(&snapshot); err != nil {
		log.Fatal(err)
	}
}
