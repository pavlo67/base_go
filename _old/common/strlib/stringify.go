package strlib

func Stringify(values []interface{}) []interface{} {
	var valuesStr []interface{}

	for _, v := range values {
		switch val := v.(type) {
		case []byte:
			valuesStr = append(valuesStr, string(val))
		default:
			valuesStr = append(valuesStr, v)
		}
	}

	return valuesStr
}
