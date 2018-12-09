package main

import (
	"aoc/common"
	"fmt"
	"log"
)

func main() {
	const day = 2
	const part = 2
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		if lines := common.ReadLines(file); error == nil {
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
