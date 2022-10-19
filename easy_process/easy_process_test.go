package easy_process

import (
	"fmt"
	"git.singularity-ai.com/backend/library/utils"
	"testing"
)

func TestName(t *testing.T) {
	ll := []int64{1, 3, 5, 7, 9, 2, 4, 6, 8, 0}
	newLL := FilterSlice(ll, func(i int64) bool {
		return i%2 == 0
	})
	fmt.Println(newLL)
}

func TestMap(t *testing.T) {
	ll := []string{}
	ll = nil
	llnew := MapSlice(ll, func(t1 string) (t2 int) {
		return len(t1)
	})
	println(utils.MustJson(llnew))
}

func TestMulti(t *testing.T) {
	ll := []int64{1, 3, 5, 7, 9, 2, 4, 6, 8, 0}
	newLL := ChoiceMulti(ll, 4)
	fmt.Println(newLL)
}
