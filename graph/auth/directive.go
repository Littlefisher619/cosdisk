package auth

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (j *JwtManager) AuthDirective(
	ctx context.Context, obj interface{}, next graphql.Resolver,
) (interface{}, error) {
	jwtAuthToken := ExtractJwtFromContext(ctx)
	if jwtAuthToken == "" {
		return nil, gqlerror.Errorf("Access denied: Login required")
	}

	token, err := j.Validate(ctx, jwtAuthToken)
	if err != nil || !token.Valid {
		return nil, gqlerror.Errorf("Access denied: %s", err)
	}

	customClaim, _ := token.Claims.(*JwtCustomClaim)
	ctx = context.WithValue(ctx, contextKeyUID, customClaim.Id)
	return next(ctx)
}
