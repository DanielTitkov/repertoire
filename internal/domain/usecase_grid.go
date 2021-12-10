package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/DanielTitkov/repertoire/internal/util"
)

func NewGrid(
	cfg GridConfig,
) *Grid {
	return &Grid{
		Config: cfg,
		Step:   GridStepTerms,
	}
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

	return nil
}

func (g *Grid) GetTriadByIndex(idx int) *Triad {
	// TODO add error on bad index
	return g.Triads[idx]
}
