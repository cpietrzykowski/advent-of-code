package main

import (
	"bufio"
	"errors"
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
		if freq, error := calibrate(file); error == nil {
			fmt.Println(freq)
		}
	} else {
		log.Println(error)
	}
}

func calibrate(file *os.File) (int, error) {
	freqs := make(map[int]int)
	rslt := 0
	ix := 0
	passix := 0
	for {
		passix++
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if number, error := strconv.Atoi(scanner.Text()); error == nil {
				rslt += number
				if _, exists := freqs[rslt]; exists {
					return rslt, nil
				}

				freqs[rslt] = ix
				ix++
			}
		}

		if _, error := file.Seek(0, 0); error != nil {
			break
		}
	}

	return 0, errors.New("calibration error")
}
