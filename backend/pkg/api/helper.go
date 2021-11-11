package api

import (
	"gritter/pkg/mission"
	"gritter/pkg/user"
)

type missionRepr struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func missionToRepr(m *mission.Mission) *missionRepr {
	return &missionRepr{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
	}
}

type userRepr struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
	Name    string `json:"name"`
	Intro   string `json:"intro"`
}

func userToRepr(u *user.User) *userRepr {
	return &userRepr{
		Id:      u.Id,
		Email:   u.Email,
		Name:    u.Name,
		Picture: u.Picture,
		Intro:   u.Intro,
	}
}
