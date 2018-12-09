package common

import (
	"fmt"
	"path"
	"runtime"
)

// CommandPath gets the directory that called the command command directory
// used to support "go run" invocations
func CommandPath() (string, error) {
	if _, file, _, ok := runtime.Caller(0); ok {
		return path.Dir(file), nil
	}

	return "", fmt.Errorf("could not determine command path")
}
