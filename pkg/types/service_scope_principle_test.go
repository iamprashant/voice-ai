// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestServiceScope_GetUserId(t *testing.T) {
	ss := &ServiceScope{UserId: &[]uint64{123}[0]}
	if got := ss.GetUserId(); got == nil || *got != 123 {
		t.Errorf("GetUserId() = %v, want %v", got, 123)
	}
}

func TestServiceScope_GetCurrentProjectId(t *testing.T) {
	ss := &ServiceScope{ProjectId: &[]uint64{456}[0]}
	if got := ss.GetCurrentProjectId(); got == nil || *got != 456 {
		t.Errorf("GetCurrentProjectId() = %v, want %v", got, 456)
	}
}

func TestServiceScope_GetCurrentOrganizationId(t *testing.T) {
	ss := &ServiceScope{OrganizationId: &[]uint64{789}[0]}
	if got := ss.GetCurrentOrganizationId(); got == nil || *got != 789 {
		t.Errorf("GetCurrentOrganizationId() = %v, want %v", got, 789)
	}
}

func TestServiceScope_HasOrganization(t *testing.T) {
	tests := []struct {
		name string
		ss   *ServiceScope
		want bool
	}{
		{
			name: "has org",
			ss:   &ServiceScope{OrganizationId: &[]uint64{1}[0]},
			want: true,
		},
		{
			name: "no org",
			ss:   &ServiceScope{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.HasOrganization(); got != tt.want {
				t.Errorf("HasOrganization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceScope_HasUser(t *testing.T) {
	tests := []struct {
		name string
		ss   *ServiceScope
		want bool
	}{
		{
			name: "has user",
			ss:   &ServiceScope{UserId: &[]uint64{1}[0]},
			want: true,
		},
		{
			name: "no user",
			ss:   &ServiceScope{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.HasUser(); got != tt.want {
				t.Errorf("HasUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceScope_HasProject(t *testing.T) {
	tests := []struct {
		name string
		ss   *ServiceScope
		want bool
	}{
		{
			name: "has project",
			ss:   &ServiceScope{ProjectId: &[]uint64{1}[0]},
			want: true,
		},
		{
			name: "no project",
			ss:   &ServiceScope{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.HasProject(); got != tt.want {
				t.Errorf("HasProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceScope_IsAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		ss   *ServiceScope
		want bool
	}{
		{
			name: "authenticated with user and org",
			ss: &ServiceScope{
				UserId:         &[]uint64{1}[0],
				OrganizationId: &[]uint64{2}[0],
			},
			want: true,
		},
		{
			name: "authenticated with project and org",
			ss: &ServiceScope{
				ProjectId:      &[]uint64{1}[0],
				OrganizationId: &[]uint64{2}[0],
			},
			want: true,
		},
		{
			name: "not authenticated no org",
			ss: &ServiceScope{
				UserId: &[]uint64{1}[0],
			},
			want: false,
		},
		{
			name: "not authenticated no user or project",
			ss: &ServiceScope{
				OrganizationId: &[]uint64{2}[0],
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.IsAuthenticated(); got != tt.want {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceScope_GetCurrentToken(t *testing.T) {
	ss := &ServiceScope{CurrentToken: "token123"}
	if got := ss.GetCurrentToken(); got != "token123" {
		t.Errorf("GetCurrentToken() = %v, want %v", got, "token123")
	}
}

func TestServiceScope_Type(t *testing.T) {
	ss := &ServiceScope{}
	if got := ss.Type(); got != "service" {
		t.Errorf("Type() = %v, want %v", got, "service")
	}
}