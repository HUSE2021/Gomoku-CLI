package main

import (
	"fmt"
	"io"
)

// null 0,  user1:○ = 1　user2: ● = 2
var haveWinner bool = false

type Board struct {
	tokens [15 * 15]int
}

func (b *Board) InitialBoard() {
	var i int

	for i = 0; i < 15*15; i++ {
		b.tokens[i] = 0
	}
}

func (b *Board) putPiece(x, y, userType int) bool {
	if checkNotOverFlow(x, y) == true {
		if b.tokens[x*15+y] == 0 {
			b.tokens[x*15+y] = userType
			if b.check5Piece(x, y, userType) {
				haveWinner = true
			}
			return true //200 is ok, 500 is not ok
		} else {
			return false
		}
	} else {
		return false
	}

}

func (b *Board) returnPieceTypeByPosition(x, y int) int {
	if checkNotOverFlow(x, y) == true {
		if b.tokens[x*15+y] != 0 {
			if b.tokens[x*15+y] == 1 {
				return 1
			} else {
				return 2
			}
		} else {
			return 0
		}
	} else {
		return 0
	}

}

func (b *Board) boardprint() int {
	fmt.Println("    0  1  2  3  4  5  6  7  8  9 10 11 12 13 14")
	for i := 0; i < 15; i++ {
		fmt.Printf("%3d", i)
		for j := 0; j < 15; j++ {
			switch b.tokens[i*15+j] {
			case 0:
				if i%4 == 3 && j%4 == 3 {
					fmt.Printf(" _ ")
				} else {
					fmt.Printf(" . ")
				}
			case 1:
				fmt.Printf(" ○ ")
			case 2:
				fmt.Printf(" ● ")
			default:
				fmt.Println("Error:Unexpected Token")
				return 1
			}
		}
		fmt.Println("")
	}
	return 0
}

func (b *Board) check5Piece(x, y, userType int) bool {
	xcount, ycount, zcount := 0, 0, 0
	x2, y2 := x, y

	for i := 0; i < 15; i++ {
		//"-"
		if xcount == 5 {
			return true
		}
		if b.tokens[x*15+i] == userType {

			xcount++
		} else {
			xcount = 0
		}
		//"|"
		if ycount == 5 {
			return true
		}
		if b.tokens[i*15+y] == userType {
			ycount++
		} else {
			ycount = 0
		}
	}
	// "/"
	for x2 > 0 && y2 < 15 {
		x2--
		y2++
	}
	for x2 < 15 && y2 > 0 {
		if zcount == 5 {
			return true
		}
		if b.tokens[x2*15+y2] == userType {
			zcount++
		} else {
			zcount = 0
		}
		x2++
		y2--
	}
	zcount = 0
	// "\"
	for x > 0 && y > 0 {
		x--
		y--
	}
	for x < 15 && y < 15 {
		if zcount == 5 {
			return true
		}
		if b.tokens[x*15+y] == userType {
			zcount++
		} else {
			zcount = 0
		}
		x++
		y++
	}
	return false
}

//func startXY1(x,y int)(i,j int) {
//	if x>y{
//		i=0
//		j=x-y
//	}else {
//		i=j-i
//		j=0
//	}
//	return i,j
//}
//func startXY2(x,y int)(i,j int) {
//	if x>y{
//		i=x+y
//		j=15
//	}else {
//		i=j-i
//		j=0
//	}
//	return i,j
//}

func checkNotOverFlow(x, y int) bool {
	if x >= 0 && x < 14 && y >= 0 && y < 15 {
		return true
	} else {
		return false
	}
}

func changeUser(nowUser *int) {
	if *nowUser == 1 {
		*nowUser++
	} else {
		*nowUser--
	}
}

func main() {
	var b Board
	var x, y int
	var nowUser int

	//Initial Game
	b.InitialBoard()
	nowUser = 1
	fmt.Println("	    ======  Game Start  ======")
	fmt.Println("	    ======  Black First ======")
	b.boardprint()
	fmt.Printf("user: %d  plz input :", nowUser)
	for {
		_, err := fmt.Scan(&x, &y)
		if err == io.EOF {
			break
		}
		if b.putPiece(x, y, nowUser) {
			fmt.Printf("user: %d  put in: %d,%d\n", nowUser, x, y)
			b.boardprint()
			if haveWinner == true {
				fmt.Printf("winner is : %d\n", nowUser)
				return
			}
			changeUser(&nowUser)
		} else {
			fmt.Printf("bad input ,again\n")
		}
		fmt.Printf("user: %d  plz input :", nowUser)

	}

	//var temp int
	//temp = b.returnPieceTypeByPosition(0, 0)
	//if temp == 0 {
	//	fmt.Println(".")
	//} else if temp == 1 {
	//	fmt.Println("○")
	//} else if temp == 2 {
	//	fmt.Println("●")
	//}
	//
	//temp = b.returnPieceTypeByPosition(0, 1)
	//if temp == 0 {
	//	fmt.Println(".")
	//} else if temp == 1 {
	//	fmt.Println("○")
	//} else if temp == 2 {
	//	fmt.Println("●")
	//}

}
