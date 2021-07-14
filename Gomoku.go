package main

import (
	"fmt"
)

// null 0,  user1:○ = 1　user2: ● = 2
type Board struct {
	tokens [15 * 15]int
}

func (b *Board) InitialBoard() {
	var i int

	for i = 0; i < 15*15; i++ {
		b.tokens[i] = 0
	}
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
	var b Board
	b.InitialBoard()
	b.putPiece(0, 0, 1)
	b.putPiece(0, 1, 2)

	var temp int
	temp = b.returnPieceTypeByPosition(0, 0)
	if temp == 0 {
		fmt.Println(".")
	} else if temp == 1 {
		fmt.Println("○")
	} else if temp == 2 {
		fmt.Println("●")
	}

	temp = b.returnPieceTypeByPosition(0, 1)
	if temp == 0 {
		fmt.Println(".")
	} else if temp == 1 {
		fmt.Println("○")
	} else if temp == 2 {
		fmt.Println("●")
	}

	b.boardprint()
}
