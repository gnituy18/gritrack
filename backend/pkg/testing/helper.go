package testing

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

var (
	ErrAny = errors.New("any err")
)

func Equal(t *testing.T, exp interface{}, res interface{}, msg string) {
	if !gomock.Eq(exp).Matches(res) {
		t.Errorf("%s: \n%v != %v", msg, exp, res)
	}
}
