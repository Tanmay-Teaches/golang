package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"math/rand"
)

var (
	Black = color.RGBA{0, 0, 0, 255}
	White = color.RGBA{255, 255, 255, 255}
)

//Canvas
type Pixels struct {
	//RGBA colors
	Pix   []uint8
	Width int
}

//Create a new canvas with dimension width x height
func NewPixels(width, height int) *Pixels {
	return &Pixels{Width: width, Pix: make([]uint8, width*height*4)}
}

func (p *Pixels) DrawRect(x, y, width, height int, rgba color.RGBA) {
	for idx := 0; idx < width; idx++ {
		for idy := 0; idy < height; idy++ {
			p.SetColor(x+idx, y+idy, rgba)
		}
	}
}

func (p *Pixels) SetColor(x, y int, rgba color.RGBA) {
	r, g, b, a := rgba.RGBA()
	index := (y*p.Width + x) * 4
	p.Pix[index] = uint8(r)
	p.Pix[index+1] = uint8(g)
	p.Pix[index+2] = uint8(b)
	p.Pix[index+3] = uint8(a)
}

type GameOfLife struct {
	gameBoard [][]int
	pixels    *Pixels
	size      int
}

//Create a new GameOfLife structure with width*height number of cell.
//Size control how big to render the board game
func NewGameOfLife(width, height, size int) *GameOfLife {
	gameBoard := make([][]int, height)
	for idx := range gameBoard {
		gameBoard[idx] = make([]int, width)
	}
	pixels := NewPixels(width*size, height*size)
	return &GameOfLife{gameBoard: gameBoard, pixels: pixels, size: size}
}
func (gol *GameOfLife) Random() {
	for idy := range gol.gameBoard {
		for idx := range gol.gameBoard[idy] {
			gol.gameBoard[idy][idx] = rand.Intn(2)
		}
	}
}
func CountNeighbours(matrix [][]int) [][]int {
	neighbours := make([][]int, len(matrix))
	for idx, val := range matrix {
		neighbours[idx] = make([]int, len(val))
	}
	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			for rowMod := -1; rowMod < 2; rowMod++ {
				newRow := row + rowMod
				if newRow < 0 || newRow >= len(matrix) {
					continue
				}
				for colMod := -1; colMod < 2; colMod++ {
					if rowMod == 0 && colMod == 0 {
						continue
					}
					newCol := col + colMod
					if newCol < 0 || newCol >= len(matrix[row]) {
						continue
					}
					neighbours[row][col] += matrix[newRow][newCol]
				}
			}
		}
	}
	return neighbours
}
func (gol *GameOfLife) PlayRound() {
	neighbours := CountNeighbours(gol.gameBoard)
	for idy := range gol.gameBoard {
		for idx, value := range gol.gameBoard[idy] {
			n := neighbours[idy][idx]
			if value == 1 && (n == 2 || n == 3) {
				continue
			} else if n == 3 {
				gol.gameBoard[idy][idx] = 1
				gol.pixels.DrawRect(idx*gol.size, idy*gol.size, gol.size, gol.size, Black)
			} else {
				gol.gameBoard[idy][idx] = 0
				gol.pixels.DrawRect(idx*gol.size, idy*gol.size, gol.size, gol.size, White)
			}
		}
	}
}

func run() {
	size := float64(2)
	width := float64(400)
	height := float64(400)
	cfg := pixelgl.WindowConfig{
		Title:  "Conway’s Game of Life",
		Bounds: pixel.R(0, 0, width*size, height*size),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	gol := NewGameOfLife(int(width), int(height), int(size))
	gol.Random()
	for !win.Closed() {
		gol.PlayRound()
		win.Canvas().SetPixels(gol.pixels.Pix)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
