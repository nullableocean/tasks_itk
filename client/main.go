package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"main/api/chatpb"
	"main/client/chat"
	"main/client/user"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connect, err := grpc.NewClient("127.0.0.1:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer connect.Close()

	// тестирование grpc для юзера
	// usClient := userpb.NewUserClient(connect)
	// userClient := user.NewUserClient(usClient)
	// testUserClient(userClient)

	chatpbClient := chatpb.NewChatClient(connect)
	chatClient := chat.NewChatClient(chatpbClient)

	startChatClient(chatClient)
}

func startChatClient(client *chat.ChatClient) {
	var err error
	var chatId int64

	args := os.Args[1:]
	fmt.Println(args)

	wg := sync.WaitGroup{}

	ctx := context.Background()

	switch args[0] {
	case "new":
		chatId, err = client.Create(ctx, []string{"anon"})
		if err != nil {
			log.Fatalln("create chat error:", err)
		}
		log.Println("CHAT_ID: ", chatId)
	case "del":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln(`incorrect chat id, wait number. "del [chatid]"`)
		}

		err = client.Delete(ctx, int64(i))
		if err != nil {
			log.Fatalln("delete chat error:", err)
		}

		log.Println("DELETED")
		return
	default:
		i, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln(`incorrect arguments. expect "new", "del [chatid]" or [chatid]`)
		}
		chatId = int64(i)
	}

	inMes := make(chan *chat.Message)
	opts := &chat.SendOpt{
		ChatId: chatId,
		In:     inMes,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = client.SendMessages(ctx, opts)
		if err != nil {
			log.Fatalln("send messages error:", err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		defer close(inMes)

		for scanner.Scan() {
			text := scanner.Text()
			if text == "/exit" {
				break
			}

			m := &chat.Message{
				Text: text,
				From: "anon",
			}

			inMes <- m
		}
	}()

	err = client.GetMessages(ctx, os.Stdout, chatId)
	if err != nil {
		log.Fatalln("get messages error:", err)
	}

	wg.Wait()
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
		log.Fatal("create user error", err)
	}

	fmt.Println("new user id:", id)

	user, err := userClient.Get(ctx, id)
	if err != nil {
		log.Fatal("get user error", err)
	}

	fmt.Println("new user info: ", user)
	fmt.Println()
	fmt.Println()
}
