package game

import mapset "github.com/deckarep/golang-set/v2"

type color string

const (
	black  color = "black"
	red    color = "red"
	blue   color = "blue"
	orange color = "orange"
)

var colors []color = []color{black, red, blue, orange}
var colorsSet mapset.Set[color] = mapset.NewSet(colors...)
