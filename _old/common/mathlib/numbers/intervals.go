package numbers

func Intersect[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](from1, to1, from2, to2 T) bool {
	if from1 < from2 {
		return from2 <= to1
	}
	return from1 <= to2
}

func Intersection[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](from1, to1, from2, to2 T) *[2]T {
	if from1 < from2 {
		if to1 < from2 {
			return nil
		} else if to1 <= to2 {
			return &[2]T{from2, to1}
		}
		return &[2]T{from2, to2}
	}

	// else from2 <= from1

	if to2 < from1 {
		return nil
	} else if to2 <= to1 {
		return &[2]T{from1, to2}
	}
	return &[2]T{from1, to1}
}

func Connection[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](from1, to1, from2, to2 T) *[2]T {
	if from1 < from2 {
		if to1 < from2 {
			return nil
		} else if to1 <= to2 {
			return &[2]T{from1, to2}
		}
		return &[2]T{from1, to1}
	}

	// else from2 <= from1

	if to2 < from1 {
		return nil
	} else if to2 <= to1 {
		return &[2]T{from2, to1}
	}
	return &[2]T{from2, to2}
}
