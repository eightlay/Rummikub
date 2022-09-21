package game

type pack []*Piece

func createInitialPack() pack {
	b := pack{}

	for d := 0; d < DecksNumber; d++ {
		for i := MinNumber; i <= MaxNumber; i++ {
			for _, c := range colors {
				p := createPiece(i, c, false)
				b = append(b, p)
			}
		}

		b = append(b, createPiece(JokerNumber, JokerColor, true))
	}

	return b
}
