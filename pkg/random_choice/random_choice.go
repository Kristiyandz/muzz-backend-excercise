package randomchoice

import (
	"math/rand"
	"time"
)

func RandomChoiceFromSlice(values []interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(values))
	return values[index]
}
