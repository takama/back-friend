package config

import "testing"

func TestConfig(t *testing.T) {
	config := New()
	err := config.LoadFromEnv()
	if err != nil {
		t.Error("Expected loading of environment vars, got", err)
	}
}
