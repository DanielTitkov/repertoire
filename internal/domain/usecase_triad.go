package domain

import (
	"errors"
)

func NewTriad() *Triad {
	return &Triad{
		Step: TriadStepRaw,
	}
}

func (t *Triad) AddLeftTerm(term *Term) {
	t.LeftTerms = append(t.LeftTerms, term)
	t.UpdateStep()
}

func (t *Triad) SetPoles(left, right string) {
	t.SetLeftPole(left)
	t.SetRightPole(right)
}

func (t *Triad) SetLeftPole(value string) {
	t.LeftPole = value
	t.UpdateStep()
}

func (t *Triad) SetRightPole(value string) {
	t.RightPole = value
	t.UpdateStep()
}

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

	t.UpdateStep()
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

	t.UpdateStep()
	return nil
}

func (t *Triad) UpdateStep() {
	t.Step = TriadStepInit
	if len(t.LeftTerms)+len(t.RightTerms) != 3 {
		t.Step = TriadStepRaw
		return
	}

	if t.LeftPole != "" && t.RightPole != "" {
		t.Step = TriadStepReady
		return
	}

	if t.LeftPole != "" && t.RightPole == "" {
		t.Step = TriadStepLeftDone
		return
	}

	if len(t.LeftTerms) == 2 && len(t.RightTerms) == 1 {
		t.Step = TriadStepChosen
		return
	}

	if len(t.LeftTerms) == 1 && len(t.RightTerms) == 2 {
		t.Step = TriadStepChosen
		return
	}
}
