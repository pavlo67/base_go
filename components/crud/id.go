package crud

import (
	"fmt"
	"strconv"
)

type ID fmt.Stringer

func NewID(idStr string) ID {
	return id(idStr)
}

func NewIDInt64(idInt64 int64) ID {
	return id(strconv.FormatInt(idInt64, 10))
}

var _ ID = id("")

type id string

func (idStr id) String() string {
	return string(idStr)
}
