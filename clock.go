package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/bclock/bcd"
)

const Box = '\u2588'

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
var blankStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorLightGray)

func drawBin(s tcell.Screen, x int, y int, bin [4]bool, onStyle tcell.Style, offStyle tcell.Style) {
	for i := 0; i < len(bin); i++ {

		// Vertical position
		pos := y + (len(bin)-1)*2 - (i * 2)

		if bin[i] {
			s.SetContent(x, pos, Box, nil, onStyle)
			s.SetContent(x+1, pos, Box, nil, onStyle)
		} else {
			s.SetContent(x, pos, Box, nil, offStyle)
			s.SetContent(x+1, pos, Box, nil, offStyle)
		}
	}
}

// Draws a binary representation of a number onto the given screen
func drawTime(s tcell.Screen, t time.Time, onStyle tcell.Style, offStyle tcell.Style) {
	b := bcd.ToBCD(t)

	w, h := s.Size()

	// TODO: fix magic numbers
	x := (w - 32) / 2
	y := (h - 7) / 2

	// Draw Markers
	for i := 0; i < 4; i++ {
		// Generate our rune by adding a number
		p := '0' + int32(math.Pow(2, float64(i)))
		pos := y + (len(b.Hours.Ones)-1)*2 - (i * 2)
		s.SetContent(x, pos, p, nil, offStyle)
	}

	// Offset from left numbers
	o := x + 3

	// Draw individual bits
	drawBin(s, o, y, b.Hours.Tens, onStyle, offStyle)
	drawBin(s, o+4, y, b.Hours.Ones, onStyle, offStyle)
	drawBin(s, o+10, y, b.Minutes.Tens, onStyle, offStyle)
	drawBin(s, o+14, y, b.Minutes.Ones, onStyle, offStyle)
	drawBin(s, o+20, y, b.Seconds.Tens, onStyle, offStyle)
	drawBin(s, o+24, y, b.Seconds.Ones, onStyle, offStyle)
}

func watcher(s tcell.Screen, events <-chan tcell.Event) {
	for {
		e := <-events

		switch ev := e.(type) {
		case *tcell.EventResize:
			s.Clear()
			drawTime(s, time.Now(), defStyle, blankStyle)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
		}
	}
}

func main() {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.SetStyle(defStyle)
	s.Clear()

	drawTime(s, time.Now(), defStyle, blankStyle)

	events := make(chan tcell.Event, 1)
	quit := make(chan struct{}, 1)

	go s.ChannelEvents(events, quit)
	go watcher(s, events)

	// event loop
	for t := range time.Tick(time.Second) {
		s.Clear()
		drawTime(s, t, defStyle, blankStyle)
		s.Show()
	}
}
