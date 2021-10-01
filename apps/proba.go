package main

import (
	"encoding/json"
	"log"

	"github.com/pavlo67/common/common"
)

func main() {
	data := `{"change":"change"}`
	var m common.Map

	log.Print(json.Unmarshal([]byte(data), &m))

}
