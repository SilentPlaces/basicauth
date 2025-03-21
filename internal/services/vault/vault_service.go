package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"github.com/google/wire"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
	"sync"
)

type SecureVaultService interface {
	GetJWTConfig() (*VaultJWTSecretConfig, error)
}

type vaultService struct {
	client *vault.Client
}

type VaultJWTSecretConfig struct {
	JwtSecret        []byte
	JwtRefreshSecret []byte
}

var (
	service *vaultService
	once    sync.Once
)

func NewSecureVaultService() SecureVaultService {
	once.Do(func() {
		vaultAddr := os.Getenv(constants.EnvKeyVaultAddr)
		vaultToken := os.Getenv(constants.EnvKeyVaultToken)
		if vaultAddr == "" || vaultToken == "" {
			log.Fatal("Vault address or token is not set in environment variables")
		}

		config := vault.DefaultConfig()
		config.Address = vaultAddr
		client, err := vault.NewClient(config)
		if err != nil {
			log.Fatalf("Unable to initialize Vault client: %v", err)
		}
		client.SetToken(vaultToken)
		service = &vaultService{client: client}
	})

	return service
}
func (s *vaultService) GetJWTConfig() (*VaultJWTSecretConfig, error) {
	path := os.Getenv(constants.EnvKeyVaultJWTPath)
	fmt.Println("####################", path)

	secret, err := s.client.KVv2("secret").Get(context.Background(), "jwt")

	if err != nil {
		log.Printf("Error reading secret from Vault at %s: %v", path, err)
		return nil, fmt.Errorf("unable to read secret from Vault: %w", err)
	}
	if secret == nil || secret.Data == nil {
		err := errors.New("no data found at the provided path")
		log.Printf("Error: %v", err)
		return nil, err
	}

	// Extract the secrets from the returned data.
	jwtSecret, ok := secret.Data[constants.VaultJWTSecretKey].(string)
	if !ok {
		err := errors.New("jwtSecret not found or not a string")
		log.Printf("Error: %v", err)
		return nil, err
	}

	jwtRefreshSecret, ok := secret.Data[constants.VaultJWTRefreshSecretKey].(string)
	if !ok {
		err := errors.New("jwtRefreshSecret not found or not a string")
		log.Printf("Error: %v", err)
		return nil, err
	}

	return &VaultJWTSecretConfig{
		JwtSecret:        []byte(jwtSecret),
		JwtRefreshSecret: []byte(jwtRefreshSecret),
	}, nil
}

var VaultServiceProviderSet = wire.NewSet(NewSecureVaultService)
