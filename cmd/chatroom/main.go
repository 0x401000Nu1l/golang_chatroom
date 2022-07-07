package main

import (
	"fmt"
	"golang_chatroom/global"
	"golang_chatroom/server"
	"log"
	"net/http"
)

var (
	addr   = ":https://github.com/0x401000Nu1l/golang_chatroom"
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
