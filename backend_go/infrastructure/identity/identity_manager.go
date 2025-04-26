package identity

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"

	"BackendGoLdap/config"
)


// Client logging in Keycloak
func () loginRestApiClient(ctx context.Context) (*gocloak.JWT, error) {
	cfg, err := config.GetConfig()
	client := gocloak.NewClient(cfg.BaseUrl)
	token, err := client.LoginClient(ctx, cfg.RestApiClientId, cfg.RestApiClientSecret, cfg.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login rest api client")
	}
	return token, nil
}

func (im *identityManager) CreateUser(ctx context.Context, user gocloak.User, password string) (*gocloak.User, error) {
	client := gocloak.NewClient(im.baseUrl)
	token, err := im.loginRestApiClient(ctx); if err != nil {
		return nil, err
	}
}
