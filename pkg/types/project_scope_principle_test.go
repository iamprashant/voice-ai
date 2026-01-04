// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales/rapida.ai for commercial usage.
package types

import (
	"testing"

	type_enums "github.com/rapidaai/pkg/types/enums"
)

func TestProjectScope_GetUserId(t *testing.T) {
	ss := &ProjectScope{}
	if ss.GetUserId() != nil {
		t.Errorf("GetUserId() = %v, want nil", ss.GetUserId())
	}
}

func TestProjectScope_GetCurrentProjectId(t *testing.T) {
	projectId := uint64(2)
	ss := &ProjectScope{ProjectId: &projectId}
	if ss.GetCurrentProjectId() == nil || *ss.GetCurrentProjectId() != projectId {
		t.Errorf("GetCurrentProjectId() = %v, want %v", ss.GetCurrentProjectId(), projectId)
	}
}

func TestProjectScope_GetCurrentOrganizationId(t *testing.T) {
	orgId := uint64(1)
	ss := &ProjectScope{OrganizationId: &orgId}
	if ss.GetCurrentOrganizationId() == nil || *ss.GetCurrentOrganizationId() != orgId {
		t.Errorf("GetCurrentOrganizationId() = %v, want %v", ss.GetCurrentOrganizationId(), orgId)
	}
}

func TestProjectScope_HasOrganization(t *testing.T) {
	tests := []struct {
		name string
		ss   *ProjectScope
		want bool
	}{
		{"has org", &ProjectScope{OrganizationId: &[]uint64{1}[0]}, true},
		{"no org", &ProjectScope{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.HasOrganization(); got != tt.want {
				t.Errorf("HasOrganization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectScope_HasUser(t *testing.T) {
	ss := &ProjectScope{}
	if ss.HasUser() {
		t.Errorf("HasUser() = %v, want false", ss.HasUser())
	}
}

func TestProjectScope_HasProject(t *testing.T) {
	tests := []struct {
		name string
		ss   *ProjectScope
		want bool
	}{
		{"has project", &ProjectScope{ProjectId: &[]uint64{2}[0]}, true},
		{"no project", &ProjectScope{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.HasProject(); got != tt.want {
				t.Errorf("HasProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectScope_IsActive(t *testing.T) {
	tests := []struct {
		name string
		ss   *ProjectScope
		want bool
	}{
		{"active", &ProjectScope{Status: type_enums.RECORD_ACTIVE.String()}, true},
		{"inactive", &ProjectScope{Status: type_enums.RECORD_INACTIVE.String()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectScope_IsAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		ss   *ProjectScope
		want bool
	}{
		{"authenticated", &ProjectScope{ProjectId: &[]uint64{2}[0], OrganizationId: &[]uint64{1}[0], Status: type_enums.RECORD_ACTIVE.String()}, true},
		{"no project", &ProjectScope{OrganizationId: &[]uint64{1}[0], Status: type_enums.RECORD_ACTIVE.String()}, false},
		{"no org", &ProjectScope{ProjectId: &[]uint64{2}[0], Status: type_enums.RECORD_ACTIVE.String()}, false},
		{"inactive", &ProjectScope{ProjectId: &[]uint64{2}[0], OrganizationId: &[]uint64{1}[0], Status: type_enums.RECORD_INACTIVE.String()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.IsAuthenticated(); got != tt.want {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectScope_GetCurrentToken(t *testing.T) {
	token := "token"
	ss := &ProjectScope{CurrentToken: token}
	if ss.GetCurrentToken() != token {
		t.Errorf("GetCurrentToken() = %v, want %v", ss.GetCurrentToken(), token)
	}
}

func TestProjectScope_Type(t *testing.T) {
	ss := &ProjectScope{}
	if ss.Type() != "project" {
		t.Errorf("Type() = %v, want %v", ss.Type(), "project")
	}
}
