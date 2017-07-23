package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/websocket"
	"xbs530.com/app-study/library/chat/handler"
	"os"
)



func main() {

	http.Handle("/chat",websocket.Handler(handler.Session))
	//http.Handle("/",http.FileServer(http.Dir("D:/")))

	bind_port := "8080"
	if len(os.Args)>1 {
		bind_port=os.Args[1]
	}

	err := http.ListenAndServe(":"+bind_port,nil)
	if err != nil {
		fmt.Printf("error: %v",err)
	}

}

