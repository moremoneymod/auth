package access_v1

import (
	"context"

	"github.com/moremoneymod/auth/pkg/access_v1"
)

type AccessService interface {
	Check(ctx context.Context, endpointAddress string) (bool, error)
}

type Implementation struct {
	access_v1.UnimplementedAccessV1Server
	accessService AccessService
}

func NewImplementation(service AccessService) *Implementation {
	return &Implementation{
		accessService: service,
	}
}
