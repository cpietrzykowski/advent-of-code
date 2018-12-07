package inputs

import (
	"path"
	"runtime"
)

func ModulePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return path.Dir(file)
	}

	return ""
}
