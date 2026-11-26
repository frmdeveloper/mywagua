// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package whatsmeow

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	waBinary "github.com/polymorfa/hypermeow/binary"
	"github.com/polymorfa/hypermeow/store"
	"github.com/polymorfa/hypermeow/types"
	waLog "github.com/polymorfa/hypermeow/util/log"
)

func TestNewsletterViewedCancellationReleasesWaiter(t *testing.T) {
	client := NewClient(&store.Device{}, waLog.Noop)
	written := make(chan string, 1)
	client.testSendHook = func(raw []byte) error {
		unpacked, err := waBinary.Unpack(raw)
		if err != nil {
			return err
		}
		node, err := waBinary.Unmarshal(unpacked)
		if err != nil {
			return err
		}
		if node.Tag != "receipt" || node.AttrGetter().String("type") != "view" {
			t.Error("not a view receipt")
		}
		written <- node.AttrGetter().String("id")
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		result <- client.NewsletterMarkViewed(ctx, types.NewJID("123", types.NewsletterServer), []types.MessageServerID{4})
	}()
	id := <-written
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation lost: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		client.receiveResponse(context.Background(), &waBinary.Node{Tag: "ack", Attrs: waBinary.Attrs{"id": id}})
		<-result
		t.Fatal("view receipt ignored cancellation")
	}
	client.responseWaitersLock.Lock()
	remaining := len(client.responseWaiters)
	client.responseWaitersLock.Unlock()
	if remaining != 0 {
		t.Fatalf("%d leaked response waiters", remaining)
	}
}

func TestNewsletterViewedAcknowledgementErrors(t *testing.T) {
	for _, tc := range []struct {
		name      string
		node      *waBinary.Node
		wantError bool
	}{
		{name: "accepted", node: &waBinary.Node{Tag: "ack"}},
		{name: "rejected", node: &waBinary.Node{Tag: "ack", Attrs: waBinary.Attrs{"error": "403"}}, wantError: true},
		{name: "disconnected", node: xmlStreamEndNode, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(&store.Device{}, waLog.Noop)
			client.testSendHook = func(raw []byte) error {
				unpacked, err := waBinary.Unpack(raw)
				if err != nil {
					return err
				}
				node, err := waBinary.Unmarshal(unpacked)
				if err != nil {
					return err
				}
				if tc.node == xmlStreamEndNode {
					client.clearResponseWaiters(tc.node)
				} else {
					reply := *tc.node
					reply.Attrs = waBinary.Attrs{"id": node.AttrGetter().String("id")}
					for k, v := range tc.node.Attrs {
						reply.Attrs[k] = v
					}
					client.receiveResponse(context.Background(), &reply)
				}
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := client.NewsletterMarkViewed(ctx, types.NewJID("123", types.NewsletterServer), []types.MessageServerID{4})
			if (err != nil) != tc.wantError {
				t.Fatalf("ack error = %v", err)
			}
		})
	}
}

func TestCancelResponseIsSafeAfterAcknowledgement(t *testing.T) {
	client := NewClient(&store.Device{}, waLog.Noop)
	for range 1000 {
		id := client.generateRequestID()
		waiter := client.waitResponse(id)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			client.receiveResponse(context.Background(), &waBinary.Node{Tag: "ack", Attrs: waBinary.Attrs{"id": id}})
		}()
		go func() { defer wg.Done(); client.cancelResponse(id, waiter) }()
		wg.Wait()
		client.cancelResponse(id, waiter)
	}
}
