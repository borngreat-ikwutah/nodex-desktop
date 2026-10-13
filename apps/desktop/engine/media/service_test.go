package media

import (
	"encoding/json"
	"image"
	"image/color"
	"testing"
)

func TestSignalRoundTrip(t *testing.T) {
	original := Signal{SessionID: "session-1", SDP: "v=0", Candidate: "candidate:1", SDPMLine: uint16Ptr(0)}
	payload, err := EncodeSignal(original)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeSignal(json.RawMessage(payload))
	if err != nil {
		t.Fatal(err)
	}
	if decoded.SessionID != original.SessionID || decoded.SDP != original.SDP || decoded.Candidate != original.Candidate || decoded.SDPMLine == nil || *decoded.SDPMLine != 0 {
		t.Fatalf("signal changed during round trip: %#v", decoded)
	}
}

func TestDeltaEncoderEmitsKeyframeThenOnlyDirtyTiles(t *testing.T) {
	encoder := NewDeltaEncoder(2)
	first := image.NewRGBA(image.Rect(0, 0, 4, 4))
	first.Set(0, 0, color.RGBA{R: 255, A: 255})
	keyframe, err := encoder.Encode(first)
	if err != nil {
		t.Fatal(err)
	}
	if string(keyframe[:4]) != "NDF1" || keyframe[5] != 1 {
		t.Fatalf("expected keyframe packet")
	}
	second := image.NewRGBA(first.Bounds())
	second.Set(0, 0, color.RGBA{R: 255, A: 255})
	second.Set(3, 3, color.RGBA{B: 255, A: 255})
	delta, err := encoder.Encode(second)
	if err != nil {
		t.Fatal(err)
	}
	if delta[5] != 0 {
		t.Fatalf("expected delta packet")
	}
	if len(delta) >= len(keyframe) {
		t.Fatalf("expected dirty tile packet to be smaller")
	}
}

func uint16Ptr(value uint16) *uint16 { return &value }
