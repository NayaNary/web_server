package main

import (
	"fmt"
	"log"
	"os"

	"test.task/src/web"
)

func main() {
	
	webRouts := web.New()

	srv := webRouts.StartServer()

	var data string

	for {
		fmt.Println("Введите 'stop' для закрытия приложения")
		n, err := fmt.Fscan(os.Stdin, &data)
		if n == 0 {
			log.Println(err.Error())
		} else {
			if data == "stop" {
				srv.Close()
				webRouts.ProcData.Timer.Stop()
				webRouts.ProcData.WritePages()
				break
			}
		}
	}
}
