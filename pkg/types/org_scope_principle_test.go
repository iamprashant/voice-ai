// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"

	type_enums "github.com/rapidaai/pkg/types/enums"
)

func TestOrganizationScope_GetUserId(t *testing.T) {
	ss := &OrganizationScope{}
	if ss.GetUserId() != nil {
		t.Errorf("GetUserId() = %v, want nil", ss.GetUserId())
	}
}

func TestOrganizationScope_GetCurrentProjectId(t *testing.T) {
	ss := &OrganizationScope{}
	if ss.GetCurrentProjectId() != nil {
		t.Errorf("GetCurrentProjectId() = %v, want nil", ss.GetCurrentProjectId())
	}
}

func TestOrganizationScope_GetCurrentOrganizationId(t *testing.T) {
	orgId := uint64(1)
	ss := &OrganizationScope{OrganizationId: &orgId}
	if ss.GetCurrentOrganizationId() == nil || *ss.GetCurrentOrganizationId() != orgId {
		t.Errorf("GetCurrentOrganizationId() = %v, want %v", ss.GetCurrentOrganizationId(), orgId)
	}
}

func TestOrganizationScope_HasOrganization(t *testing.T) {
	tests := []struct {
		name string
		ss   *OrganizationScope
		want bool
	}{
		{"has org", &OrganizationScope{OrganizationId: &[]uint64{1}[0]}, true},
		{"no org", &OrganizationScope{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.HasOrganization(); got != tt.want {
				t.Errorf("HasOrganization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrganizationScope_HasUser(t *testing.T) {
	ss := &OrganizationScope{}
	if ss.HasUser() {
		t.Errorf("HasUser() = %v, want false", ss.HasUser())
	}
}

func TestOrganizationScope_HasProject(t *testing.T) {
	ss := &OrganizationScope{}
	if ss.HasProject() {
		t.Errorf("HasProject() = %v, want false", ss.HasProject())
	}
}

func TestOrganizationScope_IsActive(t *testing.T) {
	tests := []struct {
		name string
		ss   *OrganizationScope
		want bool
	}{
		{"active", &OrganizationScope{Status: type_enums.RECORD_ACTIVE.String()}, true},
		{"inactive", &OrganizationScope{Status: type_enums.RECORD_INACTIVE.String()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrganizationScope_IsAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		ss   *OrganizationScope
		want bool
	}{
		{"authenticated", &OrganizationScope{OrganizationId: &[]uint64{1}[0], Status: type_enums.RECORD_ACTIVE.String()}, true},
		{"no org", &OrganizationScope{Status: type_enums.RECORD_ACTIVE.String()}, false},
		{"inactive", &OrganizationScope{OrganizationId: &[]uint64{1}[0], Status: type_enums.RECORD_INACTIVE.String()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.IsAuthenticated(); got != tt.want {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrganizationScope_GetCurrentToken(t *testing.T) {
	token := "token"
	ss := &OrganizationScope{CurrentToken: token}
	if ss.GetCurrentToken() != token {
		t.Errorf("GetCurrentToken() = %v, want %v", ss.GetCurrentToken(), token)
	}
}

func TestOrganizationScope_Type(t *testing.T) {
	ss := &OrganizationScope{}
	if ss.Type() != "organization" {
		t.Errorf("Type() = %v, want %v", ss.Type(), "organization")
	}
}
