package easy_process

func MergeSlice[T any](t1 []T, t2 []T) (t3 []T) {
	t3 = []T{}
	t3 = append(t3, t1...)
	t3 = append(t3, t2...)
	return t3
}
