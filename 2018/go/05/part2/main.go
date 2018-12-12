package main

import (
	"aoc/common"
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

const day = 5
const part = 2

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type polymer string

func main() {
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		p := readPolymerFile(file)
		reactions := map[rune]polymer{}
		shortest, plen := rune(0), 0
		for i := 'A'; i <= 'Z'; i++ {
			newp := p.stripUnits(i).react()
			reactions[i] = newp
			newplen := len(newp)
			if (newplen < plen) || (shortest == 0) {
				shortest = i
				plen = newplen
			}
		}

		log.Println(string(shortest), plen)
	} else {
		log.Println(error)
	}
}

func (p polymer) react() polymer {
	dirty := true
	rslt := string(p)
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

	return polymer(rslt)
}

func (p polymer) stripUnits(u rune) polymer {
	rslt := strings.Builder{}
	upperr := unicode.ToUpper(u)
	for _, r := range p {
		if !(unicode.ToUpper(r) == upperr) {
			rslt.WriteRune(r)
		}
	}

	return polymer(rslt.String())
}

func readPolymerFile(f *os.File) polymer {
	rslt := strings.Builder{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rslt.WriteString(scanner.Text())
	}

	return polymer(rslt.String())
}
