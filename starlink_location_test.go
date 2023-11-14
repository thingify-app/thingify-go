package main

import (
	"testing"
)

func TestNewStarlinkLocationProvider(t *testing.T) {
	provider, err := NewStarlinkLocationProvider()
	if err != nil {
		t.Fatal(err)
	}

	location, err := provider.GetLocation()
	if err != nil {
		t.Fatal(err)
	}

	// assert location is 0, 0
	if location.Latitude != 0 {
		t.Errorf("expected latitude to be 0, got %f", location.Latitude)
	}
	if location.Longitude != 0 {
		t.Errorf("expected longitude to be 0, got %f", location.Longitude)
	}
}
