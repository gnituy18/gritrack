package user

import (
	"go.uber.org/zap/zapcore"
)

type User struct {
	Id      string `bson:"id"`
	Email   string `bson:"email"`
	Name    string `bson:"name"`
	Picture string `bson:"picture"`
	Intro   string `bson:"intro"`

	GoogleUserId string `bson:"googleUserId"`

	Deleted bool
}

func (u *User) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", u.Id)
	encoder.AddString("email", u.Email)
	encoder.AddString("name", u.Name)
	encoder.AddString("picture", u.Picture)
	encoder.AddString("intro", u.Intro)
	encoder.AddBool("deleted", u.Deleted)
	return nil
}
