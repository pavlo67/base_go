package sets

func In[T comparable](ints []T, i T) bool {
	for _, s := range ints {
		if s == i {
			return true
		}
	}

	return false
}

//type Number interface {
//	byte | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32| float64
//}
//
//type Range[T Number] [2]T
//
//func (r Range[T]) Intersects (r1 Range[T]) Range[T] {
//	return Range[T]{}
//}
