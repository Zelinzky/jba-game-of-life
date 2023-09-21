package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type universe struct {
	cells      [][]bool
	height     int
	width      int
	liveCells  int
	generation int
}

func main() {
	var size int
	_, err := fmt.Scan(&size)
	if err != nil {
		log.Fatal(err)
	}
	u := newRandUniverse(size)
	u.PrintUniverse()
	for {
		time.Sleep(500 * time.Millisecond)
		u.nextGen()
		u.PrintUniverse()
	}
}

func newRandUniverse(size int) universe {
	return newUniverse(size, initRandCell)
}

func newEmptyUniverse(size int) universe {
	return newUniverse(size, func() bool {
		return false
	})
}

func newUniverse(size int, initFunc func() bool) universe {
	uni := make([][]bool, size)
	liveCells := 0
	for i := range uni {
		uni[i] = make([]bool, size)
		for j := range uni[i] {
			cellState := initFunc()
			uni[i][j] = cellState
			if cellState {
				liveCells++
			}
		}
	}
	return universe{cells: uni, width: size, height: size, liveCells: liveCells, generation: 1}
}

func (u *universe) PrintUniverse() {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("Generation #%v\nAlive: %v\n", u.generation, u.liveCells)
	for i, row := range u.cells {
		for j := range row {
			if u.cells[i][j] {
				fmt.Print("O")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func (u *universe) normalizeCoordinates(x, y int) (int, int) {
	x = (x + u.width) % u.width
	y = (y + u.height) % u.height
	return x, y
}

func (u *universe) isCellAlive(x, y int) bool {
	x, y = u.normalizeCoordinates(x, y)
	return u.cells[x][y]
}

func (u *universe) setCellAlive(x, y int, alive bool) {
	x, y = u.normalizeCoordinates(x, y)
	u.cells[x][y] = alive
}

func (u *universe) countLivingNeighbours(x, y int) int {
	aliveNeighbours := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && u.isCellAlive(x+i, y+j) {
				aliveNeighbours++
			}
		}
	}
	return aliveNeighbours
}

func (u *universe) nextCellLiveness(x, y int) bool {
	isAlive := u.isCellAlive(x, y)
	livingNeighbours := u.countLivingNeighbours(x, y)
	if livingNeighbours == 3 || (livingNeighbours == 2 && isAlive) {
		return true
	}
	return false
}

func (u *universe) nextGen() {
	nextUni := newEmptyUniverse(u.height)
	liveCells := 0
	for i, row := range u.cells {
		for j := range row {
			nextCellState := u.nextCellLiveness(i, j)
			nextUni.setCellAlive(i, j, nextCellState)
			if nextCellState {
				liveCells++
			}
		}
	}
	nextUni.generation = u.generation + 1
	nextUni.liveCells = liveCells
	*u = nextUni
}

func initRandCell() bool {
	if rand.Intn(2) == 0 {
		return false
	}
	return true
}
