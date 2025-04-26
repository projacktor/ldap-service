package identity

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type identityManager struct {
	baseUrl             string
	realm               string
	restApiClientId     string
	restApiClientSecret string
}

func NewIdentityManager() *identityManager {
	return &identityManager{
		baseUrl:             viper.GetString("KEYCLOAK_BASE_URL"),
		realm:               viper.GetString("KEYCLOAK_REALM"),
		restApiClientId:     viper.GetString("KEYCLOAK_REST_API_CLIENT_ID"),
		restApiClientSecret: viper.GetString("KEYCLOAK_REST_API_CLIENT_SECRET"),
	}
}

// Client logging in Keycloak
func (im *identityManager) loginRestApiClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)
	token, err := client.LoginClient(ctx, im.restApiClientId, im.restApiClientSecret, im.realm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login rest api client")
	}
	return token, nil
}

func (im *identityManager) CreateUser(ctx context.Context, user gocloak.User, password string) (*gocloak.User, error) {
	// client := gocloak.NewClient(im.baseUrl)
	// token, err := im.loginRestApiClient(ctx); if err != nil {
	// 	return nil, err
	// }
}
