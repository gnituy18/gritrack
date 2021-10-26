package mission

import (
	"testing"

	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"

	"gritter/pkg/context"
	"gritter/pkg/document"
	ptesting "gritter/pkg/testing"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDocument := document.NewMockDocument(ctrl)
	mockID := "mockID"

	cases := []struct {
		desc    string
		prepare func()
		expRes  *Mission
		expErr  error
	}{
		{
			desc: "document.GetOne failed",
			prepare: func() {
				mockDocument.EXPECT().GetOne(
					gomock.Any(),
					gomock.Eq(document.Mission),
					gomock.Eq(bson.M{"id": mockID}),
					gomock.Any(),
				).Return(ptesting.ErrAny)
			},
			expRes: nil,
			expErr: ptesting.ErrAny,
		},
		{
			desc: "mission not found",
			prepare: func() {
				mockDocument.EXPECT().GetOne(
					gomock.Any(),
					gomock.Eq(document.Mission),
					gomock.Eq(bson.M{"id": mockID}),
					gomock.Any(),
				).Return(document.ErrNotFound)
			},
			expRes: nil,
			expErr: ErrNotFound,
		},
		{
			desc: "success",
			prepare: func() {
				mockDocument.EXPECT().GetOne(
					gomock.Any(),
					gomock.Eq(document.Mission),
					gomock.Eq(bson.M{"id": mockID}),
					gomock.Any(),
				).Do(func(ctx context.Context, name document.Name, query bson.M, doc interface{}) {
					a := doc.(*Mission)
					a.ID = "id"
				}).Return(nil)
			},
			expRes: &Mission{
				ID: "id",
			},
			expErr: nil,
		},
	}

	for _, c := range cases {
		im := impl{
			doc: mockDocument,
		}
		c.prepare()
		res, err := im.Get(context.Background(), mockID)
		ptesting.Equal(t, c.expRes, res, c.desc)
		ptesting.Equal(t, c.expErr, err, c.desc)
	}
}
