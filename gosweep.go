package gosweep

import (
	"fmt"
	"math/rand"
	"time"
)

// Minefield represents minefield
type Minefield struct {
	field  [][]Cell
	width  int
	height int
	mines  int
	flags  int
}

// Cell represents the cell on the minefield
type Cell struct {
	t CellType
	s CellState
}

// CellType represents type of the cell on the minefield
type CellType = int

// Cell types
const (
	TypeEmpty CellType = 0
	Type1     CellType = 1
	Type2     CellType = 2
	Type3     CellType = 3
	Type4     CellType = 4
	Type5     CellType = 5
	Type6     CellType = 6
	Type7     CellType = 7
	Type8     CellType = 8
	TypeMine  CellType = 9
)

// CellState represents state of the cell on the minefield
type CellState = int

// Cell states
const (
	StateClosed  CellState = 0
	StateFlagged CellState = 1
	StateOpened  CellState = 2
)

// NewMinefield creates new minefield
func NewMinefield(width, height, mines int) Minefield {
	minefield := Minefield{
		width:  width,
		height: height,
		mines:  mines,
	}

	minefield.generateField()

	return minefield
}

// GetField returns 2d array that represents minefield
func (m *Minefield) GetField() [][]Cell {
	return m.field
}

// GetWidth returns number of columns of the minefield grid
func (m *Minefield) GetWidth() int {
	return m.width
}

// GetHeigth returns number of rows of the minefield grid
func (m *Minefield) GetHeigth() int {
	return m.height
}

// GetMines returns number of mines on the minefield grid
func (m *Minefield) GetMines() int {
	return m.mines
}

// GetFlags returns number of flags on the minefield grid
func (m *Minefield) GetFlags() int {
	return m.flags
}

// Print prints minefield
func (m *Minefield) Print() {
	typeChars := []string{
		TypeEmpty: " ",
		Type1:     "1",
		Type2:     "2",
		Type3:     "3",
		Type4:     "4",
		Type5:     "5",
		Type6:     "6",
		Type7:     "7",
		Type8:     "8",
		TypeMine:  "*",
	}

	stateChars := map[int]string{
		StateClosed:  "-",
		StateFlagged: "F",
	}

	for row := 0; row < m.height; row++ {
		for col := 0; col < m.width; col++ {
			cell := &m.field[row][col]
			if val, ok := stateChars[cell.s]; ok {
				fmt.Printf("%s ", val)
			} else {
				fmt.Printf("%s ", typeChars[cell.t])
			}
		}

		fmt.Println()
	}
}

// Open opens cell on the minefield
func (m *Minefield) Open(row, col int) {
	if !m.isInBounds(row, col) {
		return
	}

	cell := &m.field[row][col]
	if cell.s == StateOpened || cell.s == StateFlagged {
		return
	}

	if cell.t == TypeMine {
		// Game over
		return
	}

	m.floodFillOpen(row, col)
}

// ToggleFlag toggles state of the cell between flagged and closed
func (m *Minefield) ToggleFlag(row, col int) {
	if !m.isInBounds(row, col) {
		return
	}

	cell := &m.field[row][col]
	if cell.s == StateClosed {
		cell.s = StateFlagged
		m.flags++
	} else if cell.s == StateFlagged {
		cell.s = StateClosed
		m.flags--
	}
}

func (m *Minefield) floodFillOpen(row, col int) {
	if !m.isInBounds(row, col) {
		return
	}

	cell := &m.field[row][col]
	if cell.t == TypeMine || cell.s == StateOpened {
		return
	}

	cell.s = StateOpened
	if cell.t != TypeEmpty {
		return
	}

	m.floodFillOpen(row+1, col+1)
	m.floodFillOpen(row-1, col-1)
	m.floodFillOpen(row+1, col-1)
	m.floodFillOpen(row-1, col+1)
	m.floodFillOpen(row+1, col)
	m.floodFillOpen(row-1, col)
	m.floodFillOpen(row, col+1)
	m.floodFillOpen(row, col-1)
}

func (m *Minefield) generateField() {
	m.field = make([][]Cell, m.height)
	for row := 0; row < m.height; row++ {
		m.field[row] = make([]Cell, m.width)
	}

	rand.Seed(time.Now().Unix())

	minesSet := 0
	for minesSet < m.mines {
		// TODO: use crypto/rand to generate minefield
		randRow := rand.Intn(m.height)
		randCol := rand.Intn(m.width)

		cell := &m.field[randRow][randCol]
		if cell.t == TypeEmpty {
			cell.t = TypeMine
			cell.s = StateClosed
			minesSet++
		}
	}

	for row := 0; row < m.height; row++ {
		for col := 0; col < m.width; col++ {
			cell := &m.field[row][col]
			if cell.t == TypeEmpty {
				cell.t = m.getHint(row, col)
			}
		}
	}
}

func (m *Minefield) getHint(row, col int) CellType {
	var result CellType
	b2i := map[bool]int{true: 1, false: 0}
	result += b2i[m.isMine(row+1, col+1)]
	result += b2i[m.isMine(row+1, col-1)]
	result += b2i[m.isMine(row-1, col+1)]
	result += b2i[m.isMine(row-1, col-1)]
	result += b2i[m.isMine(row, col+1)]
	result += b2i[m.isMine(row, col-1)]
	result += b2i[m.isMine(row+1, col)]
	result += b2i[m.isMine(row-1, col)]

	return result
}

func (m *Minefield) isInBounds(row, col int) bool {
	return col >= 0 && col < m.width && row >= 0 && row < m.height
}

func (m *Minefield) isMine(row, col int) bool {
	return m.isInBounds(row, col) && m.field[row][col].t == TypeMine
}
