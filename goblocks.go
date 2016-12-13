package goblocks

import (
	"fmt"
	"golang.org/x/net/html/atom"
)

const (
	nPieces = 10
	xBits = 2
	yBits = 3
	bitsPerPiece = xBits + yBits
	xMax = 4
	yMax = 5
)

type PieceType uint8

type PieceMove struct {
	piece int
	mv    Move
}

const (
	tSmall = PieceType(0)
	tHoriz = PieceType(1)
	tVerti = PieceType(2)
	tBigSq = PieceType(3)
)

var pTypes = [nPieces]PieceType{tSmall, tSmall, tSmall, tSmall, tHoriz, tVerti, tVerti, tVerti, tVerti, tBigSq}

const (
	noMovePoss = PossMove(0)
	NoMove = Move(0)
	MoveUp = Move(iota)
	MoveDown = Move(1 << iota)
	MoveLeft = Move(1 << iota)
	MoveRight = Move(1 << iota)
)

type Move uint8

type PossMove Move

type PossMoves [nPieces]PossMove

type loc struct {
	x int;
	y int;
}

type FillMap uint32;

const (
	fmEmpty = FillMap(0)
	fmInvalid = FillMap(0xffffffff)
)

type Board uint64

var AvailableMoves = []Move{MoveUp, MoveDown, MoveLeft, MoveRight}

// Transform generates the coordinates after a move
func (m Move) Transform(x int, y int) (int, int) {
	switch(m){
	case MoveLeft:
		return x - 1, y
	case MoveRight:
		return x + 1, y
	case MoveUp:
		return x, y - 1
	case MoveDown:
		return x, y + 1
	default:
		panic("Unknown move " + m.String())
	}
}
func (m Move) Reverse() Move {
	switch(m){
	case MoveLeft:
		return MoveRight
	case MoveRight:
		return MoveLeft
	case MoveUp:
		return MoveDown
	case MoveDown:
		return MoveUp
	default:
		panic("Unknown move " + m.String())
	}
}

func (m Move) String() string {
	var mv_str = ""
	if (m == NoMove) {
		mv_str = "Ø"
	}
	if (m & MoveUp != NoMove) {
		mv_str += "⇧"
	}
	if (m & MoveDown != NoMove) {
		mv_str += "⇩"
	}
	if (m & MoveLeft != NoMove) {
		mv_str += "⇦"
	}
	if (m & MoveRight != NoMove) {
		mv_str += "⇨"
	}
	return mv_str
}

func (pm PieceMove) String() string {
	return fmt.Sprintf("%d=%s", pm.piece, pm.mv)
}

func (pm PossMoves) String() string {
	movestr := ""
	for i, move := range (pm) {
		movestr += fmt.Sprintf("%d=%s ", i, move)
	}
	return movestr
}

func (b Board) FillMap() FillMap {
	result := fmEmpty
	for i := 0; i < nPieces; i++ {
		x, y := b.PiecePos(i)
		result |= pTypes[i].FillMap(x, y)
	}
	return result
}

func (b Board) PiecePos(piece int) (x int, y int) {
	ymap := uint64((1 << yBits) - 1)
	xmap := uint64((1 << xBits) - 1)

	offset := uint(piece * bitsPerPiece)

	x = int((uint64(b) >> (offset + yBits)) & xmap)

	y = int((uint64(b) >> offset) & ymap)
	return
}

func (pm PossMove) AddMove(m Move) PossMove {
	return PossMove(uint8(pm) | uint8(m))
}

func (pm PossMove) String() string {
	return Move(pm).String()
}
func (pm PossMove) CanMove(m Move) bool {
	return pm & PossMove(m) != noMovePoss
}

func (fm FillMap) Remove(other FillMap) FillMap {
	return fm ^ other;
}

func (fm FillMap) CanPlace(other FillMap) bool {
	return fm & other == fmEmpty;
}

// PossibleMoves Gets the possible moves for this board
func (b Board) PossibleMoves() PossMoves {
	moves := [nPieces]PossMove{}

	basemap := b.FillMap();

	for i := 0; i < nPieces; i++ {
		x, y := b.PiecePos(i)
		pieceType := pTypes[i]

		piecemap := pieceType.FillMap(x, y)
		piecePasemap := basemap.Remove(piecemap)
		var posMoves PossMove = 0

		for _, m := range (AvailableMoves) {
			nx, ny := m.Transform(x, y)
			pieceFM := pieceType.FillMap(nx, ny)
			if pieceFM != fmInvalid && piecePasemap.CanPlace(pieceFM) {
				posMoves = posMoves.AddMove(m)
			}
		}
		moves[i] = posMoves
	}
	return moves
}

func locationFillMap(x int, y int) FillMap {
	return FillMap(1 << uint32((y * xMax) + x))
}

func (fm FillMap) IsValid() bool {
	return fm != fmInvalid;
}

func (pt PieceType) FillMap(x int, y int) FillMap {

	if x < 0 || y < 0 || x >= xMax || y >= yMax {
		return fmInvalid
	}
	switch pt {
	case tSmall :
		return 1 << uint32((y * xMax) + x)
	case tHoriz :
		if (x >= (xMax - 1 )) {
			return fmInvalid
		}
		return 1 << uint32((y * xMax) + x) | 1 << uint32((y * xMax) + x + 1)
	case tVerti :
		if (y >= (yMax - 1)) {
			return fmInvalid
		}
		return 1 << uint32((y * xMax) + x) | 1 << uint32(((y + 1) * xMax) + x)
	}

	if x >= (xMax - 1 ) || y >= (yMax - 1) {
		return fmInvalid;
	}
	return (1 << uint32((y * xMax) + x)) |
	  (1 << uint32(((y + 1) * xMax) + x)) |
	  (1 << uint32((y * xMax) + x + 1)) |
	  (1 << uint32(((y + 1) * xMax) + x + 1));


}

