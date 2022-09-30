package easy_process

func ReverseSlice[T any](l1 []T) (l2 []T) {
	l2 = []T{}
	for i := len(l1) - 1; i >= 0; i-- {
		l2 = append(l2, l1[i])
	}
	return l2
}
