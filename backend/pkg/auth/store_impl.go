package auth

import (
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/api/oauth2/v2"

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
	oauth2Service, err := oauth2.New(im.httpClient)
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(info.IDToken).Do()
	if err != nil {
		return nil, err
	}

	return &Result{
		Type: TypeGoogle,
		Google: &ResultGoogle{
			TokenInfo: tokenInfo,
		},
	}, nil
}
