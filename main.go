package main

import (
	"fmt"
	e "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"os"
)

const (
	CLOSE = iota
	CREATE
	CONNECT
)

func update(screen *e.Image) error {
	img, _, _ := ebitenutil.NewImageFromFile("./images/BackGround.png", e.FilterDefault)
	screen.DrawImage(img, nil)
	return nil
}

func main() {
	greetingsMainMenu()
	c := getConsoleStartChoose()
	switch c {
	case CLOSE:
		os.Exit(1)
	case CREATE:
		createGame()
	case CONNECT:
		findGames()
	}
}

func createGame() {
	e.Run(update, 800, 600, 1, "Snakes")
}

func findGames() {

}

func getConsoleStartChoose() int {
	for {
		c := -1
		fmt.Fscan(os.Stdin, &c)
		if c > 3 || c < 0 {
			fmt.Println("Wrong input")
			continue
		}
		return c
	}
}

func greetingsMainMenu() {
	fmt.Println("Hello User!")
	fmt.Println("Choose what you want:")
	fmt.Println("1. Create a new game")
	fmt.Println("2. Connect to the game")
	fmt.Println("0. Close the game")
}
