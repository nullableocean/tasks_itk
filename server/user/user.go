package user

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"main/api/userpb"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type role int

const (
	unspecified role = iota
	user
	admin
)

type User struct {
	id        int64
	Name      string
	Email     string
	Role      role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserServer struct {
	userpb.UnimplementedUserServer

	store  map[int64]*User
	nextId atomic.Int64

	mu sync.RWMutex
}

func NewUserServer() *UserServer {
	return &UserServer{
		store: make(map[int64]*User),
	}
}

func (us *UserServer) Create(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserId, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	reqID := md.Get("req-id")[0]
	fmt.Println("@user.create req id:", reqID)

	us.mu.Lock()
	defer us.mu.Unlock()

	id := us.nextId.Add(1)
	timestamp := time.Now()

	us.store[id] = &User{
		id:        id,
		Name:      req.Name,
		Email:     req.Email,
		Role:      role(req.Role),
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	return &userpb.UserId{Id: id}, nil
}

func (us *UserServer) Update(ctx context.Context, req *userpb.UpdateUserRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	reqID := md.Get("req-id")[0]
	fmt.Println("@user.update req id:", reqID)

	us.mu.Lock()
	defer us.mu.Unlock()

	user, ex := us.store[req.Id.Id]
	if !ex {
		return nil, errors.New("not found")
	}

	if req.Email != nil {
		user.Email = req.Email.Value
	}

	if req.Name != nil {
		user.Name = req.Name.Value
	}

	return &emptypb.Empty{}, nil
}

func (us *UserServer) Get(ctx context.Context, id *userpb.UserId) (*userpb.UserInfo, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	reqID := md.Get("req-id")[0]
	fmt.Println("@user.get req id:", reqID)

	us.mu.RLock()
	defer us.mu.RUnlock()

	user, ex := us.store[id.Id]
	if !ex {
		return nil, errors.New("not found")
	}

	return &userpb.UserInfo{
		Id:        &userpb.UserId{Id: user.id},
		Name:      user.Name,
		Email:     user.Email,
		Role:      userpb.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (us *UserServer) Delete(ctx context.Context, id *userpb.UserId) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	reqID := md.Get("req-id")[0]
	fmt.Println("@user.delete req id:", reqID)

	us.mu.Lock()
	defer us.mu.Unlock()

	delete(us.store, id.Id)

	return &emptypb.Empty{}, nil
}
