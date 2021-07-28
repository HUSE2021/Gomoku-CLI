package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// null 0,  user1:○ = 1　user2: ● = 2
var haveWinner bool = false
var boardSize int = 0
var regretStact []Piece

type Board struct {
	tokens []int
}

func (b *Board) InitialBoard() {
	b.tokens = make([]int, boardSize*boardSize)
}

//flash terminal
var clear map[string]func() //create a map for storing clear funcs
func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func (b *Board) regret() bool {
	if len(regretStact) == 0 {
		return false
	} else {
		fmt.Println("stack top is: ", regretStact[len(regretStact)-1].x, "and ", regretStact[len(regretStact)-1].y, "now size: ", len(regretStact))
		if b.putPiece(regretStact[len(regretStact)-1].x, regretStact[len(regretStact)-1].y, 0) {
			regretStact = regretStact[:len(regretStact)-1]
		}
		return true
	}
}

func (b *Board) putPiece(x, y, userType int) bool {
	fmt.Println(x, "+", y, "+", userType)
	if userType == 0 {
		b.tokens[x*boardSize+y] = 0
		return true
	}
	if checkNotOverFlow(x, y) == true {
		if b.tokens[x*boardSize+y] == 0 {
			b.tokens[x*boardSize+y] = userType
			if b.check5Piece(x, y, userType) {
				haveWinner = true
			}
			return true //200 is ok, 500 is not ok
		}
	}
	return false
}

func (b *Board) returnPieceTypeByPosition(x, y int) int {
	if checkNotOverFlow(x, y) == true {
		if b.tokens[x*boardSize+y] != 0 {
			if b.tokens[x*boardSize+y] == 1 {
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


func (b *Board) boardPrint() int {
	var boardSize = b.boardSize
	fmt.Printf("   ")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%2d", i)
		fmt.Printf(" ")
	}
	fmt.Println("")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%3d", i)
		for j := 0; j < boardSize; j++ {
			switch b.tokens[i*boardSize+j] {
			case 0:
				if i == 0 && j == 0 {
					fmt.Printf(" ┌─")
				} else if i == 0 && j == boardSize-1 {
					fmt.Printf("─┐ ")
				} else if i == boardSize-1 && j == 0 {
					fmt.Printf(" └─")
				} else if i == boardSize-1 && j == boardSize-1 {
					fmt.Printf("─┘ ")
				} else if j == 0 {
					fmt.Printf(" ├─")
				} else if j == boardSize-1 {
					fmt.Printf("─┤ ")
				} else if i == 0 {
					fmt.Printf("─┬─")
				} else if i == boardSize-1 {
					fmt.Printf("─┴─")
				} else {
					fmt.Printf("─┼─")
				}
			case 1:
				if j == 0 {
					fmt.Printf(" ○─")
				} else if j == boardSize-1 {
					fmt.Printf("─○ ")
				} else {
					fmt.Printf("─○─")
				}
			case 2:
				if j == 0 {
					fmt.Printf(" ●─")
				} else if j == boardSize {
					fmt.Printf("─● ")
				} else {
					fmt.Printf("─●─")
				}
			default:
				fmt.Println("Error:Unexpected Token")
				return 1
			}
		}
		fmt.Println("")
	}
	return 0
}

func (b *Board) winPrint(nowUser int) int {
	fmt.Printf("   ")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%2d", i)
		fmt.Printf(" ")
	}
	fmt.Println("")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%3d", i)
		if i == boardSize/2-3 {
			fmt.Printf("┌──")
			for k := 0; k < boardSize-2; k++ {
				fmt.Printf("───")
			}
			fmt.Printf("──┐")
		} else if i == boardSize/2+2 {
			fmt.Printf("└──")
			for k := 0; k < boardSize-2; k++ {
				fmt.Printf("───")
			}
			fmt.Printf("──┘")
		} else if i >= int(boardSize/2)-2 && i <= boardSize/2+1 {
			fmt.Printf("│  ")
			space := ((boardSize-2)*3 - 36 + 1)
			left := int(space / 2)
			right := space - left
			for k := 0; k < left; k++ {
				fmt.Printf(" ")
			}
			if nowUser == 1 {
				fmt.Printf(winmessage1[i-int(boardSize/2)+2])
			} else {
				fmt.Printf(winmessage2[i-int(boardSize/2)+2])
			}
			for k := 0; k < right; k++ {
				fmt.Printf(" ")
			}
			fmt.Printf("  │")
		} else {
			for j := 0; j < boardSize; j++ {
				switch b.tokens[i*boardSize+j] {
				case 0:
					if i == 0 && j == 0 {
						fmt.Printf(" ┌─")
					} else if i == 0 && j == boardSize-1 {
						fmt.Printf("─┐ ")
					} else if i == boardSize-1 && j == 0 {
						fmt.Printf(" └─")
					} else if i == boardSize-1 && j == boardSize-1 {
						fmt.Printf("─┘ ")
					} else if j == 0 {
						fmt.Printf(" ├─")
					} else if j == boardSize-1 {
						fmt.Printf("─┤ ")
					} else if i == 0 {
						fmt.Printf("─┬─")
					} else if i == boardSize-1 {
						fmt.Printf("─┴─")
					} else {
						fmt.Printf("─┼─")
					}
				case 1:
					if j == 0 {
						fmt.Printf(" ○─")
					} else if j == boardSize-1 {
						fmt.Printf("─○ ")
					} else {
						fmt.Printf("─○─")
					}
				case 2:
					if j == 0 {
						fmt.Printf(" ●─")
					} else if j == boardSize {
						fmt.Printf("─● ")
					} else {
						fmt.Printf("─●─")
					}
				default:
					fmt.Println("Error:Unexpected Token")
					return 1
				}
			}
		}
		fmt.Println("")
	}
	return 0
}


