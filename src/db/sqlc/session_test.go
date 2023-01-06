package db

import (
	"context"
	"testing"
	"time"

	"github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Session {
	config, err := util.LoadConfig("../../..")
	require.NoError(t, err)

	sessionID, err := uuid.NewRandom()
	require.NoError(t, err)

	user := CreateRandomUser(t)
	require.NotEmpty(t, user)

	arg := CreateSessionParams{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: util.RandomString(128),
		UserAgent:    util.RandomString(10),
		ClientIp:     util.RandomString(10),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(config.RefreshTokenDuration),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.UserID, session.UserID)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)

	require.NotZero(t, session.CreatedAt)
	require.NotZero(t, session.ExpiresAt)

	return session
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session1 := createRandomSession(t)
	session2, err := testQueries.GetSession(context.Background(), session1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.UserID, session2.UserID)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.Equal(t, session1.ExpiresAt, session2.ExpiresAt)

	require.WithinDuration(t, session1.CreatedAt, session2.CreatedAt, time.Second)
}
