package config

import (
	"math/rand"
	"time"
)

// ElectionTimeout in 1000 ms
func ElectionTimeout() int {
	return 1000
}

// HeartbeatTimeout gets the timeout which is between 150 to 300 ms
func HeartbeatTimeout() int {
	rand.Seed(time.Now().UnixNano())
	recommendedMax := 300
	recommendedMin := 150
	return rand.Intn(recommendedMax-recommendedMin) + recommendedMin
}
