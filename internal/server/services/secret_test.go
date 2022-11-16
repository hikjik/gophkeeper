package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	clientInterceptors "github.com/go-developer-ya-practicum/gophkeeper/internal/client/interceptors"
	pb "github.com/go-developer-ya-practicum/gophkeeper/internal/proto"
	serverInterceptors "github.com/go-developer-ya-practicum/gophkeeper/internal/server/interceptors"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/models"
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage"
	ms "github.com/go-developer-ya-practicum/gophkeeper/internal/server/storage/mock"
	"github.com/go-developer-ya-practicum/gophkeeper/pkg/token"
	mt "github.com/go-developer-ya-practicum/gophkeeper/pkg/token/mock"
)

func newSecretClient(accessToken string) (pb.SecretServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientInterceptors.NewAuthInterceptor(accessToken).Unary()),
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewSecretServiceClient(conn), nil
}

func TestSecretService_GetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secretStorage := ms.NewMockSecretStorage(ctrl)
	tokenManager := mt.NewMockManager(ctrl)

	secretService := &SecretService{
		SecretStorage: secretStorage,
	}

	interceptor := serverInterceptors.NewAuthInterceptor(tokenManager)

	server := NewServer(
		address,
		WithServices(secretService),
		WithUnaryInterceptors(interceptor.Unary()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Run(ctx)

	accessToken := "Token"
	userID := 0

	t.Run("EmptySecretName", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.GetSecret(context.Background(), &pb.GetSecretRequest{Name: ""})
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("SecretNotFound", func(t *testing.T) {
		secretName := "SecretName"

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		secretStorage.
			EXPECT().
			GetSecret(gomock.Any(), secretName, userID).
			Return(nil, storage.ErrSecretNotFound)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.GetSecret(context.Background(), &pb.GetSecretRequest{Name: secretName})
		checkErrorStatus(t, err, codes.NotFound)
	})
	t.Run("StorageError", func(t *testing.T) {
		secretName := "SecretName"

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		secretStorage.
			EXPECT().
			GetSecret(gomock.Any(), secretName, userID).
			Return(nil, errors.New("some error"))

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.GetSecret(context.Background(), &pb.GetSecretRequest{Name: secretName})
		checkErrorStatus(t, err, codes.Internal)
	})
	t.Run("SuccessfulGetSecret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.New(),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			GetSecret(gomock.Any(), secret.Name, secret.OwnerID).
			Return(secret, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.GetSecret(
			context.Background(), &pb.GetSecretRequest{Name: secret.Name})
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
		assert.Equal(t, secret.Content, resp.Content)
		assert.Equal(t, secret.Version.String(), resp.Version)
	})
}

func TestSecretService_CreateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secretStorage := ms.NewMockSecretStorage(ctrl)
	tokenManager := mt.NewMockManager(ctrl)

	secretService := &SecretService{
		SecretStorage: secretStorage,
	}

	interceptor := serverInterceptors.NewAuthInterceptor(tokenManager)

	server := NewServer(
		address,
		WithServices(secretService),
		WithUnaryInterceptors(interceptor.Unary()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Run(ctx)

	accessToken := "Token"
	userID := 0

	t.Run("EmptySecretName", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.CreateSecret(
			context.Background(), &pb.CreateSecretRequest{Name: "", Content: []byte("Content")})
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("EmptySecretContent", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.CreateSecret(
			context.Background(), &pb.CreateSecretRequest{Name: "Name", Content: []byte{}})
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("SecretAlreadyExists", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.UUID{},
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			CreateSecret(gomock.Any(), secret).
			Return(nil, storage.ErrSecretConflict)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.CreateSecret(
			context.Background(), &pb.CreateSecretRequest{Name: secret.Name, Content: secret.Content})
		checkErrorStatus(t, err, codes.AlreadyExists)
	})
	t.Run("StorageError", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.UUID{},
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			CreateSecret(gomock.Any(), secret).
			Return(nil, errors.New("some error"))

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.CreateSecret(
			context.Background(), &pb.CreateSecretRequest{Name: secret.Name, Content: secret.Content})
		checkErrorStatus(t, err, codes.Internal)
	})
	t.Run("SuccessfulCreateSecret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.UUID{},
			OwnerID: userID,
		}
		createdSecret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.New(),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			CreateSecret(gomock.Any(), secret).
			Return(createdSecret, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.CreateSecret(
			context.Background(), &pb.CreateSecretRequest{Name: secret.Name, Content: secret.Content})
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
		assert.Equal(t, createdSecret.Version.String(), resp.Version)
	})
}

func TestSecretService_UpdateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secretStorage := ms.NewMockSecretStorage(ctrl)
	tokenManager := mt.NewMockManager(ctrl)

	secretService := &SecretService{
		SecretStorage: secretStorage,
	}

	interceptor := serverInterceptors.NewAuthInterceptor(tokenManager)

	server := NewServer(
		address,
		WithServices(secretService),
		WithUnaryInterceptors(interceptor.Unary()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Run(ctx)

	accessToken := "Token"
	userID := 0

	t.Run("EmptySecretName", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    "",
				Content: []byte("Content"),
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("EmptySecretContent", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    "SecretName",
				Content: []byte{},
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("SecretNotExists", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			UpdateSecret(gomock.Any(), secret).
			Return(nil, storage.ErrSecretNotFound)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
			},
		)
		checkErrorStatus(t, err, codes.NotFound)
	})
	t.Run("StorageError", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			UpdateSecret(gomock.Any(), secret).
			Return(nil, errors.New("some error"))

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
			},
		)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("SuccessfulUpdateSecret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			OwnerID: userID,
		}
		secretUpdated := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.New(),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			UpdateSecret(gomock.Any(), secret).
			Return(secretUpdated, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
		assert.Equal(t, secretUpdated.Version.String(), resp.Version)
	})
}

