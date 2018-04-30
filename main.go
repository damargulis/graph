package graph

import (
	"github.com/fogleman/gg"
	"io/ioutil"
	"strconv"
	"strings"
)

type Point struct {
	p1depth, p2depth, p1wins, p2wins, ties int
}

func getPoints(fileIn string) ([]Point, int) {
	max := 0
	dat, _ := ioutil.ReadFile(fileIn)
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	x := make([]Point, len(lines))
	for i, line := range lines {
		var lineInt [5]int
		l := strings.Split(strings.TrimSpace(line), ",")
		for j := range lineInt {
			num, _ := strconv.Atoi(l[j])
			lineInt[j] = num
		}
		x[i] = Point{
			p1depth: lineInt[0],
			p2depth: lineInt[1],
			p1wins:  lineInt[2],
			p2wins:  lineInt[3],
			ties:    lineInt[4],
		}
		if x[i].p1depth > max {
			max = x[i].p1depth
		}
		if x[i].p2depth > max {
			max = x[i].p2depth
		}
	}
	return x, max + 1
}

func makeGraph(title, fileIn, fileOut string, ch chan int) {
	const S = 1024
	const P = 64
	points, maxWidth := getPoints(fileIn)
	dc := gg.NewContext(S+100, S)
	dc.InvertY()
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Translate(P, P)
	dc.Scale(S-P*2, S-P*2)

	for i := 1; i <= maxWidth; i++ {
		x := float64(i) / float64(maxWidth)
		dc.MoveTo(x, 0)
		dc.LineTo(x, 1)
		dc.MoveTo(0, x)
		dc.LineTo(1, x)
	}
	dc.SetRGBA(0, 0, 0, 0.25)
	dc.SetLineWidth(1)
	dc.Stroke()

	dc.MoveTo(0, 0)
	dc.LineTo(1, 0)
	dc.MoveTo(0, 0)
	dc.LineTo(0, 1)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(4)
	dc.Stroke()

	for _, point := range points {
		totalGames := float64(point.p1wins + point.p2wins + point.ties)
		dc.SetRGBA(float64(point.p1wins)/totalGames, float64(point.p2wins)/totalGames, float64(point.ties)/totalGames, 1)
		dc.DrawRectangle(float64(point.p1depth)/float64(maxWidth), float64(point.p2depth)/float64(maxWidth), 1.0/float64(maxWidth), 1.0/float64(maxWidth))
		dc.Fill()
	}

	dc.Identity()
	dc.SetRGB(0, 0, 0)

	dc.DrawStringAnchored(title+" - Alphabeta vs. Alphabeta", S/2, P/2, 0.5, 0.5)


	dc.DrawStringAnchored("Player 1 Lookahead", S/2, S-P/2, 0.5, 0.5)
	dc.RotateAbout(-1.5708, S/2, S/2)
	dc.DrawStringAnchored("Player 2 Lookahead", S/2, P/2, 0.5, 0.5)

	dc.Identity()
	dc.DrawStringAnchored("0", P, S-(P-10), 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(maxWidth/4), (P+S/2)/2, S-(P-10), 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(maxWidth/2), S/2, S-(P-10), 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(3*maxWidth/4), (S/2+S-P)/2, S-(P-10), 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(maxWidth), S-P, S-(P-10), 0.5, 0.5)
	dc.RotateAbout(-1.5708, S/2, S/2)
	dc.DrawStringAnchored("0", P, P-10, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(maxWidth/4), (P+S/2)/2, P-10, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(maxWidth/2), S/2, P-10, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(3*maxWidth/4), (S-P+S/2)/2, P-10, 0.5, 0.5)
	dc.DrawStringAnchored(strconv.Itoa(maxWidth), S-P, P-10, 0.5, 0.5)

	dc.Identity()

	dc.Translate(S, 0)
	dc.DrawStringAnchored("Legend", 25, 100, 0.5, 0.5)
	dc.SetRGB(1, 0, 0)
	dc.DrawRectangle(10, 150, 30, 30)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Player 1 Win", 25, 200, 0.5, 0.5)

	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(10, 250, 30, 30)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Player 2 Win", 25, 300, 0.5, 0.5)

	dc.SetRGB(0, 0, 1)
	dc.DrawRectangle(10, 350, 30, 30)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Tie", 25, 400, 0.5, 0.5)

	dc.SavePNG(fileOut)
	ch <- 1
}

type Graph struct {
	title, fileIn, fileOut string
}

func CreateGraph(title string, fileIn string, fileOut string) {
	ch := make(chan int)
	go makeGraph(title, fileIn, fileOut, ch)
	<-ch
}

func main() {
	graphs := [9]Graph{
		Graph{"Boxes", "../game/boxes.csv", "boxes.png"},
		Graph{"Checkers", "../game/checkers.csv", "checkers.png"},
		Graph{"Connect 4", "../game/connect4.csv", "connect4.png"},
		Graph{"Mancala", "../game/mancala.csv", "mancala.png"},
		Graph{"Reversi", "../game/reversi.csv", "reversi.png"},
		Graph{"TicTacToe", "../game/tictactoe.csv", "tictactoe.png"},
		Graph{"Nine Man's Morris", "../game/ninemansmorris.csv", "ninemansmorris.png"},
		Graph{"Abalone", "../game/abalone.csv", "abalone.png"},
		Graph{"New Tic Tac Toe", "../game/newtictactoe.csv", "newtictactoe.png"},
	}
	ch := make(chan int)
	for _, graph := range graphs {
		go makeGraph(graph.title, graph.fileIn, graph.fileOut, ch)
	}
	for range graphs {
		<-ch
	}
}
