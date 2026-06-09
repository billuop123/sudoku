package main

import (
	"math/rand"
	"slices"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SIZE     int = 9
	CELLSIZE int = 60
	FPS          = 60
)

func main() {
	grid := make([][]string, SIZE)
	for i := range SIZE {
		grid[i] = make([]string, SIZE)
		for j := range SIZE {
			grid[i][j] = "."
		}
	}
	populateGrid(grid)
	for !isValid(grid) {
		populateGrid(grid)
	}
	currentNumber := -1
	gameOver := false
	won := false
	options := make([][]string, 3)
	rl.InitWindow(900, 600, "Sudoku")
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		drawGrid(grid)
		drawAndPopulateOptions(options, strconv.Itoa(currentNumber))
		if gameOver {
			offsetX := 600
			offsetY := 30
			posX := offsetX
			posY := offsetY + 60*3 + offsetY
			rl.DrawRectangle(int32(posX), int32(posY), int32(CELLSIZE)+170, int32(CELLSIZE), rl.Red)
			rl.DrawText("Game Over", int32(posX+10), int32(posY+10), 40, rl.White)
			rl.DrawRectangle(int32(posX), int32(posY)+int32(CELLSIZE)+int32(offsetY), int32(CELLSIZE)+170, int32(CELLSIZE), rl.Red)
			rl.DrawText("New Game", int32(posX+10), int32(posY)+int32(CELLSIZE)+int32(offsetY)+10, 40, rl.White)
		}
		if won {
			offsetX := 600
			offsetY := 30
			posX := offsetX
			posY := offsetY + 60*3 + offsetY
			rl.DrawRectangle(int32(posX), int32(posY), int32(CELLSIZE)+170, int32(CELLSIZE), rl.Green)
			rl.DrawText("You Won!", int32(posX+10), int32(posY+10), 40, rl.White)
			rl.DrawRectangle(int32(posX), int32(posY)+int32(CELLSIZE)+int32(offsetY), int32(CELLSIZE)+170, int32(CELLSIZE), rl.Red)
			rl.DrawText("New Game", int32(posX+10), int32(posY)+int32(CELLSIZE)+int32(offsetY)+10, 40, rl.White)
		}
		rl.EndDrawing()
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			offset := 30
			optionsOffset := 600
			pos := rl.GetMousePosition()
			posX := pos.X
			posY := pos.Y
			// for playgrid
			if posX >= float32(offset) && posX <= float32(CELLSIZE*SIZE+offset) && (posY >= float32(offset) && posY <= float32(CELLSIZE*9+offset)) && (!gameOver) && (!won) {
				col := int32((posX - float32(offset)) / float32(CELLSIZE))
				row := int32((posY - float32(offset)) / float32(CELLSIZE))
				if grid[row][col] == "." && currentNumber != -1 {
					grid[row][col] = strconv.Itoa(currentNumber)
					if !isValid(grid) {
						gameOver = true
					} else {
						if youWon(grid) {
							won = true
						}
					}
				}
			}
			// for options
			if posX >= float32(optionsOffset) && (posX <= float32((optionsOffset + 3*CELLSIZE))) && (posY >= float32(offset) && posY <= float32(CELLSIZE*3+offset)) && (!gameOver) && (!won) {
				col := int32((posX - float32(optionsOffset)) / float32(CELLSIZE))
				row := int32((posY - float32(offset)) / float32(CELLSIZE))
				currentNumber, _ = strconv.Atoi(options[row][col])
			}
			// for newGame
			if posX >= float32(optionsOffset) && posX <= float32(optionsOffset+170) && posY >= float32(offset)+float32(CELLSIZE)*3+float32(offset)+float32(CELLSIZE)+float32(offset) && posY <= float32(offset)+float32(CELLSIZE)*3+float32(offset)+float32(CELLSIZE)+float32(offset)+float32(CELLSIZE) {
				resetGame(grid, &gameOver, &won, &currentNumber)
			}
		}
	}
}

func resetGame(grid [][]string, gameOver *bool, won *bool, currentNumber *int) {
	populateGrid(grid)
	*gameOver = false
	*won = false
	*currentNumber = -1
	drawGrid(grid)
}

func drawAndPopulateOptions(options [][]string, currentNumber string) {
	offsetX := 600
	offsetY := 30
	for i := range 3 {
		options[i] = make([]string, 3)
		for j := range 3 {
			num := i*3 + j + 1
			options[i][j] = strconv.Itoa(num)
			posX := offsetX + j*CELLSIZE
			posY := offsetY + i*CELLSIZE
			if options[i][j] == currentNumber {
				rl.DrawRectangle(int32(posX), int32(posY), int32(CELLSIZE), int32(CELLSIZE), rl.Green)
				rl.DrawRectangleLines(int32(posX), int32(posY), int32(CELLSIZE), int32(CELLSIZE), rl.Black)
				rl.DrawText(strconv.Itoa(num), int32(posX)+20, int32(posY)+15, 40, rl.Black)
			} else {
				rl.DrawRectangle(int32(posX), int32(posY), int32(CELLSIZE), int32(CELLSIZE), rl.LightGray)
				rl.DrawRectangleLines(int32(posX), int32(posY), int32(CELLSIZE), int32(CELLSIZE), rl.Black)
				rl.DrawText(strconv.Itoa(num), int32(posX)+20, int32(posY)+15, 40, rl.Black)
			}
		}
	}
}

func drawGrid(grid [][]string) {
	offset := 30
	for i := range SIZE {
		for j := range SIZE {
			posX := offset + j*CELLSIZE
			posY := offset + i*CELLSIZE
			rl.DrawRectangle(int32(posX), int32(posY), int32(CELLSIZE), int32(CELLSIZE), rl.LightGray)
			rl.DrawRectangleLines(int32(posX), int32(posY), int32(CELLSIZE), int32(CELLSIZE), rl.Black)
			if grid[i][j] != "." {
				rl.DrawText(grid[i][j], int32(posX)+20, int32(posY)+15, 40, rl.Black)
			}
		}
	}
}

func isValid(grid [][]string) bool {
	colExists := map[int][]string{}
	for i := range SIZE {
		rowExists := map[string]bool{}
		for j := range SIZE {
			if grid[i][j] == "." {
				continue
			}
			if rowExists[grid[i][j]] {
				return false
			}
			if slices.Contains(colExists[j], grid[i][j]) {
				return false
			}
			rowExists[grid[i][j]] = true
			colExists[j] = append(colExists[j], grid[i][j])
		}
	}
	for boxRow := range 3 {
		for boxCol := range 3 {
			exists := map[string]bool{}
			for i := range 3 {
				for j := range 3 {
					r := boxRow*3 + j
					c := boxCol*3 + i
					if grid[r][c] == "." {
						continue
					}
					if exists[grid[r][c]] {
						return false
					}
					exists[grid[r][c]] = true
				}
			}
		}
	}
	return true
}

func youWon(grid [][]string) bool {
	for i := range SIZE {
		for j := range SIZE {
			if grid[i][j] == "." {
				return false
			}
		}
	}
	return true
}

func populateGrid(grid [][]string) {
	for i := range SIZE {
		grid[i] = make([]string, SIZE)
		for j := range SIZE {
			grid[i][j] = "."
		}
	}
	for i := range SIZE {
		for j := range SIZE {
			if grid[i][j] == "." {
				randomNumber := rand.Intn(SIZE) + 1
				grid[i][j] = strconv.Itoa(randomNumber)
				if !isValid(grid) {
					grid[i][j] = "."
				}
			}
		}
	}
}
