package meta

import (
	"context"

	"github.com/google/uuid"
)

type UserData struct {
	ID        uuid.UUID `json:"sub,omitempty"`
	SessionID uuid.UUID `json:"session_id,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Role      string    `json:"role,omitempty"`
}

func User(ctx context.Context) *UserData {
	if ctx == nil {
		return nil
	}

	userData, ok := ctx.Value(UserCtxKey).(UserData)
	if !ok {
		return nil
	}

	return &userData
}
