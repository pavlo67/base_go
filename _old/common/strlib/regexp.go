package strlib

import "regexp"

var ReSemicolon = regexp.MustCompile(`\s*;\s*`)
var ReSpaces = regexp.MustCompile(`\s+`)
var ReSpacesFin = regexp.MustCompile(`(^\s+|\s+$)`)
var ReDigitsOnly = regexp.MustCompile(`^\d+$`)
var ReCheckUTF8 = regexp.MustCompile(`(?sm)\p{Cyrillic}+`)
