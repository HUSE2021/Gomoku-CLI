package main

import (
	"fmt"
)


type Board struct {
	
}

func (b *Board) InitialBoard() {
	
}

func (b *Board) putPiece(x, y, userType int) int {
	if b.tokens[x*15+y] == 0 {
		b.tokens[x*15+y] = userType
		return 200 //200 is ok, 500 is not ok
	} else {
		return 500
	}
}

func (b *Board) returnPieceTypeByPosition(x, y int) int {
	if b.tokens[x*15+y] != 0 {
		if b.tokens[x*15+y] == 1 {
			return 1
		} else {
			return 2
		}
	} else {
		return 0
	}
}

func main() {
	
}
