package armadillo_test

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	legacy "github.com/polymorfa/hypermeow/proto/waWebLabyrinthWaWasm"
	current "github.com/polymorfa/hypermeow/proto/waWebLabyrinthWaWasmV1054751086"
)

func TestLabyrinthEpochIDFormats(t *testing.T) {
	for _, tc := range []struct {
		msg  proto.Message
		kind protoreflect.Kind
	}{
		{&legacy.DeriveMessageKeyInput{EpochRootKey: []byte{1}, EpochAnonID: []byte{2}, ThreadID: proto.String("fixture")}, protoreflect.BytesKind},
		{&current.DeriveMessageKeyInput{EpochRootKey: []byte{1}, EpochAnonID: proto.Int32(2), ThreadID: proto.String("fixture")}, protoreflect.Int32Kind},
	} {
		field := tc.msg.ProtoReflect().Descriptor().Fields().ByNumber(2)
		if field == nil || field.Kind() != tc.kind {
			t.Fatalf("unexpected field 2 descriptor in %T", tc.msg)
		}
		wire, err := proto.Marshal(tc.msg)
		if err != nil {
			t.Fatal(err)
		}
		decoded := tc.msg.ProtoReflect().New().Interface()
		if err = proto.Unmarshal(wire, decoded); err != nil || !proto.Equal(tc.msg, decoded) {
			t.Fatalf("failed to round-trip %T: %v", tc.msg, err)
		}
	}
}
