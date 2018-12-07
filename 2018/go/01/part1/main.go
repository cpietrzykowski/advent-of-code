package main

import (
	"bufio"
	"fmt"
	inputs "go/common"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	inputspath := filepath.Join(inputs.ModulePath(), "../../inputs/data/01/input.txt")
	if file, error := os.Open(inputspath); error == nil {
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
