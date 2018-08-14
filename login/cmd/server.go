package main

import "github.com/saxsir/talks/2018/treasure-go/login"

func main() {
	s := login.NewServer()
	s.Init()
	s.Run(":8888")
}
