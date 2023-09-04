package utils

import (
	"github.com/google/uuid"
)

type IUuidGen interface {
	New() string
}
type uuidgen struct{}

func New() IUuidGen {
	return &uuidgen{}
}

func (u *uuidgen) New() string {
	return uuid.New().String()
}

func Parse(value string) (string, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
