package main

import (
	"log"
	"os"
)

func main() {
	logFile := "/tmp/a.log"
	w, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(w)
	log.Println("go-misc/logoutput: ログ出力先をファイルに変更")
}
