package main

import (
	"aoc/common"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	const day = 3
	const part = 2
	if file, error := common.AOCInputFile(day); error == nil {
		defer file.Close()
		validclaims := validClaims(file)
		log.Println(validclaims)
	} else {
		log.Println(error)
	}
}

type Point struct {
	X int
	Y int
}

func (point Point) String() string {
	return fmt.Sprintf("%d,%d", point.X, point.Y)
}

type Size struct {
	Width  int
	Height int
}

func (size Size) String() string {
	return fmt.Sprintf("%dx%d", size.Width, size.Height)
}

type Claim struct {
	Id       int
	Location Point
	Size     Size
}

func (claim Claim) String() string {
	return fmt.Sprintf("%06d: (%s) %s", claim.Id, claim.Location, claim.Size)
}

type Registry struct {
	Claims []*Claim
	Mapped map[int]map[int][]*Claim
}

func NewRegistry() *Registry {
	r := new(Registry)
	r.Claims = []*Claim{}
	r.Mapped = map[int]map[int][]*Claim{}
	return r
}

func (registry *Registry) Register(claim *Claim) {
	registry.Claims = append(registry.Claims, claim)
	for y := claim.Location.Y; y < (claim.Location.Y + claim.Size.Height); y++ {
		for x := claim.Location.X; x < (claim.Location.X + claim.Size.Width); x++ {
			if _, ok := registry.Mapped[y]; !ok {
				registry.Mapped[y] = make(map[int][]*Claim)
			}

			if _, ok := registry.Mapped[y][x]; !ok {
				registry.Mapped[y][x] = []*Claim{}
			}

			registry.Mapped[y][x] = append(registry.Mapped[y][x], claim)
		}
	}
}

func (registry *Registry) ValidateClaim(claim *Claim) bool {
	for y := claim.Location.Y; y < (claim.Location.Y + claim.Size.Height); y++ {
		for x := claim.Location.X; x < (claim.Location.X + claim.Size.Width); x++ {
			if len(registry.Mapped[y][x]) > 1 {
				return false
			}
		}
	}

	return true
}

func validClaims(file *os.File) []Claim {
	scanner := bufio.NewScanner(file)
	registry := NewRegistry()
	for scanner.Scan() {
		if claim, error := makeClaim(scanner.Text()); error == nil {
			registry.Register(claim)
		} else {
			log.Println(error)
		}
	}

	validclaims := []Claim{}
	for _, claim := range registry.Claims {
		if registry.ValidateClaim(claim) {
			validclaims = append(validclaims, *claim)
		}
	}

	return validclaims
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
