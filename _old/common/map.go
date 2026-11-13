package common

import (
	"reflect"
	"strconv"
)

type Map map[string]interface{}

func (p Map) StringDefault(key, defaultStr string) string {
	if reflect.TypeOf(p[key]) != nil && reflect.TypeOf(p[key]).Kind() == reflect.String {
		return reflect.ValueOf(p[key]).String()
	}

	switch value := p[key].(type) {
	case []byte:
		return string(value)
	case *[]byte:
		return string(*value)
	case *string:
		return *value
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', 1, 64)
	case float64:
		return strconv.FormatFloat(value, 'f', 1, 64)
	}

	return defaultStr
}

func (p Map) String(key string) (string, bool) {
	if reflect.TypeOf(p[key]) != nil && reflect.TypeOf(p[key]).Kind() == reflect.String {
		return reflect.ValueOf(p[key]).String(), true
	}

	switch value := p[key].(type) {
	//case string:
	//	return value, true
	case []byte:
		return string(value), true
	case *[]byte:
		return string(*value), true
	case *string:
		return *value, true
	case int:
		return strconv.Itoa(value), true
	case int64:
		return strconv.FormatInt(value, 10), true
	case int32:
		return strconv.FormatInt(int64(value), 10), true
	case int16:
		return strconv.FormatInt(int64(value), 10), true
	case int8:
		return strconv.FormatInt(int64(value), 10), true
	case uint:
		return strconv.FormatUint(uint64(value), 10), true
	case uint64:
		return strconv.FormatUint(value, 10), true
	case uint32:
		return strconv.FormatUint(uint64(value), 10), true
	case uint16:
		return strconv.FormatUint(uint64(value), 10), true
	case uint8:
		return strconv.FormatUint(uint64(value), 10), true
	case float32:
		return strconv.FormatFloat(float64(value), 'f', 1, 64), true
	case float64:
		return strconv.FormatFloat(value, 'f', 1, 64), true
	}

	return "", false
}

func (p Map) IsTrue(key string) bool {
	if reflect.TypeOf(p[key]) != nil && reflect.TypeOf(p[key]).Kind() == reflect.Bool {
		return reflect.ValueOf(p[key]).Bool()
	}

	switch value := p[key].(type) {
	case string:
		return value != ""
	case *string:
		return *value != ""
	case int:
		return value != 0
	case int64:
		return value != 0
	case int32:
		return value != 0
	case int16:
		return value != 0
	case int8:
		return value != 0
	case uint:
		return value != 0
	case uint64:
		return value != 0
	case uint32:
		return value != 0
	case uint16:
		return value != 0
	case uint8:
		return value != 0
	case float32:
		return value != 0
	case float64:
		return value != 0
	}

	return false
}

func (p Map) Float(key string) (float64, bool) {
	switch value := p[key].(type) {
	case string:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f, true
		}
		return 0, false
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		// TODO!!! check overflow
		return float64(value), true
	case int64:
		// TODO!!! check overflow
		return float64(value), true
	case int32:
		return float64(value), true
	case int16:
		return float64(value), true
	case int8:
		return float64(value), true
	case uint:
		// TODO!!! check overflow
		return float64(value), true
	case uint64:
		// TODO!!! check overflow
		return float64(value), true
	case uint32:
		return float64(value), true
	case uint16:
		return float64(value), true
	case uint8:
		return float64(value), true
	}

	return 0, false
}

func (p Map) Int64(key string) (int64, bool) {
	val := Int64(p[key])
	if val == nil {
		return 0, false
	}

	return *val, true

	//switch value := p[key].(type) {
	//case string:
	//	if i, err := strconv.Atoi(value); err == nil {
	//		return i, true
	//	}
	//	return 0, false
	//case int:
	//	return value, true
	//case int64:
	//	// TODO!!! check overflow
	//	return int(value), true
	//case int32:
	//	return int(value), true
	//case int16:
	//	return int(value), true
	//case int8:
	//	return int(value), true
	//case uint:
	//	return int(value), true
	//case uint64:
	//	// TODO!!! check overflow
	//	return int(value), true
	//case uint32:
	//	return int(value), true
	//case uint16:
	//	return int(value), true
	//case uint8:
	//	return int(value), true
	//case float32:
	//	return int(value), true
	//case float64:
	//	return int(value), true
	//}
	//return 0, false
}

