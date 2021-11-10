package mission

import (
	"go.uber.org/zap/zapcore"
)

type Mission struct {
	Id     string `bson:"id"`
	UserId string `bson:"userId"`

	Name        string `bson:"name"`
	Description string `bson:"description"`

	Deleted bool `bson:"deleted"`
}

func (m *Mission) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", m.Id)
	encoder.AddString("name", m.Name)
	encoder.AddString("userId", m.UserId)
	encoder.AddString("description", m.Description)
	encoder.AddBool("deleted", m.Deleted)
	return nil
}
