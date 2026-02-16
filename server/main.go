package main

import (
	userpb "main/api/user"
	"main/server/user"
	"net"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()

	userServ := user.NewUserServer()
	userpb.RegisterUserServer(server, userServ)

	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		panic(err)
	}

	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
