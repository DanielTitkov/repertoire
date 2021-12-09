package domain

import (
	"errors"
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
	g.TermsN = len(g.Terms)
}

func (g *Grid) AddTerm(term Term) error {
	if term.Title == "" {
		return errors.New("term cannot be empty")
	}

	g.Terms = append(g.Terms, term)
	g.TermsN = len(g.Terms)
	return nil
}
