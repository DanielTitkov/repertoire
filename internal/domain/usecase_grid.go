package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"

	"github.com/DanielTitkov/repertoire/internal/util"
)

const (
	corrDirectionTerms      = "terms"
	corrDirectionConstructs = "constructs"
)

func NewGrid(
	cfg GridConfig,
) *Grid {
	grid := &Grid{
		Config: cfg,
		// Step:   GridStepTerms,
		Step: GridStepLinking, // FIXME
		Constructs: []*Construct{
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
		},
		Terms: []Term{
			{Title: "Dad"},
			{Title: "Goose"},
			{Title: "Mom"},
			{Title: "Cat"},
			{Title: "Psychologist"},
			{Title: "Uncle Bob"},
		},
	}
	fmt.Println("init matrix", grid.InitMatrix()) // FIXME
	var data []float64
	for i := 0; i < len(grid.Constructs)*len(grid.Terms); i++ {
		data = append(data, float64(rand.Intn(grid.Config.ConstructSteps)))
	}
	grid.Matrix = mat.NewDense(len(grid.Constructs), len(grid.Terms), data) // FIXME
	fmt.Println("result", grid.CalculateResult())                           // FIXME

	err := grid.Validate()
	if err != nil { // FIXME
		panic(err)
	}

	return grid
}

func (g *Grid) RemoveTermByIndex(index int) {
	if len(g.Terms) <= index {
		return
	}

	ret := make([]Term, 0)
	ret = append(ret, g.Terms[:index]...)
	g.Terms = append(ret, g.Terms[index+1:]...)
}

func (g *Grid) AddTerm(term Term) error {
	if term.Title == "" {
		return errors.New("term cannot be empty")
	}

	if len(g.Terms) >= g.Config.MaxTerms {
		return fmt.Errorf("max number of terms is %d", g.Config.MaxTerms)
	}

	for _, t := range g.Terms {
		if t.Title == term.Title {
			return errors.New("terms must be unique")
		}
	}

	g.Terms = append(g.Terms, term)
	return nil
}

func (g *Grid) GenerateTriads() error {
	var indices []int
	for i := range g.Terms {
		indices = append(indices, i)
	}

	// get all possible term index combinations by 3
	subsets := util.Combinations(indices, 3)

	// randomize subsets
	rand.Seed(time.Now().UnixNano())
	for _, a := range subsets {
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	}
	rand.Shuffle(len(subsets), func(i, j int) { subsets[i], subsets[j] = subsets[j], subsets[i] })

	// make triads
	var triads []*Triad
	for _, subset := range subsets {
		triad := NewTriad()
		for _, idx := range subset {
			triad.AddLeftTerm(&g.Terms[idx])
		}

		switch g.Config.TriadMethod {
		case TriadMethodForced:
			// move random item from left to right
			triad.MoveFromLeft(rand.Intn(3))
		case TriadMethodChoice:
		default:
			return fmt.Errorf("unknown triad method: %s", g.Config.TriadMethod)
		}
		triads = append(triads, triad)
	}

	g.Triads = triads
	if len(triads) < g.Config.MinConstructs {
		g.Config.MinConstructs = len(triads)
	}

	g.Step = GridStepElicitation

	return nil
}

func (g *Grid) GenerateConstructs() error {
	var constructs []*Construct
	for _, triad := range g.Triads {
		if triad.Step != TriadStepReady {
			continue
		}
		constructs = append(constructs, &Construct{
			LeftPole:  triad.LeftPole,
			RightPole: triad.RightPole,
		})
	}

	g.Constructs = constructs
	err := g.InitMatrix()
	if err != nil {
		return err
	}

	g.Step = GridStepLinking

	return nil
}

func (g *Grid) InitMatrix() error {
	if len(g.Terms) == 0 {
		return errors.New("grid must have terms")
	}

	if len(g.Constructs) == 0 {
		return errors.New("grid must have constructs")
	}

	var values []float64
	for i := 0; i < len(g.Constructs)*len(g.Terms); i++ {
		values = append(values, -1)
	}

	g.Matrix = mat.NewDense(len(g.Constructs), len(g.Terms), values)

	return nil
}

func (g *Grid) GetTriadByIndex(idx int) *Triad {
	// TODO add error on bad index
	return g.Triads[idx]
}

func (g *Grid) Validate() error {
	if g.Config.ConstructSteps < 2 {
		return errors.New("construct steps cannot be less than 2")
	}

	return nil
}

func (g *Grid) IsMatrixComplete() bool {
	x, y := g.Matrix.Dims()
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			if g.Matrix.At(i, j) == -1 {
				return false
			}
		}
	}

	return true
}

func (g *Grid) CalculateResult() error {
	if !g.IsMatrixComplete() {
		return errors.New("matrix must be complete in order to proceed to results")
	}

	// analysis
	termsCorrM, err := g.getCorrelationMatrix(corrDirectionTerms)
	if err != nil {
		return err
	}

	constructsCorrM, err := g.getCorrelationMatrix(corrDirectionConstructs)
	if err != nil {
		return err
	}

	g.Analysis = &GridAnalysis{
		TermsCorrMatrix:      termsCorrM,
		ConstructsCorrMatrix: constructsCorrM,
	}
	g.Step = GridStepResult
	return nil
}

func (g *Grid) getCorrelationMatrix(direction string) (mat.Matrix, error) {
	var m mat.SymDense
	res := &m

	if direction == corrDirectionTerms {
		// without transpose will get terms correlations
		stat.CorrelationMatrix(res, g.Matrix, nil)
	} else if direction == corrDirectionConstructs {
		// to get construct correlations need transpose
		stat.CorrelationMatrix(res, g.Matrix.T(), nil)
	} else {
		return nil, fmt.Errorf("direction must be either %s or %s, got %s", corrDirectionTerms, corrDirectionConstructs, direction)
	}

	return res, nil
}
