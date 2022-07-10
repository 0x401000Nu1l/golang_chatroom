package server

import (
	"golang_chatroom/logic"
	"net/http"
)

func RegisterHandle() {
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
