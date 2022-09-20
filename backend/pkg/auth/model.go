package auth

import (
	"encoding/json"

	"go.uber.org/zap/zapcore"
	"google.golang.org/api/idtoken"
)

type Type int

const (
	TypeGoogle Type = iota + 1
)

type Info struct {
	Type Type `json:"type"`

	Google *InfoGoogle `json:"google"`
}

type InfoGoogle struct {
	IdToken string `json:"idToken"`
}

type Result struct {
	Type Type

	Google *ResultGoogle
}

func (r *Result) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("type", int(r.Type))
	encoder.AddObject("google", r.Google)
	return nil
}

type ResultGoogle struct {
	Payload *idtoken.Payload
}

func (rg *ResultGoogle) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	bytes, err := json.Marshal(rg.Payload)
	if err != nil {
		return err
	}

	encoder.AddString("payload", string(bytes))
	return nil
}
