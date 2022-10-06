package main

import (
	"fmt"
	e "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"net"
	"os"
	"time"
)

const (
	CLOSE = iota
	CREATE
	CONNECT
	srvAddr         = "239.192.0.4:9192"
	maxDatagramSize = 8192
)
const AnnouncementMsg = "AnnouncementMsg"

type CurrentGames struct {
	gamesSrc []*net.UDPAddr
}

func update(screen *e.Image) error {
	img, _, _ := ebitenutil.NewImageFromFile("./images/BackGround.png", e.FilterDefault)
	screen.DrawImage(img, nil)
	return nil
}

func main() {
	printMainMenu()
	c := getConsoleStartChoose()
	switch c {
	case CLOSE:
		os.Exit(CLOSE)
	case CREATE:
		go UDPSender(srvAddr)
		createGame()
	case CONNECT:
		findGames()
	}
}

func UDPSender(adr string) {
	MulticastUDPSender(adr)
}

func MulticastUDPSender(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			c.Write([]byte(AnnouncementMsg))
			time.Sleep(time.Second)
		}
	}()
}

func MulticastUDPListener(a string, h func(*net.UDPAddr, int, []byte, CurrentGames), cg CurrentGames) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	err = l.SetReadBuffer(maxDatagramSize)
	if err != nil {
		return
	}
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b, cg)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte, cg CurrentGames) {
	message := string(b[:n])
	log.Println(string(b[:n]) + " " + src.String())
	if message == AnnouncementMsg {
		for _, addr := range cg.gamesSrc {
			if addr == src {
				return
			}
		}
		cg.gamesSrc = append(cg.gamesSrc, src)
	}
}

func createGame() {
	e.Run(update, 800, 600, 1, "Snakes")
	time.Sleep(100 * time.Second)
}

func findGames() {
	cg := CurrentGames{gamesSrc: []*net.UDPAddr{}}
	go MulticastUDPListener(srvAddr, msgHandler, cg)
	printWaitGames()
	if len(cg.gamesSrc) == 0 {
		fmt.Println("NO GAMES")
		return
	}
	printChooseGameMenu(cg)
	gn := getConsoleGameChoose(len(cg.gamesSrc))
	connectToTheGame(cg.gamesSrc[gn])
}

func connectToTheGame(game *net.UDPAddr) {
	fmt.Println("Trying connect to: " + game.IP.String())
	time.Sleep(100 * time.Second)
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

func getConsoleGameChoose(r int) int {
	for {
		c := -1
		fmt.Fscan(os.Stdin, &c)
		if c > r || c < 0 {
			fmt.Println("Wrong input")
			continue
		}
		return c - 1
	}
}

func printMainMenu() {
	fmt.Println("Hello User!")
	fmt.Println("Choose what you want:")
	fmt.Println("1. Create a new game")
	fmt.Println("2. Connect to the game")
	fmt.Println("0. Close the game")
}

func printChooseGameMenu(cg CurrentGames) {
	fmt.Println("Choose a game")
	for i, addr := range cg.gamesSrc {
		fmt.Println(string(i+1) + "." + " " + addr.IP.String())
	}
}

func printWaitGames() {
	fmt.Println("finding games...")
	time.Sleep(time.Second / 2)
	fmt.Println("35%...")
	time.Sleep(time.Second / 3)
	fmt.Println("75%...")
	time.Sleep(time.Second / 3)
	fmt.Println("99%...")
	time.Sleep(time.Second)
	fmt.Println("done!")
}
