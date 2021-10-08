package crud

import (
	"fmt"
	"strconv"
)

var _ fmt.Stringer = ID("")

type ID string

func (id ID) String() string {
	return string(id)
}

//type ID fmt.Stringer

func NewID(idStr string) ID {
	return ID(idStr)
}

func NewIDInt64(idInt64 int64) ID {
	return ID(strconv.FormatInt(idInt64, 10))
}

//var _ ID = id("")
//
//type id string

//func (idStr id) String() string {
//	return string(idStr)
//}
