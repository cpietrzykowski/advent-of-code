package main

import (
	"fmt"
	inputs "go/common"
	"log"
	"os"
	"path/filepath"
)

func main() {
	inputspath := filepath.Join(inputs.ModulePath(), "../../inputs/data/02/input.txt")
	if file, error := os.Open(inputspath); error == nil {
		defer file.Close()
		if lines := inputs.ReadLines(file); error == nil {
			common := findCommon(lines)
			log.Println(string(common))
		}
	} else {
		log.Println(error)
	}
}

func stringDiff(l string, r string) ([]rune, error) {
	if len(l) == len(r) {
		common := []rune{}

		for i, char := range l {
			rrunes := []rune(r)
			if char == rrunes[i] {
				common = append(common, char)
			}
		}

		return common, nil
	}

	return nil, fmt.Errorf("string inconsistency")
}

func findCommon(lines []string) []rune {
	for i, line := range lines {
		for _, other := range lines[i+1:] {
			if common, error := stringDiff(line, other); error == nil {
				if (len(line) - len(common)) <= 1 {
					return common
				}
			} else {
				log.Println(error)
			}
		}
	}

	return nil
}
