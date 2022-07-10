package server

import (
	"log"
	"net/http"

	"golang_chatroom/logic"

	"github.com/gorilla/websocket"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 从客户端接受 WebSocket 握手，并将连接升级到 WebSocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error: ", err)
		return
	}
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal:", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：2-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在： ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("改昵称已存在"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}
	userHasToken := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	go userHasToken.SendMessage(req.Context())

	userHasToken.MessageChannel <- logic.NewWelcomeMessage(userHasToken)

	//避免token泄露
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	//给所有用户告知新用户的到来
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 将该用户加入广播器的用列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	// 接收用户消息
	err = user.ReceiveMessage(req.Context())

	// 用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcast.Broadcaster(msg)
	log.Println("user: ", nickname, "leave chat")

	//根据读取时的错误执行不同的close
	if err == nil {
		conn.Closr(websocket.StatusNoStatusRcvd, "")
	} else {
		log.Println("read from client error :", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}

}
