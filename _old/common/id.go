package common

import (
	"strconv"
	"strings"
)

// numerical IDs ..............................................................

type IDNum uint64

func (id IDNum) Key() IDStr {
	return IDStr(strconv.FormatUint(uint64(id), 10))
}

// string Keys/IDs ............................................................

type IDStr string

func (id IDStr) Normalize() IDStr {
	return IDStr(strings.TrimSpace(string(id)))
}

func (id IDStr) Uint64() uint64 {
	idUint64, _ := strconv.ParseUint(string(id), 10, 64)

	return idUint64
}
