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
		os.Exit(CLOSE)
	case CREATE:
		go UDPSender(srvAddr)
		createGame()
	case CONNECT:
		go MulticastUDPListener(srvAddr, msgHandler)
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

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(string(b[:n]) + " " + src.String())
}

func MulticastUDPListener(a string, h func(*net.UDPAddr, int, []byte)) {
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
		h(src, n, b)
	}
}

func createGame() {
	e.Run(update, 800, 600, 1, "Snakes")
	time.Sleep(100 * time.Second)
}

func findGames() {
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

func greetingsMainMenu() {
	fmt.Println("Hello User!")
	fmt.Println("Choose what you want:")
	fmt.Println("1. Create a new game")
	fmt.Println("2. Connect to the game")
	fmt.Println("0. Close the game")
}
