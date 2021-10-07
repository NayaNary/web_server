package main

import (
	"context"
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
				webRouts.ProcData.Timer.Stop()
				
				if len(webRouts.ProcData.Pages[webRouts.ProcData.LastId])>0{
					webRouts.ProcData.WritePages()
				}
	
				err:= srv.Shutdown(context.Background())
				if err !=nil{
					log.Println(err.Error())
				}
				break
			}
		}
	}
}
