package cli

import "runtime"

var Reset = "\033[0m"
var Bold = "\033[1m"
var Underline = "\033[4m"
var Strike = "\033[9m"
var Italic = "\033[3m"

var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Bold = ""
		Underline = ""
		Strike = ""
		Italic = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}
