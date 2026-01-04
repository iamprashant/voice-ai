// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with RapidaAI
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestCreateServiceScopeToken(t *testing.T) {
	secretKey := "test-secret"

	tests := []struct {
		name      string
		principle SimplePrinciple
		wantErr   bool
	}{
		{
			name: "valid principle with all fields",
			principle: &ServiceScope{
				UserId:         &[]uint64{1}[0],
				OrganizationId: &[]uint64{2}[0],
				ProjectId:      &[]uint64{3}[0],
			},
			wantErr: false,
		},
		{
			name: "principle with only user",
			principle: &ServiceScope{
				UserId: &[]uint64{1}[0],
			},
			wantErr: false,
		},
		{
			name: "principle with only organization",
			principle: &ServiceScope{
				OrganizationId: &[]uint64{2}[0],
			},
			wantErr: false,
		},
		{
			name: "principle with only project",
			principle: &ServiceScope{
				ProjectId: &[]uint64{3}[0],
			},
			wantErr: false,
		},
		{
			name:      "nil principle",
			principle: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateServiceScopeToken(tt.principle, secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateServiceScopeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Errorf("CreateServiceScopeToken() returned empty token")
			}
		})
	}
}

func TestExtractServiceScope(t *testing.T) {
	secretKey := "test-secret"
	principle := &ServiceScope{
		UserId:         &[]uint64{1}[0],
		OrganizationId: &[]uint64{2}[0],
		ProjectId:      &[]uint64{3}[0],
	}
	validToken, _ := CreateServiceScopeToken(principle, secretKey)

	tests := []struct {
		name     string
		token    string
		secret   string
		wantUser *uint64
		wantOrg  *uint64
		wantProj *uint64
		wantErr  bool
	}{
		{
			name:     "valid token",
			token:    validToken,
			secret:   secretKey,
			wantUser: &[]uint64{1}[0],
			wantOrg:  &[]uint64{2}[0],
			wantProj: &[]uint64{3}[0],
			wantErr:  false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			secret:  secretKey,
			wantErr: true,
		},
		{
			name:    "wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			secret:  secretKey,
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			secret:  secretKey,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope, err := ExtractServiceScope(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractServiceScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if scope == nil {
					t.Errorf("ExtractServiceScope() returned nil scope")
					return
				}
				if tt.wantUser != nil && (scope.UserId == nil || *scope.UserId != *tt.wantUser) {
					t.Errorf("ExtractServiceScope() userId = %v, want %v", scope.UserId, tt.wantUser)
				}
				if tt.wantOrg != nil && (scope.OrganizationId == nil || *scope.OrganizationId != *tt.wantOrg) {
					t.Errorf("ExtractServiceScope() orgId = %v, want %v", scope.OrganizationId, tt.wantOrg)
				}
				if tt.wantProj != nil && (scope.ProjectId == nil || *scope.ProjectId != *tt.wantProj) {
					t.Errorf("ExtractServiceScope() projId = %v, want %v", scope.ProjectId, tt.wantProj)
				}
			}
		})
	}
}

func TestToUint64(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  uint64
		ok    bool
	}{
		{
			name:  "float64",
			value: float64(123),
			want:  123,
			ok:    true,
		},
		{
			name:  "int",
			value: 456,
			want:  456,
			ok:    true,
		},
		{
			name:  "int64",
			value: int64(789),
			want:  789,
			ok:    true,
		},
		{
			name:  "string valid",
			value: "101112",
			want:  101112,
			ok:    true,
		},
		{
			name:  "string invalid",
			value: "not-a-number",
			want:  0,
			ok:    false,
		},
		{
			name:  "unsupported type",
			value: true,
			want:  0,
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := toUint64(tt.value)
			if ok != tt.ok {
				t.Errorf("toUint64() ok = %v, want %v", ok, tt.ok)
			}
			if ok && got != tt.want {
				t.Errorf("toUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
