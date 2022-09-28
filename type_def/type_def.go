package type_def

type KeyType interface {
	bool | uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int | float32 | float64
	String()
}

type PtrValueType interface {
	~uintptr
}

type BaseValueType interface {
	bool | uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int | float32 | float64 | string
}
