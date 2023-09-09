package config

import (
	"os"
	"testing"
)

func Test_LoadDotEnv(t *testing.T) {
	err := LoadEnv("../../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	foo := os.Getenv("TEST")
	if foo != "foo" {
		t.Errorf("Error loading .env file")
	}
}
