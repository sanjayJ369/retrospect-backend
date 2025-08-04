package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	mockmail "github.com/sanjayj369/retrospect-backend/mail/mock"
	"github.com/sanjayj369/retrospect-backend/ratelimiter"

	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, mailSender *mockmail.MockEmailSender) *Server {
	config := util.Config{
		SymmetricKey:        util.GetRandomString(32),
		AccessTokenDuration: time.Minute * 15,
		TemplatesDir:        "../templates",
	}
	stubratelimiter := ratelimiter.NewStubRateLimiter()
	server, err := NewServer(config, store, mailSender, stubratelimiter)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
