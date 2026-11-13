package sets

func Union[T comparable](s0, s1 []T) []T {
	for _, s := range s1 {
		if !In(s0, s) {
			s0 = append(s0, s)
		}
	}
	return s0
}

func Intersect[T comparable](s0, s1 []T) bool {
	for _, s := range s1 {
		if In(s0, s) {
			return true
		}
	}
	return false
}

func Equal[T comparable](s0, s1 []T) bool {
	for _, s := range s1 {
		if !In(s0, s) {
			return false
		}
	}
	for _, s := range s0 {
		if !In(s1, s) {
			return false
		}
	}

	return true
}
