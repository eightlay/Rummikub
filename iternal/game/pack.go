// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

// Pack of pieces
type pack []*Piece

// Create initial pack (bank)
func createInitialPack() pack {
	b := pack{}

	for d := 0; d < DecksNumber; d++ {
		b = append(b, createPiece(JokerNumber, JokerColor, true))

		for _, c := range colors {
			for i := MinNumber; i <= MaxNumber; i++ {
				p := createPiece(i, c, false)
				b = append(b, p)
			}
		}
	}

	return b
}
