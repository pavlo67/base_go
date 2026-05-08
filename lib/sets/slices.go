package sets

func Multiple[T interface{}](data [][]T) (dataMultipled [][]T) {

	for _, optionsSet := range data {
		var optionsMultipledNew [][]T

		for _, optionValue := range optionsSet {
			if len(dataMultipled) < 1 {
				optionsMultipledNew = append(optionsMultipledNew, []T{optionValue})
			} else {
				for _, opts := range dataMultipled {
					optionsMultipledNew = append(optionsMultipledNew, append(opts, optionValue))
				}
			}
		}

		dataMultipled = optionsMultipledNew
	}

	return dataMultipled
}

//func CyclicSignature(data []string) string {
//	if len(data) < 1 {
//		return ""
//	}
//
//	var sign string
//	for i := range data {
//		var signOption string
//		for j := range data {
//			signOption += "_" + data[(i+j)%len(data)]
//		}
//		if i == 0 || signOption < sign {
//			sign = signOption
//		}
//	}
//
//	return sign[1:]
//}

func Intersects[T comparable](set0, set1 []T) bool {

	for _, v0 := range set0 {
		for _, v1 := range set1 {
			if v0 == v1 {
				return true
			}
		}
	}

	return false
}