func (p Map) Int64Default(key string, def int64) int64 {
	val, ok := p.Int64(key)
	if !ok {
		return def
	}

	return val
}

//func (p Map) Strings(key string) []string {
//	switch value := p[key].(type) {
//	case string:
//		return []string{value}
//	case []string:
//		return value
//	case int:
//		return []string{strconv.Itoa(value)}
//	case int64:
//		return []string{strconv.FormatInt(value, 10)}
//	case int32:
//		return []string{strconv.FormatInt(int64(value), 10)}
//	case int16:
//		return []string{strconv.FormatInt(int64(value), 10)}
//	case int8:
//		return []string{strconv.FormatInt(int64(value), 10)}
//	case uint:
//		return []string{strconv.FormatUint(uint64(value), 10)}
//	case uint64:
//		return []string{strconv.FormatUint(value, 10)}
//	case uint32:
//		return []string{strconv.FormatUint(uint64(value), 10)}
//	case uint16:
//		return []string{strconv.FormatUint(uint64(value), 10)}
//	case uint8:
//		return []string{strconv.FormatUint(uint64(value), 10)}
//	case float32:
//		return []string{strconv.FormatFloat(float64(value), 'f', 1, 64)}
//	case float64:
//		return []string{strconv.FormatFloat(value, 'f', 1, 64)}
//	}
//
//	return nil
//}

//func (p Map) Value(key string, target interface{}) error {
//	if value, ok := p[key]; ok {
//		valueRaw, err := c.marshaler.Marshal(value)
//		if err != nil {
//			return errors.Wrapf(err, "can't marshal value (%s / %#v) to raw bytes", key, value)
//		}
//
//		return c.marshaler.Unmarshal(valueRaw, target)
//	}
//
//	return nil
//}

func (p Map) Int(key string) (int, bool) {
	switch value := p[key].(type) {
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			return i, true
		}
		return 0, false
	case int:
		return value, true
	case int64:
		// TODO!!! check overflow
		return int(value), true
	case int32:
		return int(value), true
	case int16:
		return int(value), true
	case int8:
		return int(value), true
	case uint:
		return int(value), true
	case uint64:
		// TODO!!! check overflow
		return int(value), true
	case uint32:
		return int(value), true
	case uint16:
		return int(value), true
	case uint8:
		return int(value), true
	case float32:
		return int(value), true
	case float64:
		return int(value), true
	}

	return 0, false
}

func (p Map) Strings(key string) []string {
	switch value := p[key].(type) {
	case string:
		return []string{value}
	case []string:
		return value
	case int:
		return []string{strconv.Itoa(value)}
	case int64:
		return []string{strconv.FormatInt(value, 10)}
	case int32:
		return []string{strconv.FormatInt(int64(value), 10)}
	case int16:
		return []string{strconv.FormatInt(int64(value), 10)}
	case int8:
		return []string{strconv.FormatInt(int64(value), 10)}
	case uint:
		return []string{strconv.FormatUint(uint64(value), 10)}
	case uint64:
		return []string{strconv.FormatUint(value, 10)}
	case uint32:
		return []string{strconv.FormatUint(uint64(value), 10)}
	case uint16:
		return []string{strconv.FormatUint(uint64(value), 10)}
	case uint8:
		return []string{strconv.FormatUint(uint64(value), 10)}
	case float32:
		return []string{strconv.FormatFloat(float64(value), 'f', 1, 64)}
	case float64:
		return []string{strconv.FormatFloat(value, 'f', 1, 64)}
	}

	return nil
}
