package api

import (
	"gritter/pkg/mission"
	"gritter/pkg/user"
)

type missionRepr struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func missionToRepr(m *mission.Mission) *missionRepr {
	return &missionRepr{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	}
}

type userRepr struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Intro string `json:"intro"`
}

func userToRepr(u *user.User) *userRepr {
	return &userRepr{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
		Intro: u.Intro,
	}
}
