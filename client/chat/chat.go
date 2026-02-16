package chat

import (
	"context"
	"fmt"
	"io"
	"main/api/chatpb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Text string
	From string
}
type ChatClient struct {
	client chatpb.ChatClient
}

func NewChatClient(client chatpb.ChatClient) *ChatClient {
	return &ChatClient{
		client: client,
	}
}

func (c *ChatClient) Create(ctx context.Context, usernames []string) (int64, error) {
	usrs := make([]*chatpb.Username, len(usernames))
	for _, u := range usernames {
		usrs = append(usrs, &chatpb.Username{Username: u})
	}

	r := &chatpb.CreateChatRequest{
		Usernames: usrs,
	}

	res, err := c.client.Create(ctx, r)
	if err != nil {
		return 0, err
	}

	return res.Id.Id, nil
}

func (c *ChatClient) Delete(ctx context.Context, id int64) error {
	_, err := c.client.Delete(ctx, &chatpb.ChatId{Id: id})
	return err
}

type SendOpt struct {
	ChatId int64
	In     <-chan *Message
}

func (c *ChatClient) SendMessages(ctx context.Context, opt *SendOpt) error {
	stream, err := c.client.SendMessages(ctx)
	if err != nil {
		return err
	}

	fmt.Println("ok")
	for m := range opt.In {
		message := &chatpb.ChatMessage{
			Id:        &chatpb.ChatId{Id: opt.ChatId},
			From:      m.From,
			Text:      m.Text,
			Timestamp: timestamppb.Now(),
		}

		err := stream.Send(message)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}

func (c *ChatClient) GetMessages(ctx context.Context, w io.Writer, id int64) error {
	stream, err := c.client.GetMessages(ctx, &chatpb.ChatId{Id: id})
	if err != nil {
		return err
	}

LOOP:
	for {
		message, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}

			return err
		}

		m := fmt.Sprintf("[%s]: %s\n", message.From, message.Text)
		w.Write([]byte(m))
	}

	return err
}
