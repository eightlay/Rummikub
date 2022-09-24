// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

// Player's hand
type hand pack

// Find largest piece number in the hand
func (h hand) largestPieceNumber() int {
	largest := MinNumber - 1
	if JokerNumber < largest {
		largest = JokerNumber
	}

	for _, v := range h {
		if v.Number > largest {
			largest = v.Number
		}
	}

	return largest
}
