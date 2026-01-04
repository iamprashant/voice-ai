// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateJWT creates a JWT token with the provided claims and returns the token string
func CreateServiceScopeToken(principle SimplePrinciple, secretKey string) (string, error) {
	if principle == nil {
		return "", fmt.Errorf("principle cannot be nil")
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours
	}

	if principle.GetUserId() != nil {
		claims["userId"] = *principle.GetUserId()
	}
	if principle.GetCurrentOrganizationId() != nil {
		claims["organizationId"] = *principle.GetCurrentOrganizationId()
	}
	if principle.GetCurrentProjectId() != nil {
		claims["projectId"] = *principle.GetCurrentProjectId()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error creating token: %v", err)
	}
	return tokenString, nil
}

// ExtractJWT extracts the claims from the provided JWT token string and returns the decoded PlainAuthPrinciple
func ExtractServiceScope(tokenString string, secretKey string) (*ServiceScope, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	ol := &ServiceScope{
		CurrentToken: tokenString,
	}

	if _, exists := claims["userId"]; exists {
		user, ok := toUint64(claims["userId"])
		if ok {
			ol.UserId = &user
		}
	}

	if _, exists := claims["organizationId"]; exists {
		organizationId, ok := toUint64(claims["organizationId"])
		if ok {
			ol.OrganizationId = &organizationId
		}

	}

	if _, exists := claims["projectId"]; exists {
		projectId, ok := toUint64(claims["projectId"])
		if ok {
			ol.ProjectId = &projectId
		}
	}
	return ol, nil
}

func toUint64(value interface{}) (uint64, bool) {
	switch v := value.(type) {
	case float64:
		return uint64(v), true
	case int:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case string:
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			return parsed, true
		}
		// Add more cases as needed
	}
	return 0, false
}