func (b *Board) check5Piece(x, y, userType int) bool {
	xcount, ycount, zcount := 0, 0, 0
	x2, y2 := x, y
	if userType == 0 {
		return false
	}
	for i := 0; i < boardSize; i++ {
		//"-"
		if xcount == 5 {
			return true
		}
		if b.tokens[x*boardSize+i] == userType {

			xcount++
		} else {
			xcount = 0
		}
		//"|"
		if ycount == 5 {
			return true
		}
		if b.tokens[i*boardSize+y] == userType {
			ycount++
		} else {
			ycount = 0
		}
	}
	// "/"
	for x2 > 0 && y2 < boardSize {
		x2--
		y2++
	}
	for x2 < boardSize && y2 > 0 {
		if zcount == 5 {
			return true
		}
		if b.tokens[x2*boardSize+y2] == userType {
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
	for x < boardSize && y < boardSize {
		if zcount == 5 {
			return true
		}
		if b.tokens[x*boardSize+y] == userType {
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
	if x >= 0 && x < boardSize && y >= 0 && y < boardSize {
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

type Piece struct {
	x int
	y int
}



func (b *Board) getUserName() {
	for len(b.userName[0]) == 0 {
		fmt.Println("	    ======  Plz input user1 name:  ======")
		_, errG := fmt.Scan(&b.userName[0])
		if errG != nil {
			return
		}
	}
	for len(b.userName[1]) == 0 {
		fmt.Println("	    ======  Plz input user2 name:  ======")
		_, err := fmt.Scan(&b.userName[1])
		if err != nil {
			return
		}
		if b.userName[0] == b.userName[1] {
			fmt.Println("	    ======  Cannot same name with user1 !!!  ======")
			b.userName[1] = ""
		}
	}
	for {
		fmt.Printf("	    ======  who is first ?  \n   ======  (A)%s  (B)%s (C)random  ============\n", b.userName[0], b.userName[1])
		choice := "C"
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("bad choice")
			return
		}
		if choice == "A" {
			break
		} else if choice == "B" {
			var temp = b.userName[0]
			b.userName[0] = b.userName[1]
			b.userName[1] = temp
			break
		} else if choice == "C" {
			r := rand.New(rand.NewSource(time.Now().Unix()))
			var ra = r.Intn(2)
			if ra == 1 {
				var temp = b.userName[0]
				b.userName[0] = b.userName[1]
				b.userName[1] = temp
			}
			fmt.Printf("	    ======  Fist is %s  ======\n", b.userName[0])
			break
		} else {
			fmt.Println("	    ======  Bad Input, Again  ======")
		}
	}
}


func main() {
	var b Board
	var xInput, yInput string
	var x, y int
	var nowUser int

	rand.Seed(time.Now().Unix())
	b.getUserName()
	b.getBoardSize()

	nowUser = 1
	fmt.Println("	    ======  Black First ======")
	b.boardPrint()
	for {
		fmt.Println(regretStack)
		fmt.Printf("user:%s  plz input （input 'R' to regret）:", b.userName[nowUser-1])
		_, err := fmt.Scanln(&xInput, &yInput)
		if err == io.EOF {
			break
		}
		//regret
		if xInput == "R" {
			CallClear()
			if !b.regret() {
				fmt.Println("	    ======  YOU CAN NOT REGRET ======")
			} else {
				changeUser(&nowUser)
			}
			b.boardPrint()
		} else {
			if xInput == "" || yInput == "" {
				fmt.Println("	    ======  Bad X and Y ======")
				continue
			}
			//is not regret So let input text as X and Y
			x, err = strconv.Atoi(xInput)
			y, err = strconv.Atoi(yInput)
			xInput = ""
			yInput = ""
			if err == io.EOF {
				break
			}

			if b.putPiece(x, y, nowUser) {
				regretStack = append(regretStack, Piece{x, y})
				CallClear()
				fmt.Printf("user: %d  put in: %d,%d\n", nowUser, x, y)
				b.boardPrint()
				if haveWinner == true {
					b.winPrint(nowUser)
					return
				}
				changeUser(&nowUser)
				b.winPrint(1)
			} else {
				CallClear()
				b.boardPrint()
				fmt.Printf("bad input ,again\n")
			}
		}
}
