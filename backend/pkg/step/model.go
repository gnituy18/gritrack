package step

import (
	"go.uber.org/zap/zapcore"
)

type Step struct {
	Id        string `bson:"id"`
	MissionId string `bson:"missionId"`

	Summary   string `bson:"summary"`
	Items     Items  `bson:"items"`

	CreatedAt int64 `bson:"createdAt"`
}

func (s *Step) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", s.Id)
	encoder.AddString("summary", s.Summary)
	encoder.AddArray("items", s.Items)
	encoder.AddInt64("createdAt", s.CreatedAt)
	return nil
}

type ItemType int

const (
	ItemTypeTime ItemType = iota + 1
)

type Items []*Item

func (items Items) MarshalLogArray(encoder zapcore.ArrayEncoder) error {
	for _, i := range items {
		encoder.AppendObject(i)
	}
	return nil
}

type Item struct {
	Type      ItemType `bson:"type"`
	Desc      string   `bson:"desc"`
	CreatedAt int64    `bson:"createdAt"`

	Time *ItemTime `bson:"time"`
}

func (i *Item) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("type", int(i.Type))
	encoder.AddString("desc", i.Desc)
	encoder.AddInt64("createdAt", i.CreatedAt)
	encoder.AddObject("time", i.Time)
	return nil
}

type ItemTime struct {
	Duration int `bson:"duration"`
}

func (it *ItemTime) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("duration", it.Duration)
	return nil
}
