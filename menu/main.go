package menu

import "fmt"

// Row and Column Start from the top left corner at 0, 0
type Cursor struct {
	Row int
	Col int
}

func Init() Cursor {
	fmt.Print("\033c")
	return Cursor{Row: 0, Col: 0}
}

// Moves the cursor up x rows
func (c Cursor) Up(x int) {
	c.Row -= x
	fmt.Printf("\033[%dA", x)
}

// Moves the cursor down x rows
func (c Cursor) Down(x int) {
	c.Row += x
	fmt.Printf("\033[%dB", x)
}

// Moves the cursor left x rows
func (c Cursor) Left(x int) {
	c.Col -= x
	fmt.Printf("\033[%dC", x)
}

// Moves the cursor left x rows
func (c Cursor) Right(x int) {
	c.Col += x
	fmt.Printf("\033[%dD", x)
}

// Moves cursor to Row, Col NOT WORKING ATM
func (c Cursor) GoTo(row, col int) {
	fmt.Printf("ESC[%d;%df", row, col)
}
