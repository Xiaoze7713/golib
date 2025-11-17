package easy_process

func MapMap[Key1 GKey, Key2 GKey, V1 any, V2 any](m1 map[Key1]V1, conv func(Key1, V1) (Key2, V2)) (m2 map[Key2]V2) {
	m2 = map[Key2]V2{}
	for k1, v1 := range m1 {
		k2, v2 := conv(k1, v1)
		m2[k2] = v2
	}
	return m2
}

type ComparableValue interface {
	~float32 | ~float64 | ~string | ~int64 | ~int32 | ~int8 | ~int16 | ~int
}

func Choice[T ComparableValue](t1, t2 T, choice func(t1 T, t2 T) (resT T)) (t T) {
	return choice(t1, t2)
}
