package main

// #cgo LDFLAGS: -L. -lTWAI
// #include <stdbool.h>
// long *nextBestMoves();
// void playMove(long move);
// long *renderFrame();
// void lockGame();
// bool resetGame();
import "C"
import (
	"fmt"
	"os"
	"reflect"
	"time"
	"unsafe"
	"github.com/gdamore/tcell/v2"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()
	index := 0
	var moves []int
	defer s.Fini()
	C.resetGame()
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Nanosecond):
		}
		s.Clear()

		if moves == nil || len(moves) == index {
			moves = nextBestMoves()
			index = 0
			if len(moves) == 0 {
				C.lockGame()
				moves = nextBestMoves()
			}
			if len(moves) == 0 {
				s.Fini()
			}
		}
		playNextMove(moves[index])
		index++
		frame := renderFrame()
		for y := 0; y < frame.hight; y++ {
			for x := 0; x < frame.width; x++ {
				if frame.board[(y*frame.width)+x] != 0 {
					setPoint(s, x*2, y)
					setPoint(s, x*2+1, y)
				}
			}
		}
		s.Show()
	}

}

func noUi() {
	for {
		moves := nextBestMoves()
		for _, move := range moves {
			fmt.Println("Playing move:", move)
			playNextMove(move)
		}
		C.lockGame()
	}
}

func playNextMove(move int) {
	nextMove := C.long(move)
	C.playMove(nextMove)
}

func nextBestMoves() []int {
	moves := C.nextBestMoves()
	size := int(*moves)
	p := uintptr(unsafe.Pointer(moves)) + unsafe.Sizeof(size)
	sh := &reflect.SliceHeader{Data: p, Len: size, Cap: size}
	return *(*[]int)(unsafe.Pointer(sh))
}

type Frame struct {
	board []int
	width int
	hight int
}

func renderFrame() Frame {
	moves := C.renderFrame()
	size := int(*moves) - 2
	hight := uintptr(unsafe.Pointer(moves)) + unsafe.Sizeof(size)
	width := uintptr(unsafe.Pointer(moves)) + unsafe.Sizeof(size)*2
	p := uintptr(unsafe.Pointer(moves)) + unsafe.Sizeof(size)*3
	sh := &reflect.SliceHeader{Data: p, Len: size, Cap: size}

	return Frame{
		board: *(*[]int)(unsafe.Pointer(sh)),
		width: *(*int)(unsafe.Pointer(width)),
		hight: *(*int)(unsafe.Pointer(hight)),
	}
}

func setPoint(s tcell.Screen, x, y int) {
	st := tcell.StyleDefault
	st = st.Background(tcell.NewHexColor(0xff))
	s.SetContent(x, y, ' ', nil, st)
}
