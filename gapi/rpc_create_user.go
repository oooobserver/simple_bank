package gapi

import (
	"context"
	"errors"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash the password: %s", err)
	}

	arg := db.CreateUserParams{
		Name:           req.GetUsername(),
		HashedPassword: hashword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// The user is already exits
		if errors.Is(err, pgx.ErrTooManyRows) {
			return nil, status.Errorf(codes.Internal, "user already exits: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create the user: %s", err)
	}

	resp := &pb.CreateUserResponse{
		User: converUser(user),
	}

	return resp, nil
}

func converUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Name,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordLastChange.Time),
		CreatedAt:         timestamppb.New(user.CreatedAt.Time),
	}
}
