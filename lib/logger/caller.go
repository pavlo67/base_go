package logger

import (
	"path"
	"runtime"
	"strconv"
	"strings"
)

type CallInfo struct {
	PackageName     string
	PackageFullName string
	FileName        string
	FuncName        string
	Line            string
}

func GetCallInfo() *CallInfo {
	pc, file, line, _ := runtime.Caller(1)

	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)

	packageFullName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageFullName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageFullName = strings.Join(parts[0:pl-1], ".")
	}

	packageFullNameParts := strings.Split(packageFullName, "/")
	packageName := packageFullName
	if len(packageFullNameParts) > 0 {
		packageName = packageFullNameParts[len(packageFullNameParts)-1]
	}

	return &CallInfo{
		PackageName:     packageName,
		PackageFullName: packageFullName,
		FileName:        fileName,
		FuncName:        funcName,
		Line:            strconv.Itoa(line),
	}
}
