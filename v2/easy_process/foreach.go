package easy_process

func ForeachSlice[T1 any](t1List []T1, f func(t1 T1)) {
	for _, t1 := range t1List {
		f(t1)
	}
}

func ForeachMap[Key1 GKey, V1 any](m1 map[Key1]V1, f func(Key1, V1)) {
	for k1, v1 := range m1 {
		f(k1, v1)
	}
}
