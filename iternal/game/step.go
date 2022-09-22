package game

// Step
//
// Contains information about step number, player,
// previous step and next step
type step struct {
	number   int
	player   player
	prevStep *step
	nextStep *step
}
