package goblocks

import (
	"testing"
	"fmt"


)

func CheckFillmap(t*testing.T, actual FillMap, expected string) {
	fm := actual.String();
	if (fm != expected) {
		t.Errorf("incorrect fillmap : \n%s expected: \n%s", fm, expected)
	}
}

func TestPieceFillmap(t*testing.T) {
	CheckFillmap(t, tSmall.FillMap(0, 0),
		"|x   |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, tVerti.FillMap(0, 0),
		"|x   |\n" +
		  "|x   |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, tVerti.FillMap(1, 0),
		"| x  |\n" +
		  "| x  |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, tVerti.FillMap(1, 1),
		"|    |\n" +
		  "| x  |\n" +
		  "| x  |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, tHoriz.FillMap(0, 0),
		"|xx  |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, tHoriz.FillMap(2, 4),
		"|    |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|  xx|\n")

	CheckFillmap(t, tBigSq.FillMap(0, 0),
		"|xx  |\n" +
		  "|xx  |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, tBigSq.FillMap(1, 3),
		"|    |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "| xx |\n" +
		  "| xx |\n")
}

func TestPointLocs(t *testing.T) {
	CheckFillmap(t, locationFillMap(1, 2),
		"|    |\n" +
		  "|    |\n" +
		  "| x  |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, locationFillMap(1, 2) | locationFillMap(1, 1),
		"|    |\n" +
		  "| x  |\n" +
		  "| x  |\n" +
		  "|    |\n" +
		  "|    |\n")

	CheckFillmap(t, locationFillMap(2, 3) | locationFillMap(2, 4) | locationFillMap(3, 3) | locationFillMap(3, 4),
		"|    |\n" +
		  "|    |\n" +
		  "|    |\n" +
		  "|  xx|\n" +
		  "|  xx|\n")

}

var possibleMoveTests = []struct {
	start Board
	moves PossMoves
}{
	{MakeBoard(default_start),
		[10]PossMove{noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, PossMove(MoveUp)}},
	{MakeBoard(default_start).MovePiece(9, MoveUp).MovePiece(4, MoveUp),
		[10]PossMove{noMovePoss, PossMove(MoveUp), PossMove(MoveUp), noMovePoss, PossMove(MoveDown), noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss}},
	{MakeBoard(default_start).MovePiece(9, MoveUp).MovePiece(4, MoveUp).MovePiece(1, MoveUp),
		[10]PossMove{PossMove(MoveRight), PossMove(MoveDown | MoveRight),PossMove(MoveUp | MoveLeft), noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss, noMovePoss}},
}


func TestPossibleMoves(t *testing.T) {
	for _, tt := range (possibleMoveTests) {
		possMoves := tt.start.PossibleMoves()
		if (possMoves != tt.moves) {
			t.Errorf("board \n%s has possible moves %s, expecting %s ", tt.start.FillMap(), possMoves, tt.moves)
		}
	}
}

func TestMakeFillmap(t *testing.T) {

	CheckFillmap(t, makeFillmap(default_start),
		"|x  x|\n" +
		  "|xxxx|\n" +
		  "|xxxx|\n" +
		  "|xxxx|\n" +
		  "|xxxx|\n")

}

var boardIdTests = []struct {
	b  Board
	id Board
}{
	{MakeBoard([]loc{}), 0},
	{MakeBoard([]loc{}).MovePiece(0, MoveDown), 1},
	{MakeBoard([]loc{}).MovePiece(0, MoveRight), 8},
	{MakeBoard(default_start), 317487014097564},

}

func TestMakeId(t *testing.T) {
	for _, tt := range (boardIdTests) {
		if tt.b.Identity() != tt.id {
			t.Errorf("board %s has incorrect idenity, got %d, expecting %d", tt.b, tt.b.Identity(), tt.id)
		}
	}
}



func TestSorting(t *testing.T) {
	b := MakeBoard(default_start)
	id := b.Identity()
	id3 := id.Identity()
	if(id.FillMap() != b.FillMap()){
		t.Errorf("Sorting produced different filmap %s, %s\n",id.FillMap(), b.FillMap())
	}

	if(id != id3){
		t.Errorf("Sorting is not idempotent %s != %s\n",id, id3)
	}


}

var boardLocationTests = []struct {
	b Board
	p int
	x int
	y int
}{
	{MakeBoard([]loc{{0, 0}}), 0, 0, 0},
	{MakeBoard([]loc{{0, 1}}), 0, 0, 1},
	{MakeBoard([]loc{{3, 4}}), 0, 3, 4},
	{MakeBoard([]loc{{3, 4}, {0, 0}}), 1, 0, 0},
	{MakeBoard([]loc{{3, 4}, {2, 3}}), 1, 2, 3},

}

func TestGetBoardLocations(t *testing.T) {
	for _, tt := range (boardLocationTests) {
		x, y := tt.b.PiecePos(tt.p);
		if (x != tt.x || y != tt.y) {
			t.Errorf("Invalid location for piece %d (%d, %d) on board %s, expecting (%d,%d)", tt.p, x, y, tt.b, tt.x, tt.y)
		}
	}
}



func TestSearch(T *testing.T) {

	shortestPaths := make(map[Board]int)
	boards := PopulateGraph(MakeBoard(default_start).Identity())
	shortestPaths[boards.start] =1
	fmt.Printf("Seen  %d boards \n", len (boards.transitions))

	count :=0
	boards.BoardByArcsBreadthFirst(func (b Board, via PieceMove , to Board){
		curDepth := shortestPaths[b]
		if(curDepth  < 1 ){
			panic("Invalid depth on " + b.String() + " -> " + via.String() + " -> " + to.String())
		}
		curDist := shortestPaths[to]
		if curDist ==0 || ( curDist > (curDepth +1)) {
			shortestPaths[to] = curDepth +1
		}
		count++
		if count%10000 ==0{
		}
	});

	fmt.Printf("Labelled %d nodes\n", len(shortestPaths))

	shortest := 0
	var shortestBoard Board

	for b,depth := range(shortestPaths){
		px,py :=  b.PiecePos(9)
		if px == 1 && py == 3 {
			if shortest == 0 || depth < shortest {
				shortest = depth
				shortestBoard = b
			}

		}
	}

	fmt.Printf("Shortest path is  %d  and solution is %s\n", shortest,shortestBoard)



	cur := shortestBoard
	path := []PieceMove{}
	for {
		var closest Board
		var shortestNeighbordDist = 0
		var lastMove PieceMove
		for bm,b := range(boards.transitions[cur]){
			neighborDist:= shortestPaths[b]
			if shortestNeighbordDist == 0 ||(neighborDist < shortestNeighbordDist) {
				closest = b
				shortestNeighbordDist  = neighborDist
				lastMove = bm
			}
		}

		path = append(path,lastMove)
		cur = closest
		if cur == boards.start {
			break
		}
	}

	bd := boards.start
	fmt.Printf("initial position %s\n", bd)
	for i := len(path)-1 ; i>=0; i-- {
		mv := path[i].mv.Reverse()
		bd = bd.MovePiece(path[i].piece,mv)
		fmt.Printf("%d : %s %s\n", path[i].piece, path[i].mv.Reverse(), bd)
	}
}