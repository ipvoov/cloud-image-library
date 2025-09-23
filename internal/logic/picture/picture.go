package picture

import (
	"cloud/internal/service"
)

func init() {
	service.RegisterPicture(New())
}

type sPicture struct{}

func New() *sPicture {
	return &sPicture{}
}
