// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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
