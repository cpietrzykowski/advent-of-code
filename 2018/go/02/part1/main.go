package main

import (
	"aoc/common"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	const day = 2
	const part = 1
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		fmt.Println(checksum(file))
	} else {
		log.Println(error)
	}
}

func checksum(file *os.File) int {
	a, b := 0, 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		j, k := idcalc(scanner.Text())
		if j {
			a++
		}

		if k {
			b++
		}
	}

	return (a * b)
}

func idcalc(id string) (bool, bool) {
	a, b := 0, 0
	letters := make(map[rune]int)

	for _, char := range id {
		if _, ok := letters[char]; ok {
			letters[char]++
			if letters[char] == 2 {
				a++
			} else if letters[char] == 3 {
				a--
				b++
			} else if letters[char] == 4 {
				a--
				b--
			}
		} else {
			letters[char] = 1
		}
	}

	return (a > 0), (b > 0)
}
