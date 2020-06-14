package cmd

import (
	"math/rand"
	"time"
)

func pickRandomElement(list []string) string {
	rand.Seed(time.Now().Unix())
	return list[rand.Intn(len(list))]
}
