package chat

import (
	"context"
	"errors"
	"io"
	"main/api/chatpb"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatServer struct {
	chatpb.UnimplementedChatServer
	chats  map[int64]*Chat
	nextId atomic.Int64

	mu sync.RWMutex
}

type Chat struct {
	id        int64
	usernames []string
	messages  []*chatpb.ChatMessage
	streams   []grpc.ServerStreamingServer[chatpb.ChatMessage]

	activeOps sync.WaitGroup
	isDeleted atomic.Bool

	// возмоэно нужен отдельный мьютекст для streams
	mu sync.RWMutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		chats: make(map[int64]*Chat),
		mu:    sync.RWMutex{},
	}
}

func (cs *ChatServer) Create(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.CreateChatResponse, error) {
	usernames := make([]string, 0, len(req.Usernames))
	for _, u := range req.Usernames {
		usernames = append(usernames, u.Username)
	}

	id := cs.nextId.Add(1)
	chat := &Chat{
		id:        id,
		usernames: usernames,
		messages:  make([]*chatpb.ChatMessage, 0),
		streams:   make([]grpc.ServerStreamingServer[chatpb.ChatMessage], 0),
		mu:        sync.RWMutex{},
	}

	cs.mu.Lock()
	cs.chats[id] = chat
	cs.mu.Unlock()

	return &chatpb.CreateChatResponse{Id: &chatpb.ChatId{Id: id}}, nil
}

func (cs *ChatServer) Delete(ctx context.Context, id *chatpb.ChatId) (*emptypb.Empty, error) {
	cs.mu.Lock()
	chat, exists := cs.chats[id.Id]
	if !exists {
		cs.mu.Unlock()
		return nil, errors.New("chat not found")
	}
	delete(cs.chats, id.Id)
	cs.mu.Unlock()

	// если у удаленного чата остались активные операции ждём
	// если клиент не отключается, может ждать вечно
	chat.isDeleted.Store(true)
	chat.activeOps.Wait()

	return &emptypb.Empty{}, nil
}

func (cs *ChatServer) SendMessages(stream grpc.ClientStreamingServer[chatpb.ChatMessage, emptypb.Empty]) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}

		cs.mu.RLock()
		chat, exists := cs.chats[msg.Id.Id]
		cs.mu.RUnlock()

		if !exists {
			return errors.New("chat not found")
		}

		chat.activeOps.Add(1)
		if chat.isDeleted.Load() {
			chat.activeOps.Done()
			return errors.New("chat deleted")
		}

		msg.Timestamp = timestamppb.Now()
		chat.mu.Lock()
		chat.messages = append(chat.messages, msg)
		for _, s := range chat.streams {
			s.Send(msg)
		}
		chat.mu.Unlock()

		chat.activeOps.Done()
	}
}

func (cs *ChatServer) GetMessages(id *chatpb.ChatId, stream grpc.ServerStreamingServer[chatpb.ChatMessage]) error {
	cs.mu.RLock()
	chat, exists := cs.chats[id.Id]
	cs.mu.RUnlock()

	if !exists {
		return errors.New("chat not found")
	}

	chat.activeOps.Add(1)
	defer chat.activeOps.Done()

	if chat.isDeleted.Load() {
		return errors.New("chat deleted")
	}

	chat.mu.Lock()

	// отдаём новому потоку все сохраненные сообщения
	for _, msg := range chat.messages {
		if err := stream.Send(msg); err != nil {
			chat.mu.RUnlock()
			return err
		}
	}

	chat.streams = append(chat.streams, stream)
	chat.mu.Unlock()

	// ожидаем отключение клиента
	<-stream.Context().Done()

	chat.mu.Lock()
	for i, s := range chat.streams {
		if s == stream {
			chat.streams = append(chat.streams[:i], chat.streams[i+1:]...)
			break
		}
	}
	chat.mu.Unlock()

	return nil
}
