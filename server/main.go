package main

import (
	"fmt"
	"main/api/chatpb"
	"main/api/userpb"
	"main/server/chat"
	"main/server/user"
	"net"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()

	userServ := user.NewUserServer()
	userpb.RegisterUserServer(server, userServ)

	chatServ := chat.NewChatServer()
	chatpb.RegisterChatServer(server, chatServ)

	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		panic(err)
	}

	fmt.Println("server listen///::::")
	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
