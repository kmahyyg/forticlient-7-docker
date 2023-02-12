package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("password:\n")
	var v1 string
	var v2 string
	_, _ = fmt.Scanln(&v1)
	time.Sleep(3 * time.Second)
	fmt.Printf("Confirm (y/n) [default=n]:Confirm (y/n) [default=n]:")
	_, _ = fmt.Scanln(&v2)
	fmt.Println("result: ", v1, v2)
	fmt.Println("args: ", os.Args)
	time.Sleep(15 * time.Second)
	fmt.Println("random string to test output")
	time.Sleep(9999 * time.Second)
}
