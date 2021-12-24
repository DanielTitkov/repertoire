package domain

import (
	"testing"
)

var constructs = []*Construct{
	{
		LeftPole:  "good",
		RightPole: "bad",
	},
	{
		LeftPole:  "sad",
		RightPole: "funny",
	},
	{
		LeftPole:  "sweet",
		RightPole: "sour",
	},
	{
		LeftPole:  "loved",
		RightPole: "hated",
	},
	{
		LeftPole:  "simple",
		RightPole: "complex",
	},
}

var terms = []Term{
	{Title: "Dad"},
	{Title: "Goose"},
	{Title: "Mom"},
	{Title: "Cat"},
	{Title: "Psychologist"},
	{Title: "Uncle Bob"},
}

func TestGenerateTriads(t *testing.T) {
	grid := NewGrid(
		GridConfig{
			MinTerms:       4,
			MaxTerms:       12,
			TriadMethod:    TriadMethodChoice,
			MinConstructs:  5,
			ConstructSteps: 5,
		},
	)

	grid.Terms = []Term{} // clear any preset terms

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

func TestCorrelationMatrices(t *testing.T) {
	grid := NewGrid(
		GridConfig{
			MinTerms:       4,
			MaxTerms:       12,
			TriadMethod:    TriadMethodChoice,
			MinConstructs:  5,
			ConstructSteps: 5,
		},
	)

	grid.Terms = terms
	grid.Constructs = constructs
	grid.Step = GridStepLinking // ?
	if err := grid.InitMatrix(); err != nil {
		t.Errorf("failed to init matrix: %s", err)
	}

	expTermsDim := len(grid.Terms)
	expConstructsDim := len(grid.Constructs)

	corrTermsM, err := grid.getCorrelationMatrix(corrDirectionTerms)
	if err != nil {
		t.Errorf("failed getting terms correlations: %s", err)
	}

	if x, y := corrTermsM.Dims(); x != expTermsDim || y != expTermsDim {
		t.Errorf("expected to get %d x %d matrix, but got %d x %d", expTermsDim, expTermsDim, x, y)
	}

	constructsCorrM, err := grid.getCorrelationMatrix(corrDirectionConstructs)
	if err != nil {
		t.Errorf("failed getting terms correlations: %s", err)
	}

	if x, y := constructsCorrM.Dims(); x != expConstructsDim || y != expConstructsDim {
		t.Errorf("expected to get %d x %d matrix, but got %d x %d", expConstructsDim, expConstructsDim, x, y)
	}

}
