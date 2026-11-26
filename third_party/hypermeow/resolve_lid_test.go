// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package whatsmeow

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	waBinary "github.com/polymorfa/hypermeow/binary"
	"github.com/polymorfa/hypermeow/store"
	"github.com/polymorfa/hypermeow/types"
	waLog "github.com/polymorfa/hypermeow/util/log"
)

func TestResolveLIDWithoutProfileMapping(t *testing.T) {
	for _, tc := range []struct {
		name, query, contactType, server        string
		profileLID, transportFailure, wantError bool
	}{
		{name: "registered contact without profile LID", query: "+15550001111", contactType: "in", server: types.HiddenUserServer},
		{name: "different query cannot select recipient", query: "+15550002222", contactType: "in", server: types.HiddenUserServer, wantError: true},
		{name: "unregistered contact", query: "+15550001111", contactType: "out", server: types.HiddenUserServer, wantError: true},
		{name: "phone JID is not LID", query: "+15550001111", contactType: "in", server: types.DefaultUserServer, wantError: true},
		{name: "contact lookup failure", transportFailure: true, wantError: true},
		{name: "profile LID needs no fallback", profileLID: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pn := types.NewADJID("15550001111", types.WhatsAppDomain, 7)
			lid := types.NewJID("100000011111111", types.HiddenUserServer)
			noop := &store.NoopStore{}
			client := NewClient(&store.Device{LIDs: noop, Contacts: noop}, waLog.Noop)
			calls := 0
			transportErr := errors.New("contact transport unavailable")
			client.testSendHook = func(raw []byte) error {
				unpacked, err := waBinary.Unpack(raw)
				if err != nil {
					return err
				}
				node, err := waBinary.Unmarshal(unpacked)
				if err != nil {
					return err
				}
				if node.Tag != "iq" || node.AttrGetter().OptionalString("xmlns") != "usync" {
					t.Fatal("unexpected non-query protocol operation")
				}
				calls++
				var user waBinary.Node
				switch calls {
				case 1:
					user = waBinary.Node{Tag: "user", Attrs: waBinary.Attrs{"jid": pn.ToNonAD()}}
					if tc.profileLID {
						user.Content = []waBinary.Node{{Tag: "lid", Attrs: waBinary.Attrs{"val": lid}}}
					}
				case 2:
					query := node.GetChildByTag("usync", "query", "contact")
					if query.AttrGetter().OptionalString("addressing_mode") != "lid" {
						t.Fatal("contact lookup must request LID addressing")
					}
					requested, _ := node.GetChildByTag("usync", "list", "user", "contact").Content.([]byte)
					if string(requested) != "+15550001111@c.us" {
						t.Fatalf("wrong contact query: %q", requested)
					}
					if tc.transportFailure {
						return transportErr
					}
					user = waBinary.Node{Tag: "user", Attrs: waBinary.Attrs{"jid": types.NewJID(lid.User, tc.server)}, Content: []waBinary.Node{{Tag: "contact", Attrs: waBinary.Attrs{"type": tc.contactType}, Content: []byte(tc.query + "@c.us")}}}
				default:
					t.Fatal("unexpected extra lookup")
				}
				response := &waBinary.Node{Tag: "iq", Attrs: waBinary.Attrs{"id": node.AttrGetter().String("id"), "type": "result"}, Content: []waBinary.Node{{Tag: "usync", Content: []waBinary.Node{{Tag: "list", Content: []waBinary.Node{user}}}}}}
				if !client.receiveResponse(context.Background(), response) {
					t.Fatal("missing IQ waiter")
				}
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			resolved, err := client.ResolveLID(ctx, pn)
			if tc.wantError {
				if err == nil || !resolved.IsEmpty() {
					t.Fatalf("unsafe resolution: %v, %v", resolved, err)
				}
				if tc.contactType == "out" && !errors.Is(err, ErrRecipientNotOnWhatsApp) {
					t.Fatalf("unregistered recipient error lost: %v", err)
				}
				if tc.transportFailure && !errors.Is(err, transportErr) {
					t.Fatalf("transport error lost: %v", err)
				}
			} else {
				want := lid
				want.Device = pn.Device
				if err != nil || resolved != want {
					t.Fatalf("resolution = %v, %v; want %v", resolved, err, want)
				}
			}
			wantCalls := 2
			if tc.profileLID {
				wantCalls = 1
			}
			if calls != wantCalls {
				t.Fatalf("queries = %d, want %d", calls, wantCalls)
			}
			if err != nil && strings.Contains(err.Error(), "context deadline") {
				t.Fatal("test timed out instead of resolving the response")
			}
		})
	}
}
