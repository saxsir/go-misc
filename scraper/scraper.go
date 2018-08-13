package main

import (
	"time"

	"github.com/sclevine/agouti"
)

func main() {

	url := "http://example.com"

	for i := 0; i < 2; i++ {
		page, err := agouti.NewPage("http://localhost:9515")
		if err != nil {
			panic(err)
		}
		defer page.Destroy()

		if err := page.Navigate(url); err != nil {
			panic(err)
		}

		time.Sleep(3 * time.Second)
	}
}
