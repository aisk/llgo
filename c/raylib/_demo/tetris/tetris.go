package main

import (
	"unsafe"

	"github.com/goplus/llgo/c"
	"github.com/goplus/llgo/c/raylib"
)

const (
	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 20
	BLOCK_SIZE   = 30

	SCREENWIDTH  = 300
	SCREENHEIGHT = 600
)

const MAX_BLOCKS = 4

type Shape struct {
	Blocks [MAX_BLOCKS]raylib.Vector2
	Color  raylib.Color
}

var SHAPES = []Shape{
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}}, Color: raylib.SKYBLUE},
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}, Color: raylib.YELLOW},
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}}, Color: raylib.PURPLE},
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}, Color: raylib.GREEN},
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}}, Color: raylib.RED},
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}}, Color: raylib.BLUE},
	{Blocks: [MAX_BLOCKS]raylib.Vector2{{X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}}, Color: raylib.ORANGE},
}

var board [BOARD_HEIGHT][BOARD_WIDTH]raylib.Color

var curShape Shape
var curPos raylib.Vector2

var fallTime = c.Float(0)
var fallSpeed = c.Float(0.2)
var score = 0
var scoreText = make([]c.Char, 20)
var gameOver = false

func genShape() {
	curShape = SHAPES[raylib.GetRandomValue(c.Int(0), c.Int(6))]
	curPos = raylib.Vector2{BOARD_WIDTH/2 - 1, 0}
}

func checkCollision() bool {
	for i := 0; i < MAX_BLOCKS; i++ {
		x := int(curPos.X + curShape.Blocks[i].X)
		y := int(curPos.Y + curShape.Blocks[i].Y)
		if x < 0 || x >= BOARD_WIDTH || y >= BOARD_HEIGHT || (y >= 0 && board[y][x] != raylib.BLANK) {
			return true
		}
	}
	return false
}

func lockShape() {
	for i := 0; i < MAX_BLOCKS; i++ {
		x := int(curPos.X + curShape.Blocks[i].X)
		y := int(curPos.Y + curShape.Blocks[i].Y)
		if y >= 0 {
			board[y][x] = curShape.Color
		}
	}
}

func rotateShape() {
	rotated := curShape
	for i := 0; i < MAX_BLOCKS; i++ {
		x := rotated.Blocks[i].X
		rotated.Blocks[i].X = -rotated.Blocks[i].Y
		rotated.Blocks[i].Y = x
	}

	temp := curShape
	curShape = rotated
	if checkCollision() {
		curShape = temp
	}
}

func clearLines() int {
	linesCleared := 0
	for y := BOARD_HEIGHT - 1; y >= 0; y-- {
		lineFull := true
		for x := 0; x < BOARD_WIDTH; x++ {
			if board[y][x] == raylib.BLANK {
				lineFull = false
				break
			}
		}
		if lineFull {
			for yy := y; yy > 0; yy-- {
				for x := 0; x < BOARD_WIDTH; x++ {
					board[yy][x] = board[yy-1][x]
				}
			}
			for x := 0; x < BOARD_WIDTH; x++ {
				board[0][x] = raylib.BLANK
			}
			y += 1
			linesCleared += 1
		}
	}
	return linesCleared
}

func keyPressed(key c.Int) bool {
	return raylib.IsKeyPressed(key) || raylib.IsKeyPressedRepeat(key)
}

func main() {
	raylib.InitWindow(SCREENWIDTH, SCREENHEIGHT, c.Str("tetris (powered by raylib + llgo)"))
	raylib.SetTargetFPS(c.Int(60))
	genShape()
	for !raylib.WindowShouldClose() && !gameOver {
		fallTime += raylib.GetFrameTime()
		if fallTime >= fallSpeed {
			fallTime = 0
			curPos.Y += 1
			if checkCollision() {
				curPos.Y -= 1
				lockShape()
				linesCleared := clearLines()
				score += linesCleared * 100
				genShape()
				if checkCollision() {
					gameOver = true
				}
			}
		}

		if keyPressed(raylib.KEY_LEFT) {
			curPos.X -= 1
			if checkCollision() {
				curPos.X += 1
			}
		}
		if keyPressed(raylib.KEY_RIGHT) {
			curPos.X += 1
			if checkCollision() {
				curPos.X -= 1
			}
		}
		if keyPressed(raylib.KEY_SPACE) || keyPressed(raylib.KEY_UP) || keyPressed(raylib.KEY_DOWN) {
			rotateShape()
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RAYWHITE)
		for y := 0; y < BOARD_HEIGHT; y++ {
			for x := 0; x < BOARD_WIDTH; x++ {
				raylib.DrawRectangle(c.Int(x*BLOCK_SIZE), c.Int(y*BLOCK_SIZE), c.Int(BLOCK_SIZE-1), c.Int(BLOCK_SIZE-1), board[y][x])
			}
		}

		for i := 0; i < MAX_BLOCKS; i++ {
			raylib.DrawRectangle(c.Int((curPos.X+curShape.Blocks[i].X)*BLOCK_SIZE), c.Int((curPos.Y+curShape.Blocks[i].Y)*BLOCK_SIZE),
				BLOCK_SIZE-1, BLOCK_SIZE-1, curShape.Color)
		}

		c.Sprintf(unsafe.SliceData(scoreText), c.Str("Score:%d"), score)
		raylib.DrawText(unsafe.SliceData(scoreText), 10, 10, 20, raylib.BLACK)

		raylib.EndDrawing()
	}

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RAYWHITE)
		raylib.DrawText(c.Str("Game Over"), SCREENWIDTH/2-50, SCREENHEIGHT/2-10, 20, raylib.RED)
		raylib.DrawText(unsafe.SliceData(scoreText), SCREENWIDTH/2-50, SCREENHEIGHT/2+10, 20, raylib.BLACK)
		raylib.EndDrawing()
	}
}
