package main

//#include <stdio.h>
//#include <stdlib.h>
//void PrintString(char* str){
//	printf("%s\n", str);
//}
import "C"
import "unsafe"

func main() {
	a := C.CString("This is from Golang")
	C.PrintString(a)
	C.free(unsafe.Pointer(a))
}
