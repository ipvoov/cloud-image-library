package user

import (
	"cloud/internal/service"
)

func init() {
	service.RegisterUser(New())
}

type sUser struct{}

func New() *sUser {
	return &sUser{}
}
