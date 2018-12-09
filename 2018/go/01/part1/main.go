package main

import (
	"aoc/common"
	"bufio"
	"fmt"
	"log"
	"strconv"
)

func main() {
	const day = 1
	const part = 1
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)

		rslt := 0

		for scanner.Scan() {
			if number, error := strconv.Atoi(scanner.Text()); error == nil {
				rslt += number
			}
		}

		fmt.Println(rslt)
	} else {
		log.Println(error)
	}
}
