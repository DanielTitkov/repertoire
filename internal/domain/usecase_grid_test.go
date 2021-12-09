package domain

import "testing"

func TestGenerateTriads(t *testing.T) {
	grid := NewGrid(
		GridConfig{
			MinTerms:    5,
			MaxTerms:    6,
			TriadMethod: TriadMethodForced,
		},
	)

	for _, term := range []string{"dad", "mom", "me", "psychologist", "goose"} {
		err := grid.AddTerm(Term{Title: term})
		if err != nil {
			t.Errorf("error adding term: %s", err)
		}
	}

	err := grid.GenerateTriads()
	if err != nil {
		t.Errorf("failed generating triads: %s", err)
	}

	if len(grid.Triads) != 10 {
		t.Errorf("expected to get 10 triads, got %d", len(grid.Triads))
	}
}
