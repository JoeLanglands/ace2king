package main

import (
	"math/rand"

	"github.com/JoeLanglands/ace2king/pkg/scoring"
)

// maxInt returns the larger of x or y.
func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func shufflePlayers(players *[]scoring.Player) {
	rand.Shuffle(len(*players), func(i, j int) {
		(*players)[i], (*players)[j] = (*players)[j], (*players)[i]
	})
}
