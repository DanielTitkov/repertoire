package domain

import "gonum.org/v1/gonum/mat"

type (
	User struct {
		Email string
		Age   int
	}
	// Grid encapsulates repetoire grid data
	Grid struct {
		Config     GridConfig
		Terms      []Term // TODO: maybe use pointer
		Constructs []*Construct
		Triads     []*Triad
		Step       string
		Matrix     *mat.Dense
		Analysis   *GridAnalysis
	}
	Term struct {
		Title string
	}
	Triad struct {
		LeftTerms  []*Term
		RightTerms []*Term
		LeftPole   string
		RightPole  string
		Step       string
	}
	Construct struct {
		LeftPole  string
		RightPole string
	}
	GridConfig struct {
		TriadMethod    string
		MinTerms       int
		MaxTerms       int
		MinConstructs  int
		ConstructSteps int
	}
	GridAnalysis struct {
		TermsCorrMatrix      mat.Matrix
		ConstructsCorrMatrix mat.Matrix
	}
)
