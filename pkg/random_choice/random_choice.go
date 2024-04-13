package randomchoice

import (
	"math/rand"
	"time"
)

// RandomChoiceFromSlice returns a random element from a slice of values
func RandomChoiceFromSlice(values []interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(values))
	return values[index]
}
