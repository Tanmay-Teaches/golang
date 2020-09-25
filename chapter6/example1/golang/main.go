package main

// #cgo LDFLAGS: -L. -lTWAI
// long *nextBestMoves();
// void playMove(int move);
// long *renderFrame();
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	moves := C.nextBestMoves()

	size := int(*moves)
	p := uintptr(unsafe.Pointer(moves)) + unsafe.Sizeof(size)
	sh := &reflect.SliceHeader{
		Data: p,
		Len:  size,
		Cap:  size,
	}
	data := *(*[]int)(unsafe.Pointer(sh))
	fmt.Println(size, data)
}
