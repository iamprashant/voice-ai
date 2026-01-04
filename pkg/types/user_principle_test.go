// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestUserInfo_GetId(t *testing.T) {
	u := &UserInfo{Id: 123}
	if u.GetId() != 123 {
		t.Errorf("GetId() = %v, want %v", u.GetId(), 123)
	}
}

func TestUserInfo_GetName(t *testing.T) {
	u := &UserInfo{Name: "test"}
	if u.GetName() != "test" {
		t.Errorf("GetName() = %v, want %v", u.GetName(), "test")
	}
}

func TestUserInfo_GetEmail(t *testing.T) {
	u := &UserInfo{Email: "test@example.com"}
	if u.GetEmail() != "test@example.com" {
		t.Errorf("GetEmail() = %v, want %v", u.GetEmail(), "test@example.com")
	}
}

func TestProjectRole_GetRole(t *testing.T) {
	p := &ProjectRole{Role: "admin"}
	if p.GetRole() != "admin" {
		t.Errorf("GetRole() = %v, want %v", p.GetRole(), "admin")
	}
}

func TestProjectRole_GetProjectId(t *testing.T) {
	p := &ProjectRole{ProjectId: 456}
	if p.GetProjectId() != 456 {
		t.Errorf("GetProjectId() = %v, want %v", p.GetProjectId(), 456)
	}
}