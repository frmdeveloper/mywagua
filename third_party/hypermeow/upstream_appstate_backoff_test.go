// Copyright (c) 2026 WPPConnect Team
// Licensed under the MIT License. See LICENSE-MIT for details.

package whatsmeow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/polymorfa/hypermeow/appstate"
	"github.com/polymorfa/hypermeow/store"
)

func TestAppStateBackoffDoesNotTouchStoreAndReleasesLock(t *testing.T) {
	cli := NewClient(&store.Device{}, nil)
	name := appstate.WAPatchName("regular")
	cli.appStateBackoff[name] = time.Now().Add(time.Hour)
	// No store backend is provided: a rate-limited fetch must return before IO,
	// including when a full sync was requested, and unlock for the next request.
	for range 2 {
		if err := cli.FetchAppState(context.Background(), name, true, false); !errors.Is(err, ErrAppStateRateLimited) {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
