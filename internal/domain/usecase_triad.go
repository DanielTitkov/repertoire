package domain

import (
	"errors"
)

// MoveFromLeft moves term from left to right
func (t *Triad) MoveFromLeft(index int) error {
	if index >= len(t.LeftTerms) {
		return errors.New("index out of range")
	}

	t.RightTerms = append(t.RightTerms, t.LeftTerms[index])

	newLeft := make([]*Term, 0)
	newLeft = append(newLeft, t.LeftTerms[:index]...)
	newLeft = append(newLeft, t.LeftTerms[index+1:]...)
	t.LeftTerms = newLeft

	return nil
}

func (t *Triad) MoveFromRight(index int) error {
	if index >= len(t.RightTerms) {
		return errors.New("index out of range")
	}

	t.LeftTerms = append(t.LeftTerms, t.RightTerms[index])

	newRight := make([]*Term, 0)
	newRight = append(newRight, t.RightTerms[:index]...)
	newRight = append(newRight, t.RightTerms[index+1:]...)
	t.RightTerms = newRight

	return nil
}
