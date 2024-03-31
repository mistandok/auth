package access

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mistandok/auth/pkg/access_v1"
	"google.golang.org/grpc/metadata"
)

const msgInternalError = "что-то пошло не так, мы уже работаем над решением проблемы"

var errInternal = errors.New(msgInternalError)

// Implementation user Server.
type Implementation struct {
	access_v1.UnimplementedAccessV1Server
}

// NewImplementation ..
func NewImplementation() *Implementation {
	return &Implementation{}
}

// Create ..
func (i *Implementation) Create(ctx context.Context, request *access_v1.CreateRequest) (*access_v1.CreateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], "Bearer ") {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], "Bearer ")

	return &access_v1.CreateResponse{Id: 1}, errors.New(fmt.Sprintf("%s", accessToken))
}
