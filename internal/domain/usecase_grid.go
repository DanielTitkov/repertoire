package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"

	"github.com/DanielTitkov/repertoire/internal/util"
)

func NewGrid(
	cfg GridConfig,
) *Grid {
	grid := &Grid{
		Config: cfg,
		Step:   GridStepLinking, // FIXME
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
	fmt.Println(grid.InitMatrix()) // FIXME

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

	g.Matrix = mat.NewDense(len(g.Constructs), len(g.Terms), nil)

	return nil
}

func (g *Grid) GetTriadByIndex(idx int) *Triad {
	// TODO add error on bad index
	return g.Triads[idx]
}
