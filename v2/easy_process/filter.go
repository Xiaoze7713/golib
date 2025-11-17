package easy_process

func FilterSlice[T1 any](eleList []T1, condFunc func(T1) bool) (newEleList []T1) {
	for _, v := range eleList {
		if condFunc(v) {
			newEleList = append(newEleList, v)
		}
	}
	return newEleList
}

type GKey interface {
	~float32 | ~float64 | ~string | ~int64 | ~int32 | ~int8 | ~int16 | ~int
}

func FilterMap[Key GKey, Value any](kvs map[Key]Value, condFunc func(Key, Value) bool) (kvs2 map[Key]Value) {
	kvs2 = map[Key]Value{}
	for k, v := range kvs {
		if condFunc(k, v) {
			kvs2[k] = v
		}
	}
	return
}

func MapSlice[T1 any, T2 any](t1List []T1, conv func(t1 T1) (t2 T2)) (t2List []T2) {
	t2List = []T2{}
	for _, t1 := range t1List {
		t2 := conv(t1)
		t2List = append(t2List, t2)
	}
	return t2List
}
