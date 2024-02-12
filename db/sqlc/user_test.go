package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hash_word, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Name:           util.RandomOwner(),
		HashedPassword: hash_word,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.PasswordLastChange)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	// Random test
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Name)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordLastChange.Time, user2.PasswordLastChange.Time)

	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	olduser := createRandomUser(t)

	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Name: olduser.Name,
		FullName: pgtype.Text{
			String: "ele",
			Valid:  true,
		},
	})
	updatedUser, err := testQueries.GetUser(context.Background(), olduser.Name)

	require.NoError(t, err)
	require.NotEqual(t, olduser.FullName, updatedUser.FullName)
	require.Equal(t, olduser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, olduser.Name, updatedUser.Name)
	require.Equal(t, olduser.Email, updatedUser.Email)
}
