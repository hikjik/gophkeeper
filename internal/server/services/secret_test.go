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
	"github.com/go-developer-ya-practicum/gophkeeper/internal/server/token"
	mt "github.com/go-developer-ya-practicum/gophkeeper/internal/server/token/mock"
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
			Return(uuid.Nil, storage.ErrSecretNameConflict)

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
			Return(uuid.Nil, errors.New("some error"))

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

		version := uuid.New()

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			CreateSecret(gomock.Any(), secret).
			Return(version, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.CreateSecret(
			context.Background(), &pb.CreateSecretRequest{Name: secret.Name, Content: secret.Content})
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
		assert.Equal(t, version.String(), resp.Version)
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
				Version: uuid.New().String(),
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
				Version: uuid.New().String(),
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("InvalidSecretVersion", func(t *testing.T) {
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
				Content: []byte("Content"),
				Version: "Invalid",
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("SecretNotExists", func(t *testing.T) {
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
			UpdateSecret(gomock.Any(), secret).
			Return(uuid.Nil, storage.ErrSecretNotFound)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
				Version: secret.Version.String(),
			},
		)
		checkErrorStatus(t, err, codes.NotFound)
	})
	t.Run("SecretVersionConflict", func(t *testing.T) {
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
			UpdateSecret(gomock.Any(), secret).
			Return(uuid.Nil, storage.ErrSecretVersionConflict)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
				Version: secret.Version.String(),
			},
		)
		checkErrorStatus(t, err, codes.NotFound)
	})
	t.Run("StorageError", func(t *testing.T) {
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
			UpdateSecret(gomock.Any(), secret).
			Return(uuid.Nil, errors.New("some error"))

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
				Version: secret.Version.String(),
			},
		)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("SuccessfulUpdateSecret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Content: []byte("SecretContent"),
			Version: uuid.New(),
			OwnerID: userID,
		}

		version := uuid.New()

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			UpdateSecret(gomock.Any(), secret).
			Return(version, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		resp, err := client.UpdateSecret(
			context.Background(),
			&pb.UpdateSecretRequest{
				Name:    secret.Name,
				Content: secret.Content,
				Version: secret.Version.String(),
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
		assert.Equal(t, version.String(), resp.Version)
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
				Name:    "",
				Version: uuid.New().String(),
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("InvalidSecretVersion", func(t *testing.T) {
		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: userID}, nil)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.DeleteSecret(
			context.Background(),
			&pb.DeleteSecretRequest{
				Name:    "SecretName",
				Version: "Invalid",
			},
		)
		checkErrorStatus(t, err, codes.InvalidArgument)
	})
	t.Run("SecretNotExists", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Version: uuid.New(),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			DeleteSecret(gomock.Any(), secret).
			Return(storage.ErrSecretNotFound)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.DeleteSecret(
			context.Background(),
			&pb.DeleteSecretRequest{
				Name:    secret.Name,
				Version: secret.Version.String(),
			},
		)
		checkErrorStatus(t, err, codes.NotFound)
	})
	t.Run("SecretVersionConflict", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Version: uuid.New(),
			OwnerID: userID,
		}

		tokenManager.
			EXPECT().
			Validate(accessToken).
			Return(&token.Payload{UserID: secret.OwnerID}, nil)

		secretStorage.
			EXPECT().
			DeleteSecret(gomock.Any(), secret).
			Return(storage.ErrSecretVersionConflict)

		client, err := newSecretClient(accessToken)
		require.NoError(t, err)

		_, err = client.DeleteSecret(
			context.Background(),
			&pb.DeleteSecretRequest{
				Name:    secret.Name,
				Version: secret.Version.String(),
			},
		)
		checkErrorStatus(t, err, codes.NotFound)
	})
	t.Run("StorageError", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Version: uuid.New(),
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
				Name:    secret.Name,
				Version: secret.Version.String(),
			},
		)
		checkErrorStatus(t, err, codes.Internal)
	})
	t.Run("SuccessfulDeleteSecret", func(t *testing.T) {
		secret := &models.Secret{
			Name:    "SecretName",
			Version: uuid.New(),
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
				Name:    secret.Name,
				Version: secret.Version.String(),
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, secret.Name, resp.Name)
	})
}

func TestSecretService_FetchSecrets(t *testing.T) {
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

		_, err = client.FetchSecrets(
			context.Background(),
			&pb.FetchSecretsRequest{},
		)
		checkErrorStatus(t, err, codes.Internal)
	})

	t.Run("SuccessfulFetchSecrets", func(t *testing.T) {
		secrets := []*models.Secret{
			{
				Name:    "Name1",
				Version: uuid.New(),
			},
			{
				Name:    "Name2",
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

		resp, err := client.FetchSecrets(
			context.Background(),
			&pb.FetchSecretsRequest{},
		)
		assert.NoError(t, err)
		assert.Equal(t, len(secrets), len(resp.Secrets))
		for i, secret := range secrets {
			assert.Equal(t, secret.Name, resp.Secrets[i].Name)
			assert.Equal(t, secret.Version.String(), resp.Secrets[i].Version)
		}
	})
}
