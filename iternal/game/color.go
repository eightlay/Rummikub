package game

import mapset "github.com/deckarep/golang-set/v2"

// Pice color
type color string

const (
	// Black piece color
	black color = "black"
	// Red piece color
	red color = "red"
	// Blue piece color
	blue color = "blue"
	// Orange piece color
	orange color = "orange"
)

// All colors
var colors []color = []color{black, red, blue, orange}

// All colors set
var colorsSet mapset.Set[color] = mapset.NewSet(colors...)
