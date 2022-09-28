package easy_process

import (
	"math/rand"
	"time"
)

func ChoiceOne[T any](choices []T) (choice T) {
	rand.Seed(time.Now().UnixMilli())
	return choices[rand.Intn(len(choices))]
}

func ChoiceOneFromMap[T any, K string | int | int64 | int32 | int8 | float64 | float32 | uint | uint8 | uint16 | uint32 | uint64 | uintptr](choices map[K]T) (choice T) {
	for _, v := range choices {
		return v
	}
	panic("no item")
}
