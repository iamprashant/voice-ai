// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import (
	"strconv"
	"strings"
)

const (
	VERSION_PREFIX = "vrsn_"
)

func GetVersionDefinition(version string) *uint64 {
	if version == "" || version == "latest" {
		return nil
	}
	_vrsn := strings.Replace(version, VERSION_PREFIX, "", 1)
	_pid, err := strconv.ParseUint(_vrsn, 10, 64)
	if err != nil {
		return nil
	}
	return &_pid
}
