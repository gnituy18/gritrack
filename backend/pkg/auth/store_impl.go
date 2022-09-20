package auth

import (
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/api/idtoken"

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
	payload, err := idtoken.Validate(ctx, info.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, ErrTokenAudienceInvalid
	}

	return &Result{
		Type: TypeGoogle,
		Google: &ResultGoogle{
			Payload: payload,
		},
	}, nil
}
