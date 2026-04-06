package menu

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/x/term"
	"github.com/eiannone/keyboard"
)

// Row and Column Start from the top left corner at 1, 1
type Menu struct {
	Row    int
	Col    int
	Pane   int
	Height int
	Width  int
	Panes  []Pane
}

type Pane struct {
	Lines []string
}

func clear() {
	fmt.Print("\033c")
}

func Init() *Menu {
	stdout := os.Stdout.Fd()
	w, h, err := term.GetSize(stdout)
	if err != nil {
		h = 24
		w = 80
	}

	clear()

	panes := make([]Pane, 2)
	for i := range panes {
		panes[i] = Pane{Lines: make([]string, h)}
	}

	//fmt.Printf("\033[%dB", (h-1)/2)
	return &Menu{Row: 1, Col: 1, Pane: 0, Height: h, Width: w, Panes: panes}
}

// Moves the cursor up x rows
func (m *Menu) up(x int) {
	if m.Row-x > 0 {
		m.Row -= x
	} else {
		m.Row = 0
	}
	fmt.Printf("\033[%dA", x)
}

// Moves the cursor down x rows
func (m *Menu) down(x int) {
	if m.Row+x < m.Height {
		m.Row += x
	} else {
		m.Row = m.Height
	}
	fmt.Printf("\033[%dB", x)
}

// Moves the cursor left x columns
func (m *Menu) left(x int) {
	if m.Col-x > 0 {
		m.Col -= x
	} else {
		m.Col = 0
	}
	fmt.Printf("\033[%dD", x)
}

// Moves the cursor right x columns
func (m *Menu) right(x int) {
	if m.Col+x < m.Width {
		m.Col += x
	} else {
		m.Col = m.Width
	}
	fmt.Printf("\033[%dC", x)
}

func (m *Menu) nextPane() {
	if m.Pane+1 < len(m.Panes)-1 {
		m.Pane++
	} else {
		m.Pane = len(m.Panes) - 1
	}
}

func (m *Menu) prevPane() {
	if m.Pane-1 > 0 {
		m.Pane--
	} else {
		m.Pane = 0
	}
}

func (m *Menu) Listen() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyEsc {
			clear()
			//fmt.Println(m)
			break
		}
		if key == keyboard.KeyArrowDown {
			m.down(1)
		}
		if key == keyboard.KeyArrowUp {
			m.up(1)
		}
		if key == keyboard.KeyArrowRight {
			m.nextPane()
			m.right(40)
		}
		if key == keyboard.KeyArrowLeft {
			m.prevPane()
			m.left(40)
		}
		if key == keyboard.KeySpace {
			m.write(strconv.Itoa(m.Row) + "," + strconv.Itoa(m.Pane))
			m.draw()
		}
	}
}

// func (m Menu) MainLoop() {
// 	for {
// 		if listen(&m) == 1 {
// 			break
// 		}
// 	}
// }

func (m Menu) printr(text string) {
	fmt.Print(text + "\r")
	m.Col = 0
}

func (m Menu) print(text string) {
	fmt.Print(text)
	m.Col += len(text)
}

func (m Menu) write(text string) {
	m.Panes[m.Pane].Lines[m.Row] = text
}

func (m Menu) draw() {
	m.goTo(1, 1)
	clear()
	for i := 1; i < m.Height-1; i++ {
		fmt.Print(m.Panes[0].Lines[i])
		fmt.Printf("\033[%dC", 40-len(m.Panes[0].Lines[i]))
		// Return carriage on last line
		fmt.Println(m.Panes[len(m.Panes)-1].Lines[i])
	}
	m.goTo(m.Row, m.Col)
}

// Moves cursor to Row, Col NOT WORKING ATM
func (m Menu) goTo(row, col int) {
	fmt.Printf("\033[%d;%df", row, col)
}
