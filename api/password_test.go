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

func TestResetPasswordAPI(t *testing.T) {
	user := randomUser()
	newPassword := util.GetRandomString(10)

	// Create a valid token for testing
	tokenMaker, err := token.NewPasetoMaker(util.GetRandomString(32))
	require.NoError(t, err)

	validToken, _, err := tokenMaker.CreateToken(user.ID.Bytes, time.Minute*15, token.PurposeResetPassword)
	require.NoError(t, err)

	invalidPurposeToken, _, err := tokenMaker.CreateToken(user.ID.Bytes, time.Minute*15, token.PurposeLogin)
	require.NoError(t, err)

	expiredToken, _, err := tokenMaker.CreateToken(user.ID.Bytes, -time.Minute, token.PurposeResetPassword)
	require.NoError(t, err)

	resetReq := resetPasswordRequest{
		Token:       validToken,
		NewPassword: newPassword,
	}

	invalidPurposeReq := resetPasswordRequest{
		Token:       invalidPurposeToken,
		NewPassword: newPassword,
	}

	expiredTokenReq := resetPasswordRequest{
		Token:       expiredToken,
		NewPassword: newPassword,
	}

	invalidTokenReq := resetPasswordRequest{
		Token:       "invalid-token",
		NewPassword: newPassword,
	}

	shortPasswordReq := resetPasswordRequest{
		Token:       validToken,
		NewPassword: "123", // Too short
	}

	marshalledReq, err := json.Marshal(resetReq)
	require.NoError(t, err)
	validResetDetails := bytes.NewReader(marshalledReq)

	marshalledInvalidPurposeReq, err := json.Marshal(invalidPurposeReq)
	require.NoError(t, err)
	invalidPurposeDetails := bytes.NewReader(marshalledInvalidPurposeReq)

	marshalledExpiredReq, err := json.Marshal(expiredTokenReq)
	require.NoError(t, err)
	expiredTokenDetails := bytes.NewReader(marshalledExpiredReq)

	marshalledInvalidTokenReq, err := json.Marshal(invalidTokenReq)
	require.NoError(t, err)
	invalidTokenDetails := bytes.NewReader(marshalledInvalidTokenReq)

	marshalledShortPasswordReq, err := json.Marshal(shortPasswordReq)
	require.NoError(t, err)
	shortPasswordDetails := bytes.NewReader(marshalledShortPasswordReq)

	testCases := []struct {
		name          string
		resetDetails  *bytes.Reader
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			resetDetails: validResetDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), EqUpdateUserHashedPasswordParams(db.UpdateUserHashedPasswordParams{
						ID: pgtype.UUID{Bytes: user.ID.Bytes, Valid: true},
					}, newPassword)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, "Password reset successfully", response["message"])
			},
		},
		{
			name:         "InvalidToken",
			resetDetails: invalidTokenDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:         "ExpiredToken",
			resetDetails: expiredTokenDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:         "InvalidTokenPurpose",
			resetDetails: invalidPurposeDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "ShortPassword",
			resetDetails: shortPasswordDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "DatabaseError",
			resetDetails: validResetDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), EqUpdateUserHashedPasswordParams(db.UpdateUserHashedPasswordParams{
						ID: pgtype.UUID{Bytes: user.ID.Bytes, Valid: true},
					}, newPassword)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "BadRequestMalformedJSON",
			resetDetails: bytes.NewReader([]byte("invalid json")),
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "MissingToken",
			resetDetails: bytes.NewReader([]byte(`{"new_password":"password123"}`)),
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "MissingPassword",
			resetDetails: bytes.NewReader([]byte(`{"token":"some-token"}`)),
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateUserHashedPassword(gomock.Any(), gomock.Any()).
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
			tc.resetDetails.Seek(0, 0)

			server := newTestServer(t, store, nil)
			// Override the tokenMaker with our test tokenMaker to ensure consistent token verification
			server.tokenMaker = tokenMaker
			recorder := httptest.NewRecorder()

			url := "/users/reset-password"

			req, err := http.NewRequest(http.MethodPost, url, tc.resetDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestForgotPasswordAPI(t *testing.T) {
	user := randomUser()

	forgotReq := forgotPasswordRequest{
		Email: user.Email,
	}

	invalidEmailReq := forgotPasswordRequest{
		Email: "invalid-email",
	}

	marshalledReq, err := json.Marshal(forgotReq)
	require.NoError(t, err)
	validEmailDetails := bytes.NewReader(marshalledReq)

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
					SendMail(
						gomock.Eq("Password Reset Request"),
						gomock.Any(),
						gomock.Eq([]string{user.Email}),
						gomock.Eq([]string(nil)),
						gomock.Eq([]string(nil)),
						gomock.Eq([]string(nil)),
					).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, "Password reset email sent successfully", response["message"])
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
		{
			name:         "MissingEmail",
			emailDetails: bytes.NewReader([]byte(`{}`)),
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

			url := "/users/forgot-password"

			req, err := http.NewRequest(http.MethodPost, url, tc.emailDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestSendPasswordResetMail(t *testing.T) {
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
						gomock.Eq("Password Reset Request"),
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
				require.Contains(t, err.Error(), "failed to send email")
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

			err := SendPasswordResetMail(
				emailSender,
				user.ID.Bytes,
				user.Email,
				tokenMaker,
				time.Minute*15,
				"https://example.com/users/forgot-password",
				"../templates/password_reset.html",
			)

			tc.checkResponse(t, err)
		})
	}
}

// Helper function to create a matcher for UpdateUserHashedPasswordParams
type eqUpdateUserHashedPasswordParamsMatcher struct {
	arg      db.UpdateUserHashedPasswordParams
	password string
}

func (e eqUpdateUserHashedPasswordParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.UpdateUserHashedPasswordParams)
	if !ok {
		return false
	}

	if err := util.CheckPassword(e.password, arg.HashedPassword); err != nil {
		return false
	}

	return arg.ID.Bytes == e.arg.ID.Bytes && arg.ID.Valid == e.arg.ID.Valid
}

func (e eqUpdateUserHashedPasswordParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v (password should match)", e.arg)
}

func EqUpdateUserHashedPasswordParams(arg db.UpdateUserHashedPasswordParams, password string) gomock.Matcher {
	return eqUpdateUserHashedPasswordParamsMatcher{
		arg:      arg,
		password: password,
	}
}
