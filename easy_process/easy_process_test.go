package easy_process

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	ll := []int64{1, 3, 5, 7, 9, 2,4,6,8,0}
	newLL := FilterSlice(ll, func(i int64) bool {
		return  i % 2 == 0
	})
	fmt.Println(newLL)
}