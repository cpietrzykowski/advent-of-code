package main

import (
	"aoc/common"
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	const day = 3
	const part = 1
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		overclaimed := overlappingClaims(file)
		log.Println(overclaimed)
	} else {
		log.Println(error)
	}
}

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

type Claim struct {
	Id       int
	Location Point
	Size     Size
}

type Registry map[int]map[int][]*Claim

func (registry *Registry) Register(claim Claim) {
	for y := claim.Location.Y; y < (claim.Location.Y + claim.Size.Height); y++ {
		for x := claim.Location.X; x < (claim.Location.X + claim.Size.Width); x++ {
			if _, ok := (*registry)[y]; !ok {
				(*registry)[y] = make(map[int][]*Claim)
			}

			if _, ok := (*registry)[y][x]; !ok {
				(*registry)[y][x] = []*Claim{}
			}

			(*registry)[y][x] = append((*registry)[y][x], &claim)
		}
	}
}

func overlappingClaims(file *os.File) int {
	scanner := bufio.NewScanner(file)
	registry := Registry{}
	for scanner.Scan() {
		if claim, error := makeClaim(scanner.Text()); error == nil {
			registry.Register(*claim)
		} else {
			log.Println(error)
		}
	}

	overclaimed := 0
	for _, row := range registry {
		for _, col := range row {
			if len(col) > 1 {
				overclaimed++
			}
		}
	}

	return overclaimed
}

func makeClaim(entry string) (*Claim, error) {
	// #12 @ 269,129: 25x14
	re := regexp.MustCompile(`#(\d+)\s*@\s*(\d+),(\d+):\s+(\d+)x(\d+)`)
	fields := re.FindAllStringSubmatch(entry, -1)

	if (fields != nil) && (len(fields[0]) == 6) {
		if id, error := strconv.Atoi(fields[0][1]); error == nil {
			var point Point
			if xpos, error := strconv.Atoi(fields[0][2]); error == nil {
				if ypos, error := strconv.Atoi(fields[0][3]); error == nil {
					point = Point{xpos, ypos}
				} else {
					return nil, error
				}
			} else {
				return nil, error
			}

			var size Size
			if width, error := strconv.Atoi(fields[0][4]); error == nil {
				if height, error := strconv.Atoi(fields[0][5]); error == nil {
					size = Size{width, height}
				} else {
					return nil, error
				}
			} else {
				return nil, error
			}

			return &Claim{id, point, size}, nil
		} else {
			return nil, error
		}
	}

	return nil, nil
}
