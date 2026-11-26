// Copyright (c) 2025 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package upgrades

import (
	"context"
	"embed"

	"go.mau.fi/util/dbutil"
)

//go:embed *.sql
var upgrades embed.FS

var Table = dbutil.BuildUpgradeTable().WithFS(upgrades).
	WithRaw(18, 19, 8, "Reconcile upstream and HyperMeow chat settings", dbutil.TxnModeOn, reconcileUpstreamSchema).Finish()

func reconcileUpstreamSchema(ctx context.Context, db *dbutil.Database) error {
	// Upstream v16 already has this column, while HyperMeow v16 introduced
	// migration indexes instead. Detect the schema rather than trusting the
	// overlapping version number. Preserve any existing root-secret values.
	exists, err := db.ColumnExists(ctx, "whatsmeow_chat_settings", "wasa_root_secret_id")
	if err != nil {
		return err
	}
	if !exists {
		if _, err = db.Exec(ctx, "ALTER TABLE whatsmeow_chat_settings ADD COLUMN wasa_root_secret_id TEXT NOT NULL DEFAULT ''"); err != nil {
			return err
		}
	}
	// Importing upstream v16 skips the fork's v15->v16 index migration.
	_, err = db.Exec(ctx, "CREATE INDEX IF NOT EXISTS whatsmeow_sender_keys_sender_idx ON whatsmeow_sender_keys (our_jid, sender_id)")
	return err
}
