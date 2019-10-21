package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	widthPtr := flag.Int("width", 10, "The width of the board")
	heightPtr := flag.Int("height", 10, "The height of the board")
	densityPtr := flag.Int("density", 25, "The percentage of cells to start as alive")
	delayPtr := flag.Int("delay", 1000, "How many milliseconds to wait between generations")
	generationPtr := flag.Int("generations", 100, "How many generations to run")

	flag.Parse()

	width := *widthPtr
	height := *heightPtr
	density := (*densityPtr) * width * height / 100
	delay := time.Duration(*delayPtr)
	generations := *generationPtr

	boardA := initBoard(width, height)
	populateBoard(boardA, width, height, density)
	boardB := initBoard(width, height)
	for generation := 0; generation < generations; generation++ {
		time.Sleep(delay * time.Millisecond)
		if generation%2 == 0 {
			simulateGeneration(boardA, boardB, width, height)
			printBoard(boardB, width, height)
		} else {
			simulateGeneration(boardB, boardA, width, height)
			printBoard(boardA, width, height)
		}
	}
}

func printBoard(board [][]bool, xDim int, yDim int) {
	var buf bytes.Buffer
	for y := 0; y < yDim; y++ {
		for x := 0; x < xDim; x++ {
			if board[y][x] {
				buf.WriteByte('*')
			} else {
				buf.WriteByte(' ')
			}
		}
		buf.WriteByte('\n')
	}
	fmt.Print("\033[H\033[2J", buf.String())
}

func neighborCount(board [][]bool, x int, y int, xDim int, yDim int) int {
	neighbors := 0
	for yChange := -1; yChange <= 1; yChange++ {
		for xChange := -1; xChange <= 1; xChange++ {
			if xChange != 0 || yChange != 0 {
				if board[makeInBounds(y+yChange, yDim)][makeInBounds(x+xChange, xDim)] {
					neighbors++
				}
			}
		}
	}
	return neighbors
}

func determineNextState(board [][]bool, x int, y int, xDim int, yDim int) bool {
	neighbors := neighborCount(board, x, y, xDim, yDim)
	if board[y][x] {
		return neighbors >= 2 && neighbors <= 3
	}
	return neighbors == 3
}

func initBoard(xDim, yDim int) (board [][]bool) {
	board = make([][]bool, yDim)
	for i := 0; i < yDim; i++ {
		board[i] = make([]bool, xDim)
	}
	return
}

func populateBoard(board [][]bool, xDim int, yDim int, population int) {
	for population > 0 {
		x := rand.Intn(xDim)
		y := rand.Intn(yDim)
		if !board[y][x] {
			board[y][x] = true
			population--
		}
	}
}

func simulateGeneration(prevBoard [][]bool, nextBoard [][]bool, xDim int, yDim int) {
	for x := 0; x < xDim; x++ {
		for y := 0; y < yDim; y++ {
			nextBoard[y][x] = determineNextState(prevBoard, x, y, xDim, yDim)
		}
	}
}

func makeInBounds(val, max int) int {
	if val < 0 {
		return max - 1
	} else if val >= max {
		return 0
	}
	return val
}
