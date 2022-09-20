package game

type step struct {
	number   int
	player   player
	prevStep *step
	nextStep *step
}
