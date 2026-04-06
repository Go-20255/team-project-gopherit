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

	// Initialize with cursor centered
	return &Menu{Row: h / 2, Col: 1, Pane: 0, Height: h, Width: w, Panes: panes}
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

// Navigates to the next pane, allocates a new one
// If one is not present
func (m *Menu) nextPane() {
	if m.Pane+1 < len(m.Panes)-1 {
		m.Pane++
	} else {
		m.Panes = append(m.Panes, Pane{Lines: make([]string, m.Height)})
		m.Pane++
	}
}

// Navigates to the previous pane
func (m *Menu) prevPane() {
	if m.Pane-1 > 0 {
		m.Pane--
	} else {
		m.Pane = 0
	}
}

// Gets a list of strings representing the
// contents of the next directory over
func readDir(path string) []string {
	dir, err := os.ReadDir(path)
	if err != nil {
		return make([]string, 0)
	}
	files := make([]string, 0)

	for _, entry := range dir {
		files = append(files, entry.Name())
	}
	return files
}

// Sets the initial information in the menu before
// control is given to the user
// TODO: Same with run, lets accept a starting dir
func (m *Menu) prerun() {
	// Get initial directory
	home, _ := os.UserHomeDir()
	files := readDir(home + "/classes")
	m.up(len(files) / 2)
	m.batchWrite(files, 0)

	nextDir := home + "/classes/" + files[0]
	files = readDir(nextDir)
	m.batchWrite(files, 1)

	m.draw()
}

// The meat and potatoes. Run does key listening and dispatches
// commands based on them. Press esc to quit
// TODO, have run accept a starting directory
func (m *Menu) Run() {
	m.prerun()

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
		// TODO: When arrow keys are pressed, make a query
		// to getDir to look into the next directory
		// m.Panes[m.Pane].Lines[m.Row] will be helpful
		if key == keyboard.KeyArrowDown {
			m.down(1)
			m.draw()
		}
		if key == keyboard.KeyArrowUp {
			m.up(1)
			m.draw()
		}
		if key == keyboard.KeyArrowRight {
			m.nextPane()
			m.draw()
		}
		if key == keyboard.KeyArrowLeft {
			m.prevPane()
			m.draw()
		}
		if key == keyboard.KeySpace {
			m.write(strconv.Itoa(m.Row)+","+strconv.Itoa(m.Pane), m.Row, m.Pane)
			m.draw()
		}
	}
}

// Write a single line of text to a line on the specified pane
// Mostly a helper function
func (m *Menu) write(text string, row int, pane int) {
	m.Panes[pane].Lines[row] = text
}

// Writes a list of strings to the described pane
// Lines will be start from m.Row
func (m *Menu) batchWrite(lines []string, pane int) {
	for i, line := range lines {
		m.write(line, m.Row+i, pane) // Change this to format better
	}
}

// Draws the Menu struct's information to the screen
// TODO: Make the line at m.Row a different color
// (ANSI escape sequences will be your friend here)
func (m Menu) draw() {
	m.goTo(1, 1)
	clear()
	for i := 1; i < m.Height-1; i++ {
		fmt.Print(m.Panes[m.Pane].Lines[i])
		fmt.Printf("\033[%dC|", 39-len(m.Panes[m.Pane].Lines[i]))
		// Return carriage on last line
		if m.Pane+1 < len(m.Panes) {
			fmt.Println(m.Panes[m.Pane+1].Lines[i])
		} else {
			fmt.Println("")
		}
	}
	m.goTo(m.Row, m.Col)
}

// Moves cursor to Row, Col
// goTo does *NOT* update the cursor position in m
// This is by design, and is used mostly for the
// draw function
func (m Menu) goTo(row, col int) {
	fmt.Printf("\033[%d;%df", row, col)
}
