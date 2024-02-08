package api

import (
	"errors"
	"fmt"
	"net/http"
	"simplebank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type renewAccessTokenrRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenrResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessTokenr(ctx *gin.Context) {
	var req renewAccessTokenrRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, util.ConvertGU(refreshPayload.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check if blocked
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Check user name
	if session.Username != refreshPayload.Name {
		err := fmt.Errorf("incorrect user name")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Check token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("incorrect token")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Check expired
	if time.Now().After(session.ExpiresAt.Time) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Name,
		server.config.ACCESS_DURATION,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := renewAccessTokenrResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