const (
	bm3 = uint64(7)
	bm2 = uint64(3)
)

func sort4(b Board, from int) Board {
	sorts := []struct{ a int; b int }{{0, 1}, {2, 3}, {0, 3}, {1, 2}, {0, 1}, {2, 3}}
	for _, st := range (sorts) {
		ax, ay := b.PiecePos(st.a + from);
		bx, by := b.PiecePos(st.b + from);
		cmp := (ay * xMax + ax) - (by * xMax + bx);
		if cmp < 0 {
			b = b.placePieceInternal(st.a + from, bx, by).
			  placePieceInternal(st.b + from, ax, ay)
		}
	}
	return b
}

func (b Board) identity() Board {
	return sort4(sort4(b, 0), 5)
}

func (b Board) placePieceInternal(piece int, x int, y int) Board {
	piecemask := ^(uint64(31) << uint(piece * bitsPerPiece))
	pieceval := (((uint64(x) & bm2) << yBits) | (uint64(y) & bm3)) << (uint32(piece) * bitsPerPiece)
	return Board((uint64(b) & piecemask) | pieceval)
}

func (b Board) PlacePiece(piece int, x int, y int) Board {

	return b.placePieceInternal(piece,x,y).identity();
}

func (b Board) MovePiece(piece int, m Move) Board {
	x, y := b.PiecePos(piece)

	nx, ny := m.Transform(x, y)
	return b.PlacePiece(piece, nx, ny)

}

func MakeBoard(locs []loc) Board {
	var bd = Board(0)
	for idx, pos := range (locs) {
		bd = bd.placePieceInternal(idx, pos.x, pos.y)
	}
	return bd.identity()
}

func makeFillmap(locs []loc) FillMap {
	result := fmEmpty
	for i, pos := range (locs) {
		result |= pTypes[i].FillMap(pos.x, pos.y)
	}
	return result
}

func (fm FillMap) Get(x int, y int) bool {
	return (fm & (1 << uint32((y * xMax) + x))) != fmEmpty;
}

func (fm FillMap) String() string {
	var val = ""

	for y := 0; y < yMax; y++ {
		val += "|"
		for x := 0; x < xMax; x++ {
			if (fm.Get(x, y)) {
				val += "x"

			} else {
				val += " "
			}
		}
		val += "|\n"
	}
	return val
}

func (b Board )  String() string {
	bd := [5][4]string{}
	typeGlyphs := []string{"▣", "▭", "▯", "▢"}

	for i := 0; i < nPieces; i++ {
		x, y := b.PiecePos(i)
		t := pTypes[i]
		pfm := t.FillMap(x, y)
		for dx := 0; dx < xMax; dx ++ {
			for dy := 0; dy < yMax; dy ++ {
				if (pfm.Get(dx, dy)) {
					bd[dy][dx] = typeGlyphs[uint(t)]
				}
			}
		}
	}

	str := "\n"
	for dy := 0; dy < yMax; dy ++ {
		str += "|"
		for dx := 0; dx < xMax; dx ++ {
			v := bd[dy][dx]
			if (v == "") {
				v = " "
			}
			str += v
		}
		str += "|\n"
	}
	for i := 0; i < nPieces; i++ {
		x, y := b.PiecePos(i)
		str += fmt.Sprintf("%d %s =(%d,%d) ", i, typeGlyphs[pTypes[i]], x, y)
	}

	return str + "\n"
}

var default_start = []loc{{0, 4}, {1, 4}, {2, 4}, {3, 4},
	{1, 3},
	{0, 0}, {0, 2}, {3, 0}, {3, 2},
	{1, 1},
}

type BoardGraph struct {
	start       Board
	transitions map[Board]map[PieceMove]Board
}

func PopulateGraph(start Board) *BoardGraph {

	seen := make(map[Board]map[PieceMove]Board)
	next := []Board{start}

	for {
		if (len(next) == 0) {
			break
		}
		cur := next[0]
		//	fmt.Printf("cur %s\n" ,cur)
		next = next[1:]
		if seen[cur] == nil {
			seen[cur] = make(map[PieceMove]Board)
		} else {
			continue
		}

		moves := cur.PossibleMoves()
		for i := 0; i < nPieces; i++ {
			for _, mv := range (AvailableMoves) {
				if (moves[i].CanMove(mv)) {
					board := cur.MovePiece(i, mv).Identity();
					seen[cur][PieceMove{i, mv}] = board
					if seen[board] == nil {
						next = append(next, board)
					}
				}
			}
		}
	}

	return &BoardGraph{start, seen}
}

func (b BoardGraph) BoardByArcsBreadthFirst(visit func(from Board, via PieceMove, to Board)) {
	seen := make(map[Board]map[PieceMove]bool)
	next := []Board{b.start}
	for {
		if (len(next) == 0) {
			break
		}

		cur := next[0]

		next = next[1:]

		if seen[cur] == nil {
			seen[cur] = make(map[PieceMove]bool)
		}

		links := b.transitions[cur]
		if links != nil {
			for pm, destBoard := range (links) {
				if (!seen[cur][pm]) {
					visit(cur, pm, destBoard)
					seen[cur][pm] = true
					next = append(next, destBoard)
				}
			}
		}
	}
}
