package main

import "fmt"

func main() {
	go InitiateBot()
	server()
}

func ErrCheck(err error) {
	if err != nil {
		fmt.Println(err);
		return
	}
}