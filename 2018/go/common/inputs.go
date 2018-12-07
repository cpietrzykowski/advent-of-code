package inputs

import (
	"bufio"
	"os"
	"path"
	"runtime"
)

// ModulePath gets the "run" (caller) command directory
func ModulePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return path.Dir(file)
	}

	return ""
}

// ReadLines is a convenience for reading a file into an array of strings
func ReadLines(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
