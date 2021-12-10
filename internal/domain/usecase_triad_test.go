package domain

import "testing"

func TestMoveFromLeft(t *testing.T) {
	var triad = Triad{
		LeftTerms: []*Term{
			&Term{Title: "first"},
			&Term{Title: "second"},
			&Term{Title: "third"},
		},
	}

	err := triad.MoveFromLeft(1) // move second
	if err != nil {
		t.Errorf("error while performing move from left %s", err)
	}
	if len(triad.LeftTerms) != 2 {
		t.Errorf("left len must be 2 after move")
	}
	if len(triad.RightTerms) != 1 {
		t.Errorf("right len must be 1 after move")
	}
	if first := triad.LeftTerms[0]; first.Title != "first" {
		t.Errorf("expeted to get 'first', got '%s'", first.Title)
	}
	if third := triad.LeftTerms[1]; third.Title != "third" {
		t.Errorf("expeted to get 'third', got '%s'", third.Title)
	}
	if second := triad.RightTerms[0]; second.Title != "second" {
		t.Errorf("expeted to get 'second', got '%s'", second.Title)
	}
}

func TestTriadStepInit(t *testing.T) {
	exp := TriadStepInit
	var triad = Triad{
		LeftTerms: []*Term{
			&Term{Title: "first"},
			&Term{Title: "second"},
			&Term{Title: "third"},
		},
		RightTerms: []*Term{},
	}

	triad.UpdateStep()
	if triad.Step != exp {
		t.Errorf("expected '%s' but got '%s' step value", exp, triad.Step)
	}
}

func TestTriadStepChosen(t *testing.T) {
	exp := TriadStepChosen
	var triad = Triad{
		LeftTerms: []*Term{
			&Term{Title: "first"},
			&Term{Title: "second"},
		},
		RightTerms: []*Term{
			&Term{Title: "third"},
		},
	}

	triad.UpdateStep()
	if triad.Step != exp {
		t.Errorf("expected '%s' but got '%s' step value", exp, triad.Step)
	}
}

func TestTriadStepReady(t *testing.T) {
	exp := TriadStepReady
	var triad = Triad{
		LeftTerms: []*Term{
			&Term{Title: "first"},
			&Term{Title: "second"},
		},
		RightTerms: []*Term{
			&Term{Title: "third"},
		},
		LeftPole:  "foo",
		RightPole: "bar",
	}

	triad.UpdateStep()
	if triad.Step != exp {
		t.Errorf("expected '%s' but got '%s' step value", exp, triad.Step)
	}
}

func TestTriadStepLeftDone(t *testing.T) {
	exp := TriadStepLeftDone
	var triad = Triad{
		LeftTerms: []*Term{
			&Term{Title: "first"},
			&Term{Title: "second"},
		},
		RightTerms: []*Term{
			&Term{Title: "third"},
		},
		LeftPole: "foo",
	}

	triad.UpdateStep()
	if triad.Step != exp {
		t.Errorf("expected '%s' but got '%s' step value", exp, triad.Step)
	}
}

func TestTriadStepRaw(t *testing.T) {
	exp := TriadStepRaw
	var triad = Triad{
		LeftTerms:  []*Term{},
		RightTerms: []*Term{},
	}

	triad.UpdateStep()
	if triad.Step != exp {
		t.Errorf("expected '%s' but got '%s' step value", exp, triad.Step)
	}
}

func TestMoveFromRight(t *testing.T) {
	var triad = Triad{
		LeftTerms: []*Term{
			&Term{Title: "first"},
			&Term{Title: "second"},
		},
		RightTerms: []*Term{
			&Term{Title: "third"},
		},
	}

	err := triad.MoveFromRight(0) // move "third"
	if err != nil {
		t.Errorf("error while performing move from left %s", err)
	}
	if len(triad.LeftTerms) != 3 {
		t.Errorf("left len must be 3 after move")
	}
	if len(triad.RightTerms) != 0 {
		t.Errorf("right len must be 0 after move")
	}
	if first := triad.LeftTerms[0]; first.Title != "first" {
		t.Errorf("expeted to get 'first', got '%s'", first.Title)
	}
	if second := triad.LeftTerms[1]; second.Title != "second" {
		t.Errorf("expeted to get 'second', got '%s'", second.Title)
	}
	if third := triad.LeftTerms[2]; third.Title != "third" {
		t.Errorf("expeted to get 'third', got '%s'", third.Title)
	}

}
