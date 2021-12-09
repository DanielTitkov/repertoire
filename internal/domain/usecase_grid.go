package domain

import (
	"errors"
	"fmt"
)

func NewGrid(
	cfg GridConfig,
) *Grid {
	return &Grid{
		Config: cfg,
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
