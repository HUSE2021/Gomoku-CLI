package main

import (
	"bufio"
	"fmt"
	term "github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// run "go get github.com/nsf/termbox-go"
// null 0,  user1:○ = 1　user2: ● = 2  select time =-1
var haveWinner bool = false
var regretStack []Piece
var boardSize int = 0
var errorMessage string = ""

var winmessage1 = [...]string{
	"██████░░░███░██░░░░░░░░██░██░███░░░██",
	"██░░░██░████░██░░░██░░░██░██░████░░██",
	"██████░██░██░░██░████░██░░██░██░██░██",
	"██░░░░░░░░██░░████░░████░░██░██░░████",
	"██░░░░░███████░██░░░░██░░░██░██░░░███",
	"»»—————————————-　★　—————————————-««",
	"░░░░░░░░░░Congratulations!!░░░░░░░░░░",
}
var winmessage2 = [...]string{
	"██████░█████░██░░░░░░░░██░██░███░░░██",
	"██░░░██░░░██░██░░░██░░░██░██░████░░██",
	"██████░░███░░░██░████░██░░██░██░██░██",
	"██░░░░░██░░░░░████░░████░░██░██░░████",
	"██░░░░███████░░██░░░░██░░░██░██░░░███",
	"»»—————————————-　★　—————————————-««",
	"░░░░░░░░░░Congratulations!!░░░░░░░░░░",
}

type Piece struct {
	x    int
	y    int
	user int
}

type Board struct {
	tokens    []int
	userName  [2]string
	boardSize int
}

func (b *Board) InitialBoard(boardSize int) {
	b.boardSize = boardSize
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
	if len(regretStack) == 0 {
		return false
	} else {
		fmt.Println("stack top is: ", regretStack[len(regretStack)-1].x, "and ", regretStack[len(regretStack)-1].y, "by: ", regretStack[len(regretStack)-1].user, "\nnow size: ", len(regretStack))
		if b.putPiece(regretStack[len(regretStack)-1].x, regretStack[len(regretStack)-1].y, 0) {
			regretStack = regretStack[:len(regretStack)-1]
		}
		return true
	}
}

func (b *Board) putPiece(x, y, userType int) bool {
	var boardSize = b.boardSize
	fmt.Println(x, "+", y, "+", userType)
	if userType == 0 {
		b.tokens[x*boardSize+y] = 0
		return true
	}
	if b.checkNotOverFlow(x, y) == true {
		b.tokens[x*boardSize+y] = userType
		if b.check5Piece(x, y, userType) {
			haveWinner = true
		}
		return true //200 is ok, 500 is not ok
	}
	return false
}

func (b *Board) returnPieceTypeByPosition(x, y int) int {
	var boardSize = b.boardSize
	if b.checkNotOverFlow(x, y) == true {
		if b.tokens[x*boardSize+y] != 0 {
			if b.tokens[x*boardSize+y] == 1 {
				return 1
			} else if b.tokens[x*boardSize+y] == 2 {
				return 2
			} else {
				return -1
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
				} else if j == boardSize-1 {
					fmt.Printf("─● ")
				} else {
					fmt.Printf("─●─")
				}
			case -1:
				if j == 0 {
					fmt.Printf(" X─")
				} else if j == boardSize-1 {
					fmt.Printf("─X ")
				} else {
					fmt.Printf("─X─")
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
	var boardSize = b.boardSize
	magicNumber := 0 //if you say what's is, IDK,but work
	if boardSize == 15 {
		magicNumber = 4
	} else {
		magicNumber = 6
	}
	fmt.Printf("   ")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%2d", i)
		fmt.Printf(" ")
	}
	fmt.Println("")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%3d", i)
		if i == boardSize/2-4 {
			fmt.Printf(" ╔══")
			for k := 0; k < boardSize-3; k++ {
				fmt.Printf("═══")
			}
			fmt.Printf("═══╗")
		} else if i == boardSize/2+4 {
			fmt.Printf(" ╚══")
			for k := 0; k < boardSize-3; k++ {
				fmt.Printf("═══")
			}
			fmt.Printf("═══╝")
		} else if i >= boardSize/2-3 && i <= (boardSize/2)+3 {
			fmt.Printf(" ║")
			space := ((boardSize-2)*3 - 36 + 1)
			left := int(space / 2)
			right := space - left
			for k := 0; k < left; k++ {
				fmt.Printf(" ")
			}
			if nowUser == 1 {
				fmt.Printf(winmessage1[i-magicNumber])
			} else {
				fmt.Printf(winmessage2[i-magicNumber])
			}
			for k := 0; k < right; k++ {
				fmt.Printf(" ")
			}
			fmt.Printf("║ ")
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
	var boardSize = b.boardSize
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

func reset() {
	term.Sync() // cosmestic purpose
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

func (b *Board) checkNotOverFlow(x, y int) bool {
	var boardSize = b.boardSize
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

func (b *Board) getBoardSize() {
	getBoardSize := "A"
	fmt.Println("=====  Game Start  =====")
	for b.boardSize == 0 {
		fmt.Println("=====  Select Board Size  =====\n   =====  (A)15  (B)19   =====")
		_, err := fmt.Scan(&getBoardSize)
		if err != nil {
			return
		}
		if getBoardSize == "A" || getBoardSize == "a" {
			b.InitialBoard(15)
		} else if getBoardSize == "B" || getBoardSize == "b" {
			b.InitialBoard(19)
		} else {
			fmt.Println("======  Bad Input, Again  ======")
		}
	}
}

func (b *Board) getUserName() {
	for len(b.userName[0]) == 0 {
		fmt.Println("=====  Plz input user1 name:  =====")
		_, errG := fmt.Scan(&b.userName[0])
		if errG != nil {
			return
		}
	}
	for len(b.userName[1]) == 0 {
		fmt.Println("=====  Plz input user2 name:  =====")
		_, err := fmt.Scan(&b.userName[1])
		if err != nil {
			return
		}
		if b.userName[0] == b.userName[1] {
			fmt.Println("=====  Cannot same name with user1 !!!  =====")
			b.userName[1] = ""
		}
	}
	for {
		fmt.Printf("=====  who is first ?  =====\n=====  (A)%s  (B)%s (C)random  =====\n", b.userName[0], b.userName[1])
		choice := "C"
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("===== bad choice =====")
			return
		}
		if choice == "A" || choice == "a" {
			break
		} else if choice == "B" || choice == "b" {
			var temp = b.userName[0]
			b.userName[0] = b.userName[1]
			b.userName[1] = temp
			break
		} else if choice == "C" || choice == "c" {
			r := rand.New(rand.NewSource(time.Now().Unix()))
			var ra = r.Intn(2)
			if ra == 1 {
				var temp = b.userName[0]
				b.userName[0] = b.userName[1]
				b.userName[1] = temp
			}
			fmt.Printf("=====  Fist is %s  =====\n", b.userName[0])
			break
		} else {
			fmt.Println("=====  Bad Input, Again  =====")
		}
	}
}

func (b *Board) keyGet() string {
	CallClear()
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	b.boardPrint()
	fmt.Println(regretStack)
	fmt.Println(errorMessage)
	fmt.Println("Enter ↑ ↓ ← → to select, SPACE to regret")
keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyArrowUp:
				reset()
				fmt.Println("Arrow Up pressed")
				return "up"
			case term.KeyArrowDown:
				reset()
				fmt.Println("Arrow Down pressed")
				return "down"
			case term.KeyArrowLeft:
				reset()
				fmt.Println("Arrow Left pressed")
				return "left"
			case term.KeyArrowRight:
				reset()
				fmt.Println("Arrow Right pressed")
				return "right"
			case term.KeyEnter:
				reset()
				fmt.Println("Enter pressed")
				return "enter"
			case term.KeySpace:
				reset()
				fmt.Println("Backspace pressed")
				return "backspace"
			default:
				// we only want to read a single character or one key pressed event
				reset()
				return "other"
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
	return "error"
}

func main() {
	var b Board
	var nowUser int
	var keyWait string
	var nowPositionUser int
	var nowSelect Piece
	//var nowPosition Piece

	fmt.Println("Gomoku  rulu\n Players alternate turns placing a stone of their color on an empty " +
		"intersection. Black plays first. The winner is the first player to form an unbroken chain of " +
		"five stones horizontally, vertically, or diagonally. Placing so that a line of more than five " +
		"stones of the same color is created does not result in a win. These are called overlines.")
	fmt.Println("　　　　　　ﾊ,,ﾊ\n　　　　　( ﾟωﾟ ) ")
	fmt.Print("  Press 'Enter' to continue... ")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	rand.Seed(time.Now().Unix())
	b.getUserName()
	b.getBoardSize()

	nowUser = 1
	fmt.Println("	    ======  Black First ======")
	b.boardPrint()

	nowSelect.x = 7
	nowSelect.y = 7
	nowSelect.user = -1

	nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
	b.putPiece(nowSelect.x, nowSelect.y, -1)

	for {
		keyWait = ""
		keyWait = b.keyGet()
		errorMessage = ""
		b.putPiece(nowSelect.x, nowSelect.y, nowPositionUser)

		fmt.Println(regretStack)

		if keyWait == "left" {
			if b.checkNotOverFlow(nowSelect.x, nowSelect.y-1) {
				nowSelect.y = nowSelect.y - 1
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			} else {
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			}

		}
		if keyWait == "right" {
			if b.checkNotOverFlow(nowSelect.x, nowSelect.y+1) {
				nowSelect.y = nowSelect.y + 1
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			} else {
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			}
		}
		if keyWait == "down" {
			if b.checkNotOverFlow(nowSelect.x+1, nowSelect.y) {
				nowSelect.x = nowSelect.x + 1
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			} else {
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			}
		}
		if keyWait == "up" {
			if b.checkNotOverFlow(nowSelect.x-1, nowSelect.y) {
				nowSelect.x = nowSelect.x - 1
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			} else {
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			}
		}
		if keyWait == "enter" {
			if b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y) == 0 {
				b.putPiece(nowSelect.x, nowSelect.y, nowUser)
				regretStack = append(regretStack, Piece{nowSelect.x, nowSelect.y, nowUser})
				nowPositionUser = nowUser
				CallClear()
				errorMessage = "user " + b.userName[nowUser-1] + " put in " + strconv.Itoa(nowSelect.x) + " " + strconv.Itoa(nowSelect.y)
				b.boardPrint()
				if haveWinner == true {
					CallClear()
					b.winPrint(nowUser)
					println("Congratulations!! " + b.userName[nowUser-1])
					return
				}
				changeUser(&nowUser)
			} else {
				errorMessage = "pity! There are already chess pieces on it"
				nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
				b.putPiece(nowSelect.x, nowSelect.y, -1)
			}
		}
		if keyWait == "backspace" {
			if !b.regret() {
				errorMessage = ("	YOU CAN NOT REGRET ")
			} else {
				changeUser(&nowUser)
				errorMessage = b.userName[nowUser-1] + "regret!"
			}
			nowPositionUser = b.returnPieceTypeByPosition(nowSelect.x, nowSelect.y)
			b.putPiece(nowSelect.x, nowSelect.y, -1)
		}
		if keyWait == "other" {
			errorMessage = "You have entered the wrong command"
		}
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
