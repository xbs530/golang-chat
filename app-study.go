package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/websocket"
	"xbs530.com/app-study/library/chat/handler"
)



func main() {

	http.Handle("/chat",websocket.Handler(handler.Session))
	//http.Handle("/",http.FileServer(http.Dir("D:/")))

	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		fmt.Printf("error: %v",err)
	}

}

