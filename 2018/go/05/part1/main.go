package main

import (
	"aoc/common"
	"bufio"
	"log"
	"os"
	"strings"
)

const day = 5
const part = 1

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		polymer := readPolymerFile(file)
		polymer = calcPolymer(polymer)
		log.Println(polymer, len(polymer))
	} else {
		log.Println(error)
	}
}

func calcPolymer(polymer string) string {
	dirty := true
	rslt := polymer
	for dirty {
		buffer := strings.Builder{}
		dirty = false
		var prevr rune = -1
		for _, r := range rslt {
			if prevr > (-1) {
				if ((r + ('a' - 'A')) == prevr) || ((r - ('a' - 'A')) == prevr) {
					// reaction
					dirty = true
					prevr = (-1)
				} else {
					buffer.WriteRune(prevr)
					prevr = r
				}
			} else {
				prevr = r
			}
		}

		if prevr > (-1) {
			buffer.WriteRune(prevr)
		}

		rslt = buffer.String()
	}

	return rslt
}

func readPolymerFile(f *os.File) string {
	rslt := strings.Builder{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rslt.WriteString(scanner.Text())
	}

	return rslt.String()
}
