//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"github.com/Littlefisher619/cosdisk/graph/auth"
	"github.com/Littlefisher619/cosdisk/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	Service    *service.CosDisk
	JwtManager *auth.JwtManager
}
