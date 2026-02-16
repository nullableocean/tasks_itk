package main

import (
	"context"
	"fmt"
	"log"
	userpb "main/api/user"
	"main/client/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connect, err := grpc.NewClient("127.0.0.1:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer connect.Close()

	usClient := userpb.NewUserClient(connect)
	userClient := user.NewUserClient(usClient)
	testUserClient(userClient)
}

func testUserClient(userClient *user.UserClient) {
	ctx := context.Background()

	crd := user.CreateUserData{
		Name:            "new_user",
		Email:           "email@email.com",
		Password:        "passw1",
		PasswordConfirm: "passw1",
		Role:            user.USER,
	}

	id, err := userClient.Create(ctx, crd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new user id:", id)

	user, err := userClient.Get(ctx, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new user info: ", user)
}
