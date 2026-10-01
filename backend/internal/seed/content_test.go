package seed

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadContent(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "content.json"), []byte(`{
  "items": [{"id":"water","name":"Water","category":"supply","stackable":true,"maxStack":10}],
  "jobs": [{"id":"general-work","name":"General Workforce","department":"Operations","positions":[]}],
  "courses": [{"id":"orientation","name":"World Orientation","durationHours":1,"educationGain":1,"requiredLevel":1}],
  "population": {"seed":42,"names":["Aster"],"jobWeights":{"general-work":100},"personalityBias":{"curiosity":5}}
}`), 0o644); err != nil {
		t.Fatal(err)
	}

	content, err := LoadContent(dir)
	if err != nil {
		t.Fatal(err)
	}
	if content.Population.Seed != 42 || len(content.Population.Names) != 1 {
		t.Fatalf("unexpected population fixture: %+v", content.Population)
	}
}

func TestLoadContentRejectsEmptyPopulation(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "content.json"), []byte(`{
  "items": [{"id":"water"}],
  "jobs": [{"id":"job"}],
  "courses": [{"id":"course"}],
  "population": {"seed":42}
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadContent(dir); err == nil {
		t.Fatal("expected validation error")
	}
}
