package main

import (
	"fmt"

	"gopherit/menu"
)

func main() {
	c := menu.Init()
	c.Down(20)
	fmt.Println(c.Row, c.Col)
	fmt.Println("Hello World!")
	for {
	}
}
