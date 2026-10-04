package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Antesser/trainingFinder/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	// errNoMetadata returned when no metadata in request context.
	errNoMetadata = status.Error(codes.Unauthenticated, "no metadata in the request context")
	// errKeyNotFound returned when authorization key not found in metadata.
	errKeyNotFound = status.Error(codes.Unauthenticated, "authorization token is empty")
	// errInvalidAuthHeader return when auth header have wrong format.
	errInvalidAuthHeader = status.Error(codes.Unauthenticated, "invalid authorization header format")
	// errInvalidToken returned when auth token not match.
	errInvalidToken = status.Error(codes.Unauthenticated, "invalid token")
	// errClientUnauthenticated
	errClientUnauthenticated = status.Error(codes.Unauthenticated, "client is not authenticated to use handler")
	// errUnexpectedSigningMethod
	errUnexpectedSigningMethod = status.Error(codes.Unauthenticated, "unexpected signing method")
	// errFailedToDecodeClaims
	errFailedToDecodeClaims = status.Error(codes.Unauthenticated, "failed to decode token claims")
	errAccessDenied         = status.Error(codes.Unauthenticated, "failed to find a proper role")
)

type userIDKeyType struct{}

var userIDKey userIDKeyType
var coachIDKey userIDKeyType

func WithAuth(secretKey string, cfg config.AuthConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		allowedRoles, needsAuth := cfg.BearerSet[info.FullMethod]
		if !needsAuth {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errNoMetadata
		}

		authMD := md.Get("authorization")
		if len(authMD) == 0 {
			return nil, errKeyNotFound
		}

		authHeader := authMD[0]
		const bearerPrefix = "Bearer "
		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if len(token) == 0 {
			return nil, errInvalidAuthHeader
		}

		jwToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errUnexpectedSigningMethod
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			return nil, errClientUnauthenticated
		}
		if !jwToken.Valid {
			return nil, errInvalidToken
		}
		claims, ok := jwToken.Claims.(jwt.MapClaims)
		fmt.Println("claims", claims)
		if !ok {
			return nil, errFailedToDecodeClaims
		}
		userRole, ok := claims["role"]
		if !ok {
			fmt.Println("userRole", claims["role"])
		}
		fmt.Println("userRole", userRole)
		userRoleID, ok := userRole.(string)
		if !ok {
			fmt.Println("userRoleID", userRoleID)
		}
		if !hasAllowedRole(userRoleID, allowedRoles) {
			return nil, errAccessDenied
		}
		id := claims["id"]

		switch {
		case IsCoach(userRoleID):
			ctx = context.WithValue(ctx, coachIDKey, id)
		case IsUser(userRoleID):
			ctx = context.WithValue(ctx, userIDKey, id)
		}

		return handler(ctx, req)
	}
}

func GetUserID(ctx context.Context) *string {
	res, ok := ctx.Value(userIDKey).(string)
	if ok {
		return &res
	}

	return nil
}
func GetCoachID(ctx context.Context) *string {
	res, ok := ctx.Value(coachIDKey).(string)
	if ok {
		return &res
	}

	return nil
}

func hasAllowedRole(userRole string, allowed map[string]struct{}) bool {
	role := roles[userRole]
	if _, ok := allowed[role]; ok {
		return ok
	}

	return false
}
