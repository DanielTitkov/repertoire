package domain

const (
	TriadMethodForced = "forcedTriads"
	TriadMethodChoice = "choiceTriads"

	GridStepTerms       = "terms"
	GridStepElicitation = "elicit"
	GridStepLinking     = "linking"

	TriadStepRaw      = "raw"      // empty object
	TriadStepInit     = "init"     // with terms on one side
	TriadStepChosen   = "chosen"   // with differing term chosend
	TriadStepLeftDone = "leftDone" // with left pole done
	TriadStepReady    = "ready"    // with both poles done
)
