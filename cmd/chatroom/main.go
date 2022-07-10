package main

import (
	"fmt"
	"golang_chatroom/global"
	"golang_chatroom/server"
	"log"
	"net/http"
)

var (
	addr   = ":2020" //端口
	banner = `
    ____              _____ 
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |
ChatRoom，start ：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner, addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
