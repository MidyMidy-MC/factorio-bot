package rconcontroller

import (
	"encoding/json"
	"testing"
)

func TestEvent(t *testing.T) {
	b := []byte(`{"type":"empty"}`)

	e := &Event{}
	if err := json.Unmarshal(b, &e.Munch); err != nil {
		t.Error(err)
	}

	if e.EventType() != EventEmpty {
		t.Error("bad event!")
	}
}
