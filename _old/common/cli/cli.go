package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pavlo67/common/common/mathlib/sets"
)

func Confirm(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question + "? ")

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)

	return sets.In([]string{"y", "yes"}, strings.ToLower(strings.TrimSpace(text)))
}
