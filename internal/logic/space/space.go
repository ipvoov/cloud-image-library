package space

import (
	"cloud/internal/service"
)

func init() {
	service.RegisterSpace(New())
}

type sSpace struct{}

func New() *sSpace {
	return &sSpace{}
}
