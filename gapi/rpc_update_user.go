package gapi

import (
	"context"
	"errors"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %s", err)
	}

	if authPayload.Name != req.GetName() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's  info: %s", err)
	}

	if vio := validateUpdateUserRequest(req); vio != nil {
		err := InvalidArgumentError(vio)
		return nil, err
	}

	arg := db.UpdateUserParams{
		Name: req.GetName(),
	}

	if req.FullName != nil {
		arg.FullName = util.ConvertString(req.GetFullName())
	}

	if req.Email != nil {
		arg.Email = util.ConvertString(req.GetEmail())
	}

	// Depends on the input
	if req.Password != nil {
		hashword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash the password: %s", err)
		}

		arg.HashedPassword = util.ConvertString(hashword)
		arg.PasswordLastChange = util.ConverTime(time.Now())
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Internal, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update the user: %s", err)
	}

	resp := &pb.UpdateUserResponse{
		User: converUser(user),
	}

	return resp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (vio []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetName()); err != nil {
		vio = append(vio, fieldValidation("username", err))
	}

	if req.FullName != nil {
		if err := val.ValidFullName(req.GetFullName()); err != nil {
			vio = append(vio, fieldValidation("fullname", err))
		}
	}

	if req.Password != nil {
		if err := val.ValidatePasswaord(req.GetPassword()); err != nil {
			vio = append(vio, fieldValidation("password", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			vio = append(vio, fieldValidation("email", err))
		}
	}

	return vio
}
