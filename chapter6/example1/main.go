package main
//#include <stdio.h>
//void PrintHelloWorld(){
//	printf("Hello, World!\n");
//}
import "C"

func main() {
	C.PrintHelloWorld()
}
