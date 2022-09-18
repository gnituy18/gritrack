package auth

import (
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"gritter/pkg/context"
)

func New() Store {
	return &impl{
		httpClient: &http.Client{},
	}
}

type impl struct {
	httpClient *http.Client
}

func (im *impl) Auth(ctx context.Context, info *Info) (*Result, error) {
	switch info.Type {
	case TypeGoogle:
		return im.authWithGoogle(ctx, info.Google)
	default:
		ctx.With(
			zap.Int("type", int(info.Type)),
		).Error("auth type invalid")
		return nil, ErrTypeInvalid
	}
}

func (im *impl) authWithGoogle(ctx context.Context, info *InfoGoogle) (*Result, error) {
	// validate google id token
	_, err := idtoken.Validate(ctx, info.AccessToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, ErrTokenAudienceInvalid
	}

	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(im.httpClient))

	userInfoCall := oauth2Service.Userinfo.Get()
	userInfoCall.Header().Set("Authorization", "Bearer "+info.AccessToken)
	userInfo, err := userInfoCall.Do()
	if err != nil {
		return nil, err
	}

	return &Result{
		Type: TypeGoogle,
		Google: &ResultGoogle{
			UserInfo: userInfo,
		},
	}, nil
}
