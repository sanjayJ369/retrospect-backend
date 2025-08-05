package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	mockDB "github.com/sanjayj369/retrospect-backend/db/mock"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	mockmail "github.com/sanjayj369/retrospect-backend/mail/mock"
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestResendVerificationEmailAPI(t *testing.T) {
	user := randomUser()
	user.IsVerified = false // Ensure user is not verified for testing

	verifiedUser := randomUser()
	verifiedUser.IsVerified = true

	resendReq := ResendVerificationEmailRequest{
		Email: user.Email,
	}

	verifiedUserReq := ResendVerificationEmailRequest{
		Email: verifiedUser.Email,
	}

	marshalledReq, err := json.Marshal(resendReq)
	require.NoError(t, err)
	validEmailDetails := bytes.NewReader(marshalledReq)

	marshalledVerifiedReq, err := json.Marshal(verifiedUserReq)
	require.NoError(t, err)
	verifiedEmailDetails := bytes.NewReader(marshalledVerifiedReq)

	invalidEmailReq := ResendVerificationEmailRequest{
		Email: "invalid-email",
	}
	marshalledInvalidReq, err := json.Marshal(invalidEmailReq)
	require.NoError(t, err)
	invalidEmailDetails := bytes.NewReader(marshalledInvalidReq)

	testCases := []struct {
		name          string
		emailDetails  *bytes.Reader
		buildStub     func(store *mockDB.MockStore, mail *mockmail.MockEmailSender)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			emailDetails: validEmailDetails,
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Eq([]string{user.Email}), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, "Verification email sent successfully", response["message"])
			},
		},
		{
			name:         "UserAlreadyVerified",
			emailDetails: verifiedEmailDetails,
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(verifiedUser.Email)).
					Times(1).
					Return(verifiedUser, nil)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, "Email already verified", response["message"])
			},
		},
		{
			name:         "UserNotFound",
			emailDetails: validEmailDetails,
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:         "InvalidEmailFormat",
			emailDetails: invalidEmailDetails,
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "DatabaseError",
			emailDetails: validEmailDetails,
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code) // API returns 404 for any GetUserByEmail error
			},
		},
		{
			name:         "EmailSendError",
			emailDetails: validEmailDetails,
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Eq([]string{user.Email}), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("failed to send email"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "BadRequestMalformedJSON",
			emailDetails: bytes.NewReader([]byte("invalid json")),
			buildStub: func(store *mockDB.MockStore, mail *mockmail.MockEmailSender) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)

				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			emailSender := mockmail.NewMockEmailSender(ctrl)

			tc.buildStub(store, emailSender)
			tc.emailDetails.Seek(0, 0)

			server := newTestServer(t, store, emailSender)
			recorder := httptest.NewRecorder()

			url := "/users/resend-verification"

			req, err := http.NewRequest(http.MethodPost, url, tc.emailDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestVerifyEmailAPI(t *testing.T) {
	user := randomUser()
	user.IsVerified = false

	updatedUser := user
	updatedUser.IsVerified = true

	// Create a valid token for testing
	tokenMaker, err := token.NewPasetoMaker(util.GetRandomString(32))
	require.NoError(t, err)

	validToken, _, err := tokenMaker.CreateToken(user.ID.Bytes, time.Minute*15, token.PurposeVerifyEmail)
	require.NoError(t, err)

	invalidPurposeToken, _, err := tokenMaker.CreateToken(user.ID.Bytes, time.Minute*15, token.PurposeLogin)
	require.NoError(t, err)

	expiredToken, _, err := tokenMaker.CreateToken(user.ID.Bytes, -time.Minute, token.PurposeVerifyEmail)
	require.NoError(t, err)

	verifyReq := VerifyEmailRequest{
		Token: validToken,
	}

	invalidPurposeReq := VerifyEmailRequest{
		Token: invalidPurposeToken,
	}

	expiredTokenReq := VerifyEmailRequest{
		Token: expiredToken,
	}

	invalidTokenReq := VerifyEmailRequest{
		Token: "invalid-token",
	}

	marshalledReq, err := json.Marshal(verifyReq)
	require.NoError(t, err)
	validTokenDetails := bytes.NewReader(marshalledReq)

	marshalledInvalidPurposeReq, err := json.Marshal(invalidPurposeReq)
	require.NoError(t, err)
	invalidPurposeDetails := bytes.NewReader(marshalledInvalidPurposeReq)

	marshalledExpiredReq, err := json.Marshal(expiredTokenReq)
	require.NoError(t, err)
	expiredTokenDetails := bytes.NewReader(marshalledExpiredReq)

	marshalledInvalidTokenReq, err := json.Marshal(invalidTokenReq)
	require.NoError(t, err)
	invalidTokenDetails := bytes.NewReader(marshalledInvalidTokenReq)

	testCases := []struct {
		name          string
		tokenDetails  *bytes.Reader
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			tokenDetails: validTokenDetails,
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateUserIsVerifiedParams{
					ID:         pgtype.UUID{Bytes: user.ID.Bytes, Valid: true},
					IsVerified: true,
				}
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, "Email verified successfully", response["message"])
			},
		},
		{
			name:         "InvalidToken",
			tokenDetails: invalidTokenDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:         "ExpiredToken",
			tokenDetails: expiredTokenDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:         "InvalidTokenPurpose",
			tokenDetails: invalidPurposeDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "DatabaseError",
			tokenDetails: validTokenDetails,
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateUserIsVerifiedParams{
					ID:         pgtype.UUID{Bytes: user.ID.Bytes, Valid: true},
					IsVerified: true,
				}
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "BadRequestMalformedJSON",
			tokenDetails: bytes.NewReader([]byte("invalid json")),
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "MissingToken",
			tokenDetails: bytes.NewReader([]byte(`{"token":""}`)),
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserIsVerified(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)

			tc.buildStub(store)
			tc.tokenDetails.Seek(0, 0)

			server := newTestServer(t, store, nil)
			// Override the tokenMaker with our test tokenMaker to ensure consistent token verification
			server.tokenMaker = tokenMaker
			recorder := httptest.NewRecorder()

			url := "/users/verify-email"

			req, err := http.NewRequest(http.MethodPost, url, tc.tokenDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestSendVerificationMail(t *testing.T) {
	user := randomUser()
	tokenMaker, err := token.NewPasetoMaker(util.GetRandomString(32))
	require.NoError(t, err)

	testCases := []struct {
		name          string
		buildStub     func(mail *mockmail.MockEmailSender)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			buildStub: func(mail *mockmail.MockEmailSender) {
				mail.EXPECT().
					SendMail(
						gomock.Eq("Verify your email address"),
						gomock.Any(),
						gomock.Eq([]string{user.Email}),
						gomock.Eq([]string(nil)),
						gomock.Eq([]string(nil)),
						gomock.Eq([]string(nil)),
					).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "EmailSendError",
			buildStub: func(mail *mockmail.MockEmailSender) {
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("failed to send email"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "sending email failed")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			emailSender := mockmail.NewMockEmailSender(ctrl)
			tc.buildStub(emailSender)

			err := SendVerificationMail(
				emailSender,
				user.ID.Bytes,
				user.Email,
				tokenMaker,
				time.Minute*15,
				"https://example.com/verify-email",
				"../templates/email_verification.html",
			)

			tc.checkResponse(t, err)
		})
	}
}
