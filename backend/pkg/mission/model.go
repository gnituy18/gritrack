package mission

import (
	"go.uber.org/zap/zapcore"
)

type Mission struct {
	ID     string `bson:"id"`
	UserID string `bson:"userID"`

	Name        string `bson:"name"`
	Description string `bson:"description"`

	Deleted bool `bson:"deleted"`
}

func (m *Mission) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", m.ID)
	encoder.AddString("name", m.Name)
	encoder.AddString("userID", m.UserID)
	encoder.AddString("description", m.Description)
	encoder.AddBool("deleted", m.Deleted)
	return nil
}
