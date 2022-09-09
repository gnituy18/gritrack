package api

import (
	"gritter/pkg/mission"
	"gritter/pkg/step"
	"gritter/pkg/user"
)

type missionRepr struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func missionToRepr(m *mission.Mission) *missionRepr {
	return &missionRepr{
		Id:          m.Id,
		UserId:      m.UserId,
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

type stepsRepr struct {
	Steps []*stepRepr `json:"steps"`
	More  bool        `json:"more"`
}

func stepsToRepr(ss []*step.Step, more bool) *stepsRepr {
	steps := make([]*stepRepr, len(ss))
	for i, s := range ss {
		steps[i] = stepToRepr(s)
	}
	return &stepsRepr{
		Steps: steps,
		More:  more,
	}
}

type stepRepr struct {
	Id        string     `json:"id"`
	Time      int64      `json:"time"`
	Summary   string     `json:"summary"`
	Items     step.Items `json:"items"`
	CreatedAt int64      `json:"createdAt"`
}

func stepToRepr(s *step.Step) *stepRepr {
	return &stepRepr{
		Id:        s.Id,
		Time:      s.Time,
		Summary:   s.Summary,
		Items:     s.Items,
		CreatedAt: s.CreatedAt,
	}
}
