// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package roundrobin_keyrotators

import (
	"sync"
	"sync/atomic"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/keyrotators"
)

type roundRobinKeyRotator struct {
	keys     []string
	counter  uint32
	keyMutex sync.RWMutex
	logger   commons.Logger
}

func NewRoundRobinKeyRotator(logger commons.Logger, keys ...string) keyrotators.KeyRotator {
	return &roundRobinKeyRotator{
		keys:   keys,
		logger: logger,
	}
}

func (r *roundRobinKeyRotator) Next() string {
	r.keyMutex.RLock()
	defer r.keyMutex.RUnlock()

	if len(r.keys) == 0 {
		r.logger.Infof("roundRobinKeyRotator: No keys available")
		return ""
	}

	if len(r.keys) == 1 {
		r.logger.Infof("roundRobinKeyRotator: Returning single key: %s", r.keys[0])
		return r.keys[0]
	}

	index := atomic.AddUint32(&r.counter, 1) % uint32(len(r.keys))
	selectedKey := r.keys[int(index)]
	r.logger.Infof("roundRobinKeyRotator: Returning key: %s (index: %d)", selectedKey, index)
	return selectedKey
}

func (r *roundRobinKeyRotator) Add(key string) error {
	r.keyMutex.Lock()
	defer r.keyMutex.Unlock()

	r.keys = append(r.keys, key)
	return nil
}

func (r *roundRobinKeyRotator) Remove(key string) error {
	r.keyMutex.Lock()
	defer r.keyMutex.Unlock()

	for i, k := range r.keys {
		if k == key {
			r.keys = append(r.keys[:i], r.keys[i+1:]...)
			break
		}
	}
	return nil
}

func (r *roundRobinKeyRotator) GetAll() []string {
	r.keyMutex.RLock()
	defer r.keyMutex.RUnlock()

	keys := make([]string, len(r.keys))
	copy(keys, r.keys)
	return keys
}
