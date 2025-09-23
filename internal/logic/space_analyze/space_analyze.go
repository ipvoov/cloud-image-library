package space_analyze

import (
	"cloud/internal/service"
)

func init() {
	service.RegisterSpaceAnalyze(New())
}

type sSpaceAnalyze struct{}

func New() *sSpaceAnalyze {
	return &sSpaceAnalyze{}
}
