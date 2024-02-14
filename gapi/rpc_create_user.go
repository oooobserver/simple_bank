package gapi

import (
	"context"
	"errors"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"
	"simplebank/worker"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if vio := validateCreateUserRequest(req); vio != nil {
		err := InvalidArgumentError(vio)
		return nil, err
	}
	hashword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash the password: %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Name:           req.GetUsername(),
			HashedPassword: hashword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendEmail{
				Username: user.Name,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendEmail(ctx, taskPayload, opts...)
		},
	}

	txres, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		// The user is already exits
		if errors.Is(err, pgx.ErrTooManyRows) {
			return nil, status.Errorf(codes.Internal, "user already exits: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create the user: %s", err)
	}

	resp := &pb.CreateUserResponse{
		User: converUser(txres.User),
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

func validateCreateUserRequest(req *pb.CreateUserRequest) (vio []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		vio = append(vio, fieldValidation("username", err))
	}
	if err := val.ValidFullName(req.GetFullName()); err != nil {
		vio = append(vio, fieldValidation("fullname", err))
	}
	if err := val.ValidatePasswaord(req.GetPassword()); err != nil {
		vio = append(vio, fieldValidation("password", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		vio = append(vio, fieldValidation("email", err))
	}

	return vio
}
