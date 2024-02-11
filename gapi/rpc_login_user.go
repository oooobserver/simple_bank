package gapi

import (
	"context"
	"errors"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if vio := validateLoginUserRequest(req); vio != nil {
		err := InvalidArgumentError(vio)
		return nil, err
	}
	user, err := server.store.GetUser(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to search the user: %s", err)
	}

	err = util.VerifyPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "wrong password: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Name,
		server.config.ACCESS_DURATION,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Name,
		server.config.REFRESH_TOKEN_DURATION,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the refresh token: %s", err)
	}

	// Get context information
	meta := server.extractMetaData(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           util.ConvertGU(refreshPayload.ID),
		Username:     user.Name,
		RefreshToken: refreshToken,
		UserAgent:    meta.UserAgent,
		ClientIp:     meta.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    util.ConverTime(refreshPayload.ExpiredAt),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the session: %s", err)
	}

	resp := &pb.LoginUserResponse{
		SessionId:             uuid.UUID(session.ID.Bytes[:]).String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		User:                  converUser(user),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return resp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetName()); err != nil {
		violations = append(violations, fieldValidation("username", err))
	}

	if err := val.ValidatePasswaord(req.GetPassword()); err != nil {
		violations = append(violations, fieldValidation("password", err))
	}

	return violations
}
