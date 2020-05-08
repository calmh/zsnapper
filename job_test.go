package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoadJobs(t *testing.T) {
	yaml := `
- family:    test
  # comment
  datasets:
    - data/test
  schedule: "*/20 * * * * *"
  keep:      5
  recursive: true
# comment
- family:    minute
  datasets:
    - data/test
    - data/foo/*
  schedule: "0 * * * * *"
  keep:      12
  recursive: false`

	expected := []Job{
		Job{
			Family:    "test",
			Datasets:  []string{"data/test"},
			Schedule:  "*/20 * * * * *",
			Keep:      5,
			Recursive: true,
		},
		Job{
			Family:    "minute",
			Datasets:  []string{"data/test", "data/foo/*"},
			Schedule:  "0 * * * * *",
			Keep:      12,
			Recursive: false,
		},
	}

	r := strings.NewReader(yaml)
	jobs, err := LoadJobs(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != len(expected) {
		t.Fatalf("Number of jobs %d != %d", len(jobs), len(expected))
	}
	for i := range expected {
		if !reflect.DeepEqual(jobs[i], expected[i]) {
			t.Errorf("%d: %+v != %+v", i, jobs[i], expected[i])
		}
	}
}
