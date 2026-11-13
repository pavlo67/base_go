package sqllib_pg

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pavlo67/common/common/sqllib"
)

var _ sqllib.CorrectWildcards = CorrectWildcards

var reQuestionMark = regexp.MustCompile("\\?")

func CorrectWildcards(query string) string {
	n := 1

LOOP:
	loc := reQuestionMark.FindStringIndex(query)
	if len(loc) > 0 {
		query = query[:loc[0]] + "$" + strconv.Itoa(n) + query[loc[1]:]
		n++

		goto LOOP
	}

	return query
}

func WildcardsForUpdate(fields []string) string {
	var wildcards []string

	for n, f := range fields {
		wildcards = append(wildcards, f+"=$"+strconv.Itoa(n+1))
	}

	return strings.Join(wildcards, ",")
}

func WildcardsForInsert(fields []string) string {
	var wildcards []string
	for n := 1; n <= len(fields); n++ {
		wildcards = append(wildcards, "$"+strconv.Itoa(n))
	}

	return strings.Join(wildcards, ",")
}
