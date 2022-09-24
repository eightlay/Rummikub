// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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
