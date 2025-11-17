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

func ChoiceMulti[T any](choices []T, cnt int) (multi []T) {
	if cnt > len(choices) {
		cnt = len(choices)
	}
	idxMap := map[int]bool{}
	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < cnt; {
		idx := rand.Intn(len(choices))
		_, ok := idxMap[idx]
		if ok {
			continue
		}
		idxMap[idx] = true
		multi = append(multi, choices[idx])
		i++
	}
	return multi
}
