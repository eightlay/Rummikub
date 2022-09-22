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
