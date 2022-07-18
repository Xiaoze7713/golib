package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestDur(t *testing.T) {
	start := time.Now().UnixNano()
	time.Sleep(time.Second)
	tx := UTimeMSDurF(start)
	fmt.Printf("%v\n", tx)
}

func TestFind(t *testing.T) {
	s := "哈哈哈哈 ？。"
	idxList, err := Find("(\\.|。|\\?|？|!|！){1,1}", s)
	if len(idxList) > 0 {
		//return s[]
	}
	fmt.Printf("%v %v", idxList, err)
}