package main

import (
	"aoc/common"
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if problem, error := common.GetAOCProblem(); error == nil {
		doProblem(problem)
	} else {
		log.Println(error)
	}
}

// using floats due to unimplemented int utilities (min/max/abs,etc.)
// might be expensive? if anything it complicates relatively simple int statements
type point struct {
	X float64
	Y float64
}

func (p point) String() string {
	return fmt.Sprintf("(% 4.1f, % 4.1f)", p.X, p.Y)
}

type size struct {
	Width  float64
	Height float64
}

type bounds struct {
	origin point
	size   size
}

func (p point) taxicabDistanceTo(other point) float64 {
	xd := math.Abs(p.X - other.X)
	yd := math.Abs(p.Y - other.Y)
	return xd + yd
}

func newPointFromComponents(components []string) (*point, error) {
	x, error := strconv.Atoi(strings.TrimSpace(components[0]))
	if error != nil {
		return nil, error
	}

	y, error := strconv.Atoi(strings.TrimSpace(components[1]))
	if error != nil {
		return nil, error
	}

	return &point{float64(x), float64(y)}, nil
}

type distanceMap struct {
	pois   []point
	grid   [][]point
	bounds bounds
}

func (m distanceMap) nearestPoi(p point) point {
	nearest, d := point{-1, -1}, (m.bounds.size.Width + m.bounds.size.Height)
	for _, n := range m.pois {
		if n == p {
			return n
		}

		ndist := p.taxicabDistanceTo(n)
		if ndist < d {
			d = ndist
			nearest = n
		} else if ndist == d {
			nearest = point{-1, -1}
		}

	}

	return nearest
}

func (m *distanceMap) addPointOfInterest(p point) {
	m.pois = append(m.pois, p)

	if len(m.pois) > 1 {
		// update bounds
		if p.X < m.bounds.origin.X {
			m.bounds.origin.X = p.X
		}

		if p.X > (m.bounds.origin.X + m.bounds.size.Width) {
			m.bounds.size.Width = p.X - m.bounds.origin.X + 1
		}

		if p.Y < m.bounds.origin.Y {
			m.bounds.origin.Y = p.Y
		}

		if p.Y > (m.bounds.origin.Y + m.bounds.size.Height) {
			m.bounds.size.Height = p.Y - m.bounds.origin.Y + 1
		}
	} else {
		m.bounds = bounds{point{p.X, p.Y}, size{1, 1}}
	}
}

func (m *distanceMap) calculateNearestNeighbors() {
	m.grid = make([][]point, int(m.bounds.size.Height))
	for y := range m.grid {
		m.grid[y] = make([]point, int(m.bounds.size.Width))
		for x := range m.grid[y] {
			m.grid[y][x] = point{-1, -1}
			thisp := point{float64(x) + m.bounds.origin.X, float64(y) + m.bounds.origin.Y}
			nearestp := m.nearestPoi(thisp)
			m.grid[y][x] = nearestp
		}
	}
}

func (m distanceMap) edgePois() map[point]bool {
	edgepois := map[point]bool{}
	for y := range m.grid {
		for x := range m.grid[y] {
			poi := m.grid[y][x]
			if _, ok := edgepois[poi]; !ok {
				if x == 0 || y == 0 || x == len(m.grid[y])-1 || y == len(m.grid)-1 {
					edgepois[poi] = true
				}
			}
		}
	}

	return edgepois
}

func (m distanceMap) mostIsolated() (point, int) {
	m.calculateNearestNeighbors()
	counts := map[point]int{}
	mx, mxpoi := 0, point{-1, -1}
	edgeset := m.edgePois()
	for y := range m.grid {
		for x := range m.grid[y] {
			if _, ok := edgeset[m.grid[y][x]]; ok {
				continue
			}

			_, ok := counts[m.grid[y][x]]
			if !ok {
				counts[m.grid[y][x]] = 0
			}

			counts[m.grid[y][x]]++
			if counts[m.grid[y][x]] > mx {
				mx = counts[m.grid[y][x]]
				mxpoi = m.grid[y][x]
			}
		}
	}

	return mxpoi, mx
}

func doProblem(problem *common.AOCProblem) {
	log.Printf("%+v", problem)
	if file, error := problem.OpenInput(); error == nil {
		defer file.Close()
		coords := coordinatesFromFile(file)
		distmap := distanceMap{}
		for _, c := range coords {
			distmap.addPointOfInterest(c)
		}

		log.Println(distmap.mostIsolated())
	}
}

func coordinatesFromFile(file *os.File) []point {
	coordinates := []point{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// each line is a coordinate pair
		// example: xxx, yyy
		line := scanner.Text()
		coord, error := newPointFromComponents(strings.Split(line, ","))

		if error != nil {
			log.Println("unable to parse line", error)
			continue
		}

		coordinates = append(coordinates, *coord)
	}

	return coordinates
}