func TestSecretService_DeleteSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secretStorage := ms.NewMockSecretStorage(ctrl)
	tokenManager := mt.NewMockManager(ctrl)

	secretService := &SecretService{
		SecretStorage: secretStorage,
	}

	interceptor := serverInterceptors.NewAuthInterceptor(tokenManager)

	server := NewServer(
		address,
		WithServices(secretService),
		WithUnaryInterceptors(interceptor.Unary()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Run(ctx)

	accessToken := "Token"
	userID := 0

	t.Run("EmptySecretName", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.DeleteSecret(
			context.Background(),
			&pb.DeleteSecretRequest{
				Name: "",
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("StorageError", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			DeleteSecret(gomock.Any(), secret).
			Return(errors.New("some error"))

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.DeleteSecret(
			context.Background(),
			&pb.DeleteSecretRequest{
				Name: secret.Name,
			},
		)
		checkErrorStatus(t, err, codes.Internal)
	})
	t.Run("SuccessfulDeleteSecret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			DeleteSecret(gomock.Any(), secret).
			Return(nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.DeleteSecret(
			context.Background(),
			&pb.DeleteSecretRequest{
				Name: secret.Name,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
	})
}

func TestSecretService_ListSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secretStorage := ms.NewMockSecretStorage(ctrl)
	tokenManager := mt.NewMockManager(ctrl)

	secretService := &SecretService{
		SecretStorage: secretStorage,
	}

	interceptor := serverInterceptors.NewAuthInterceptor(tokenManager)

	server := NewServer(
		address,
		WithServices(secretService),
		WithUnaryInterceptors(interceptor.Unary()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Run(ctx)

	accessToken := "Token"
	userID := 0

	t.Run("StorageError", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		secretStorage.
			EXPECT().
			ListSecrets(gomock.Any(), userID).
			Return(nil, errors.New("some error"))

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.ListSecrets(
			context.Background(),
			&pb.ListSecretsRequest{},
		)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("SuccessfulListSecrets", func(t *testing.T) {
		secrets := []*models.Secret{
			{
				Name:    "Name1",
				Content: []byte("Content1"),
				Version: uuid.New(),
			},
			{
				Name:    "Name2",
				Content: []byte("Content2"),
				Version: uuid.New(),
			},
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		secretStorage.
			EXPECT().
			ListSecrets(gomock.Any(), userID).
			Return(secrets, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.ListSecrets(
			context.Background(),
			&pb.ListSecretsRequest{},
		)
		assert.NoError(t, err)
		assert.Equal(t, len(secrets), len(resp.Secrets))
		for i, secret := range secrets {
			assert.Equal(t, secret.Name, resp.Secrets[i].Name)
			assert.Equal(t, secret.Content, resp.Secrets[i].Content)
			assert.Equal(t, secret.Version.String(), resp.Secrets[i].Version)
		}
	})
}
