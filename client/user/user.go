package user

import (
	"context"
	"main/api/userpb"

	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Role int

const (
	unspecified Role = iota
	USER
	ADMIN
)

type User struct {
	id        int64
	Name      string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserClient struct {
	client userpb.UserClient
}

func NewUserClient(grpcClient userpb.UserClient) *UserClient {
	return &UserClient{
		client: grpcClient,
	}
}

func (c *UserClient) Get(ctx context.Context, id int64) (*User, error) {
	md := metadata.New(map[string]string{
		"req-id": "qwerty-1",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	data, err := c.client.Get(ctx, &userpb.UserId{Id: id})
	if err != nil {
		return nil, err
	}

	user := &User{
		id:        id,
		Name:      data.Name,
		Email:     data.Email,
		Role:      Role(data.Role),
		CreatedAt: data.CreatedAt.AsTime(),
		UpdatedAt: data.UpdatedAt.AsTime(),
	}

	return user, nil
}

type CreateUserData struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            Role
}

func (c *UserClient) Create(ctx context.Context, data CreateUserData) (int64, error) {
	md := metadata.New(map[string]string{
		"req-id": "qwerty-2",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	r := userpb.CreateUserRequest{
		Name:            data.Name,
		Email:           data.Email,
		Password:        data.Password,
		PasswordConfirm: data.PasswordConfirm,
		Role:            userpb.Role(data.Role),
	}

	id, err := c.client.Create(ctx, &r)
	return id.Id, err
}

type UpdateUserData struct {
	Name  string
	Email string
}

func (c *UserClient) Update(ctx context.Context, id int64, data UpdateUserData) error {
	md := metadata.New(map[string]string{
		"req-id": "qwerty-3",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	r := userpb.UpdateUserRequest{
		Id:    &userpb.UserId{Id: id},
		Name:  &wrapperspb.StringValue{Value: data.Name},
		Email: &wrapperspb.StringValue{Value: data.Email},
	}

	_, err := c.client.Update(ctx, &r)
	return err
}

func (c *UserClient) Delete(ctx context.Context, id int64) error {
	md := metadata.New(map[string]string{
		"req-id": "qwerty-4",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := c.client.Delete(ctx, &userpb.UserId{Id: id})
	return err
}
