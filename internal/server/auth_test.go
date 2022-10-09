package server

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

	mh "github.com/go-developer-ya-practicum/gophkeeper/internal/hasher/mock"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/storage"
	ms "github.com/go-developer-ya-practicum/gophkeeper/internal/storage/mock"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/token"
	mt "github.com/go-developer-ya-practicum/gophkeeper/internal/token/mock"
)

const address = "127.0.0.1:5050"

func newClient() (pb.GophKeeperClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewGophKeeperClient(conn), nil
}

func newServer(t *testing.T) (*Server, func()) {
	ctrl := gomock.NewController(t)

	server := &Server{
		Storage:        ms.NewMockStorage(ctrl),
		TokenGenerator: mt.NewMockGenerator(ctrl),
		Hasher:         mh.NewMockHasher(ctrl),
		Address:        address,
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run(ctx)
	}()

	return server, func() {
		cancel()
		wg.Wait()
		ctrl.Finish()
	}
}

func TestServer_SignUp(t *testing.T) {
	server, cancel := newServer(t)
	defer cancel()

	testHasher, ok := server.Hasher.(*mh.MockHasher)
	require.True(t, ok)

	testStorage, ok := server.Storage.(*ms.MockStorage)
	require.True(t, ok)

	testTokenGenerator, ok := server.TokenGenerator.(*mt.MockGenerator)
	require.True(t, ok)

	client, err := newClient()
	require.NoError(t, err)

	t.Run("Empty Email", func(t *testing.T) {
		request := &pb.SignUpRequest{Email: "", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("Empty Password", func(t *testing.T) {
		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: ""}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("Hasher Error", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("", errors.New("hasher error"))

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Email Exists", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), gomock.Any()).
			Return(0, storage.ErrEmailIsAlreadyInUse)

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.AlreadyExists)
	})

	t.Run("Storage PutUser Error", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), gomock.Any()).
			Return(0, errors.New("storage error"))

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Token Create Error", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), gomock.Any()).
			Return(0, nil)

		testTokenGenerator.
			EXPECT().
			Create(gomock.Any()).
			Return("", errors.New("failed to create token"))

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignUp(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Success Register", func(t *testing.T) {
		userID := 0
		accessToken := "aaa.bbb.ccc"

		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			PutUser(gomock.Any(), gomock.Any()).
			Return(userID, nil)

		testTokenGenerator.
			EXPECT().
			Create(gomock.Any()).
			Return(accessToken, nil)

		request := &pb.SignUpRequest{Email: "test@mail.ru", Password: "password"}
		resp, err := client.SignUp(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, accessToken, resp.AccessToken)
	})
}

func TestServer_SignIn(t *testing.T) {
	server, cancel := newServer(t)
	defer cancel()

	testHasher, ok := server.Hasher.(*mh.MockHasher)
	require.True(t, ok)

	testStorage, ok := server.Storage.(*ms.MockStorage)
	require.True(t, ok)

	testTokenGenerator, ok := server.TokenGenerator.(*mt.MockGenerator)
	require.True(t, ok)

	client, err := newClient()
	require.NoError(t, err)

	t.Run("Empty Email", func(t *testing.T) {
		request := &pb.SignInRequest{Email: "", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("Empty Password", func(t *testing.T) {
		request := &pb.SignInRequest{Email: "test@mail.ru", Password: ""}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})

	t.Run("Hasher Error", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("", errors.New("hasher error"))

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), gomock.Any()).
			Return(0, storage.ErrInvalidCredentials)

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("Storage GetUser Error", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), gomock.Any()).
			Return(0, errors.New("storage error"))

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Token Create Error", func(t *testing.T) {
		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), gomock.Any()).
			Return(0, nil)

		testTokenGenerator.
			EXPECT().
			Create(gomock.Any()).
			Return("", errors.New("failed to create token"))

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		_, err = client.SignIn(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Success Login", func(t *testing.T) {
		userID := 0
		accessToken := "aaa.bbb.ccc"

		testHasher.
			EXPECT().
			Hash(gomock.Any()).
			Return("PasswordHash", nil)

		testStorage.
			EXPECT().
			GetUser(gomock.Any(), gomock.Any()).
			Return(userID, nil)

		testTokenGenerator.
			EXPECT().
			Create(gomock.Any()).
			Return(accessToken, nil)

		request := &pb.SignInRequest{Email: "test@mail.ru", Password: "password"}
		resp, err := client.SignIn(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, accessToken, resp.AccessToken)
	})
}

func TestServer_VerifyToken(t *testing.T) {
	server, cancel := newServer(t)
	defer cancel()

	testTokenGenerator, ok := server.TokenGenerator.(*mt.MockGenerator)
	require.True(t, ok)

	client, err := newClient()
	require.NoError(t, err)

	t.Run("Empty Token", func(t *testing.T) {
		request := &pb.VerifyTokenRequest{AccessToken: ""}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("Expired Token", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"

		testTokenGenerator.
			EXPECT().
			Validate(accessToken).
			Return(nil, token.ErrExpiredToken)

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"

		testTokenGenerator.
			EXPECT().
			Validate(accessToken).
			Return(nil, token.ErrInvalidToken)

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Unauthenticated)
	})

	t.Run("Validate Token Error", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"

		testTokenGenerator.
			EXPECT().
			Validate(accessToken).
			Return(nil, errors.New("validate token error"))

		request := &pb.VerifyTokenRequest{AccessToken: accessToken}
		_, err = client.VerifyToken(context.Background(), request)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("Valid Token", func(t *testing.T) {
		accessToken := "aaa.bbb.ccc"
		userID := 0

		testTokenGenerator.
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
