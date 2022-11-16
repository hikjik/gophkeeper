package services

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
	ms "github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage/mock"
	mh "github.com/go-developer-ya-practicum/gophkeeper/pkg/hasher/mock"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/token"
	mt "github.com/go-developer-ya-practicum/gophkeeper/pkg/token/mock"
)

const address = "127.0.0.1:5050"

func newAuthClient() (pb.AuthServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewAuthServiceClient(conn), nil
}

func newAuthService(t *testing.T) (*AuthService, func()) {
	const address = "127.0.0.1:5050"

	ctrl := gomock.NewController(t)

	authService := &AuthService{
		UserStorage:  ms.NewMockUserStorage(ctrl),
		TokenManager: mt.NewMockManager(ctrl),
		Hasher:       mh.NewMockHasher(ctrl),
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := NewServer(address, WithServices(authService))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run(ctx)
	}()

	return authService, func() {
		cancel()
		wg.Wait()
		ctrl.Finish()
	}
}

func TestServer_SignUp(t *testing.T) {
	authService, cancel := newAuthService(t)
	defer cancel()

	testHasher, ok := authService.Hasher.(*mh.MockHasher)
	require.True(t, ok)

	testStorage, ok := authService.UserStorage.(*ms.MockUserStorage)
	require.True(t, ok)

	testTokenManager, ok := authService.TokenManager.(*mt.MockManager)
	require.True(t, ok)

	client, err := newAuthClient()
	require.NoError(t, err)

	t.Run("EmptyEmail", func(t *testing.T) {
		request := &pb.SignUpRequest{Email: "", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: ""}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("HasherError", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("", errors.New("hasher error"))

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("EmailExists", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), gomock.Any()).
			Return(nil, storage.ErrUserConflict)

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.AlreadyExists)
	})

	t.Run("StoragePutUserError", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("storage error"))

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("TokenCreateError", func(t *testing.T) {
		user := &models.User{
			ID:           0,
			Email:        "test@mail.ru",
			PasswordHash: "hash",
		}
		password := "password"

		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return(user.PasswordHash, nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), user).
			Return(user, nil)

		testTokenManager.
			EXPECT().
			Create(gomock.Any()).
			Return("", errors.New("failed to create token"))

		request := &pb.SignUpRequest{Email: user.Email, Password: password}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("SuccessfulRegister", func(t *testing.T) {
		user := &models.User{
			ID:           0,
			Email:        "test@mail.ru",
			PasswordHash: "hash",
		}
		password := "password"
		accessToken := "aaa.bbb.ccc"

		testHasher.
			EXPECT().
			Hash(password).
			Return(user.PasswordHash, nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), user).
			Return(user, nil)

		testTokenManager.
			EXPECT().
			Create(gomock.Any()).
			Return(accessToken, nil)

		request := &pb.SignUpRequest{Email: user.Email, Password: password}
		resp, err := client.SignUp(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, accessToken, resp.AccessToken)
	})
}

func TestServer_SignIn(t *testing.T) {
	authService, cancel := newAuthService(t)
	defer cancel()

	testHasher, ok := authService.Hasher.(*mh.MockHasher)
	require.True(t, ok)

	testStorage, ok := authService.UserStorage.(*ms.MockUserStorage)
	require.True(t, ok)

	testTokenManager, ok := authService.TokenManager.(*mt.MockManager)
	require.True(t, ok)

	client, err := newAuthClient()
	require.NoError(t, err)

	t.Run("EmptyEmail", func(t *testing.T) {
		request := &pb.SignInRequest{Email: "", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		request := &pb.SignInRequest{Email: "test@mail.ru", Password: ""}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("HasherError", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("", errors.New("hasher error"))

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), gomock.Any()).
			Return(nil, storage.ErrUserNotFound)

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("StorageGetUserError", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("storage error"))

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("TokenCreateError", func(t *testing.T) {
		user := &models.User{
			ID:           0,
			Email:        "test@mail.ru",
			PasswordHash: "hash",
		}
		password := "password"

		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return(user.PasswordHash, nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), user).
			Return(user, nil)

		testTokenManager.
			EXPECT().
			Create(gomock.Any()).
			Return("", errors.New("failed to create token"))

		request := &pb.SignInRequest{Email: user.Email, Password: password}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("SuccessLogin", func(t *testing.T) {
		user := &models.User{
			ID:           0,
			Email:        "test@mail.ru",
			PasswordHash: "hash",
		}
		password := "password"
		accessToken := "aaa.bbb.ccc"

		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return(user.PasswordHash, nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), user).
			Return(user, nil)

		testTokenManager.
			EXPECT().
			Create(gomock.Any()).
			Return(accessToken, nil)

		request := &pb.SignInRequest{Email: user.Email, Password: password}
		resp, err := client.SignIn(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, accessToken, resp.AccessToken)
	})
}

func TestServer_VerifyToken(t *testing.T) {
	authService, cancel := newAuthService(t)
	defer cancel()

	testTokenManager, ok := authService.TokenManager.(*mt.MockManager)
	require.True(t, ok)

	client, err := newAuthClient()
	require.NoError(t, err)

	t.Run("EmptyToken", func(t *testing.T) {
		request := &pb.VerifyTokenRequest{AccessToken: ""}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"

		testTokenManager.
			EXPECT().
			Validate(accessToken).
			Return(nil, token.ErrExpiredToken)

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"

		testTokenManager.
			EXPECT().
			Validate(accessToken).
			Return(nil, token.ErrInvalidToken)

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("ValidateTokenError", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"

		testTokenManager.
			EXPECT().
			Validate(accessToken).
			Return(nil, errors.New("validate token error"))

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("ValidToken", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"
		userID := 0

		testTokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		resp, err := client.VerifyToken(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, userID, int(resp.UserId))
	})
}

func checkErrorStatus(t *testing.T, err error, code codes.Code) {
	require.Error(t, err)

	errStatus, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, code, errStatus.Code())
}
