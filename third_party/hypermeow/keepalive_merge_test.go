package whatsmeow

import (
	"context"
	"testing"

	waBinary "github.com/polymorfa/hypermeow/binary"
	waLog "github.com/polymorfa/hypermeow/util/log"
)

type keepAliveTestLogger struct {
	waLog.Logger
	count int
}

func (l *keepAliveTestLogger) Debugf(string, ...any) { l.count++ }

func TestKeepAliveLoggingPreservesRawHookAndIQDelivery(t *testing.T) {
	relay := &fakeShadowRelay{}
	client := newTestShadow(t, relay)
	log := &keepAliveTestLogger{Logger: waLog.Noop}
	client.recvLog, client.sendLog = log, log
	pingID := "fixture-ping"
	client.lastKeepAliveID.Store(&pingID)
	rawCalls := 0
	client.RawNodeHandler = func(_ context.Context, _ *waBinary.Node) (*waBinary.Node, bool) {
		rawCalls++
		return nil, false
	}
	for _, id := range []string{pingID, "other-iq"} {
		req := waBinary.Node{Tag: "iq", Attrs: waBinary.Attrs{"id": id, "type": "get"}}
		if _, err := client.sendNodeAndGetData(t.Context(), req); err != nil {
			t.Fatal(err)
		}
		result := waBinary.Node{Tag: "iq", Attrs: waBinary.Attrs{"id": id, "type": "result"}}
		for _, injected := range []bool{false, true} {
			waiter := client.waitResponse(id)
			if injected {
				if err := client.InjectNode(t.Context(), &result); err != nil {
					t.Fatal(err)
				}
			} else {
				frame, err := waBinary.Marshal(result)
				if err != nil {
					t.Fatal(err)
				}
				client.handleFrame(t.Context(), frame, nil)
			}
			select {
			case <-waiter:
			default:
				t.Fatal("IQ response did not reach waiter")
			}
		}
		if id == pingID && log.count != 0 {
			t.Fatal("keepalive appeared in normal node trace")
		}
	}
	if rawCalls != 4 || relay.sentCount() != 2 || log.count != 3 {
		t.Fatalf("raw callbacks=%d, relayed requests=%d, normal trace entries=%d", rawCalls, relay.sentCount(), log.count)
	}
}
