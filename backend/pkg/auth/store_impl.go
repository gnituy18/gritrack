package auth

import (
	"net/http"
	"os"

	"go.uber.org/zap"
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
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(im.httpClient))
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(info.AccessToken)
	tokenInfo, err := tokenInfoCall.Do()
	if tokenInfo.Audience != os.Getenv("GOOGLE_CLIENT_ID") {
		return nil, ErrTokenAudienceInvalid
	} else if err != nil {
		return nil, err
	}

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
