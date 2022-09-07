package constraints

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Complex interface {
	~complex64 | ~complex128
}

type Ordered interface {
	Integer | Float | ~string
}

// Comparable 比较大小接口
type Comparable[E any] interface {
	// Compare 比较大小
	// if v1.Compare(v2) > 0 {
	// 		v1 > v2
	// } else if v1.Compare(v2) < 0 {
	// 		v1 < v2
	// } else {
	// 		v1 == v2
	//}
	Compare(E) int
}
