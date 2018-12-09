package common

import (
	"bufio"
	"os"
)

// ReadLines is a convenience for reading a file into an array of strings
func ReadLines(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
