package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/FarhanDH/sidingin-rating-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type AccessClaims struct {
	UserID   string          `json:"sub"`
	FullName string          `json:"fullName"`
	Role     domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type JwtVerifier struct {
	secret []byte
	issuer string
}

func NewJwtVerifier(secret, issuer string) (*JwtVerifier, error) {
	if secret == "" || issuer == "" {
		return nil, errors.New("JWT secret and issuer are required")
	}
	return &JwtVerifier{
		secret: []byte(secret),
		issuer: issuer,
	}, nil
}

func (v *JwtVerifier) VerifyIncoming(ctx context.Context) (*AccessClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing gRPC metadata")
	}

	values := md.Get("authorization")
	if len(values) != 1 {
		return nil, errors.New("mission authorization metadata")
	}

	parts := strings.Fields(values[0])
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid authorization format")
	}

	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return v.secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience("rating-service"),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("verify access token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return claims, nil
}
