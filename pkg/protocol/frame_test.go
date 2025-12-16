package protocol

import (
	"bytes"
	"testing"
)

func TestFrameReadWrite(t *testing.T) {
	payload := []byte{0x01, 0x02, 0x03}
	var buf bytes.Buffer

	if err := WriteFrame(&buf, payload); err != nil {
		t.Fatalf("WriteFrame: %v", err)
	}

	frame, err := ReadFrame(&buf)
	if err != nil {
		t.Fatalf("ReadFrame: %v", err)
	}

	if frame.Length != int32(len(payload)) {
		t.Fatalf("unexpected frame length: %d", frame.Length)
	}
	if !bytes.Equal(frame.Payload, payload) {
		t.Fatalf("payload mismatch: %v vs %v", frame.Payload, payload)
	}
}
