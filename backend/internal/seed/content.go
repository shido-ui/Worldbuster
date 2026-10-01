package seed

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shido-ui/Worldbuster/backend/internal/inventory"
	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
	"github.com/shido-ui/Worldbuster/backend/internal/simulation"
)

type Content struct {
	Items      []inventory.Item
	Jobs       []job.Job
	Courses    []progression.Course
	Population PopulationProfile
}

type PopulationProfile struct {
	Seed            int64
	Names           []string
	JobWeights      map[string]int
	PersonalityBias simulation.Personality
}

func LoadContent(dir string) (Content, error) {
	if dir == "" {
		dir = "database/seed"
	}
	path := filepath.Join(dir, "content.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return Content{}, fmt.Errorf("read seed fixture %s: %w", path, err)
	}

	var content Content
	if err := json.Unmarshal(raw, &content); err != nil {
		return Content{}, fmt.Errorf("decode seed fixture %s: %w", path, err)
	}
	if err := validateContent(content); err != nil {
		return Content{}, err
	}
	return content, nil
}

func validateContent(content Content) error {
	if len(content.Items) == 0 {
		return fmt.Errorf("seed fixture contains no items")
	}
	if len(content.Jobs) == 0 {
		return fmt.Errorf("seed fixture contains no jobs")
	}
	if len(content.Courses) == 0 {
		return fmt.Errorf("seed fixture contains no courses")
	}
	if content.Population.Seed == 0 {
		return fmt.Errorf("seed fixture population seed must be non-zero")
	}
	if len(content.Population.Names) == 0 {
		return fmt.Errorf("seed fixture population contains no names")
	}
	return nil
}
