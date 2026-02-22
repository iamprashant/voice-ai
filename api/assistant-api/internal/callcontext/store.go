// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_callcontext

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
)

// Store provides operations to save and retrieve call contexts from Postgres.
//
// Call contexts are session-scoped records that live for the entire duration of a call.
// Telephony providers (Twilio, Vonage, Exotel, etc.) send event/status callbacks
// asynchronously — these can arrive at any time, including after the media stream
// has disconnected and the context has been marked "completed". Therefore, the
// row is never deleted during the call lifecycle; it is only transitioned through
// statuses: pending/queued → claimed → completed/failed.
type Store interface {
	// Save stores a call context with a generated contextId (UUID).
	// Returns the generated contextId.
	Save(ctx context.Context, cc *CallContext) (string, error)

	// Get retrieves a call context by contextId regardless of its current status
	// (pending, queued, claimed, completed, or failed). This is intentional:
	// event/status callbacks from upstream telephony providers are asynchronous
	// and may arrive after the call has already ended (status="completed").
	// The row must remain readable for the full lifetime of the context.
	Get(ctx context.Context, contextID string) (*CallContext, error)

	// Claim atomically transitions a call context from "pending" or "queued"
	// to "claimed". Inbound contexts start as "pending"; outbound contexts
	// start as "queued" (set by the outbound call creator). Only one concurrent
	// media connection can win the claim — subsequent callers get an error
	// because the row is no longer in a claimable status.
	// Returns the CallContext or an error if not found or already claimed.
	Claim(ctx context.Context, contextID string) (*CallContext, error)

	// Delete removes a call context row. This is only intended for cleanup
	// (e.g. TTL-based garbage collection), NOT during active call flows,
	// because async event callbacks may still reference the contextId.
	Delete(ctx context.Context, contextID string) error

	// Complete marks a call context as completed. Called when the call/session
	// ends. The row remains in the database so that late-arriving async event
	// callbacks from the telephony provider can still resolve the context.
	Complete(ctx context.Context, contextID string) error

	// UpdateField sets a single column on an existing call context.
	// Used to patch the channel UUID after the telephony provider returns it.
	UpdateField(ctx context.Context, contextID, field, value string) error
}

type postgresStore struct {
	postgres connectors.PostgresConnector
	logger   commons.Logger
}

// NewStore creates a new call context store backed by Postgres.
func NewStore(postgres connectors.PostgresConnector, logger commons.Logger) Store {
	return &postgresStore{
		postgres: postgres,
		logger:   logger,
	}
}

// Save stores a call context in Postgres with a generated UUID as the contextId.
func (s *postgresStore) Save(ctx context.Context, cc *CallContext) (string, error) {
	if cc.ContextID == "" {
		cc.ContextID = uuid.New().String()
	}
	if cc.Status == "" {
		cc.Status = StatusPending
	}

	db := s.postgres.DB(ctx)
	if err := db.Create(cc).Error; err != nil {
		return "", fmt.Errorf("failed to save call context %s: %w", cc.ContextID, err)
	}

	s.logger.Infof("saved call context: contextId=%s, assistant=%d, conversation=%d, direction=%s",
		cc.ContextID, cc.AssistantID, cc.ConversationID, cc.Direction)

	return cc.ContextID, nil
}

// Get retrieves a call context by contextId regardless of its status.
// Used by event/status callbacks which need the context throughout the call.
// This deliberately reads any status (pending, queued, claimed, completed, failed)
// because upstream telephony providers fire event webhooks asynchronously — a
// "completed" callback from Twilio can arrive well after the media stream ends.
func (s *postgresStore) Get(ctx context.Context, contextID string) (*CallContext, error) {
	db := s.postgres.DB(ctx)
	var cc CallContext
	if err := db.Where("context_id = ?", contextID).First(&cc).Error; err != nil {
		return nil, fmt.Errorf("call context not found: %s: %w", contextID, err)
	}

	s.logger.Debugf("resolved call context: contextId=%s, assistant=%d, conversation=%d, status=%s",
		cc.ContextID, cc.AssistantID, cc.ConversationID, cc.Status)

	return &cc, nil
}

// Claim atomically transitions a call context from "pending" or "queued" to "claimed"
// using an atomic UPDATE ... WHERE status IN ('pending','queued'). Only one concurrent
// caller can win. The context remains in the database so event callbacks can still read it.
// Both "pending" (inbound) and "queued" (outbound) are valid pre-claim states.
func (s *postgresStore) Claim(ctx context.Context, contextID string) (*CallContext, error) {
	db := s.postgres.DB(ctx)

	// Atomic update: only succeeds if status is still "pending" or "queued"
	result := db.Model(&CallContext{}).
		Where("context_id = ? AND status IN ?", contextID, []string{StatusPending, StatusQueued}).
		Updates(map[string]interface{}{
			"status":       StatusClaimed,
			"updated_date": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("failed to claim call context %s: %w", contextID, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("call context %s not found or already claimed", contextID)
	}

	// Fetch the full row after claiming
	var cc CallContext
	if err := db.Where("context_id = ?", contextID).First(&cc).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch claimed call context %s: %w", contextID, err)
	}

	s.logger.Debugf("claimed call context: contextId=%s, assistant=%d, conversation=%d",
		cc.ContextID, cc.AssistantID, cc.ConversationID)

	return &cc, nil
}

// Delete removes a call context from Postgres.
func (s *postgresStore) Delete(ctx context.Context, contextID string) error {
	db := s.postgres.DB(ctx)
	if err := db.Where("context_id = ?", contextID).Delete(&CallContext{}).Error; err != nil {
		return fmt.Errorf("failed to delete call context %s: %w", contextID, err)
	}

	s.logger.Debugf("deleted call context: contextId=%s", contextID)
	return nil
}

// Complete marks a call context as completed. Called when the call/session ends.
func (s *postgresStore) Complete(ctx context.Context, contextID string) error {
	db := s.postgres.DB(ctx)
	result := db.Model(&CallContext{}).
		Where("context_id = ?", contextID).
		Updates(map[string]interface{}{
			"status":       StatusCompleted,
			"updated_date": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to complete call context %s: %w", contextID, result.Error)
	}

	s.logger.Debugf("completed call context: contextId=%s", contextID)
	return nil
}

// UpdateField sets a single column on an existing call context row.
func (s *postgresStore) UpdateField(ctx context.Context, contextID, field, value string) error {
	db := s.postgres.DB(ctx)

	// Allowlist of updatable fields to prevent SQL injection
	allowed := map[string]bool{
		"channel_uuid": true,
		"status":       true,
		"provider":     true,
	}
	if !allowed[field] {
		return fmt.Errorf("field %q is not updatable on call context", field)
	}

	result := db.Model(&CallContext{}).
		Where("context_id = ?", contextID).
		Update(field, value)

	if result.Error != nil {
		return fmt.Errorf("failed to update field %s on call context %s: %w", field, contextID, result.Error)
	}

	s.logger.Debugf("updated call context field: contextId=%s, %s=%s", contextID, field, value)
	return nil
}
