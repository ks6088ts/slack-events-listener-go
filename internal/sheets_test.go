package internal

import (
	"testing"
)

func TestNewSheetClient(t *testing.T) {
	_, err := NewSheetClient("invalid", "invalid")
	if err == nil {
		t.Error("calling NewSheetClient() with invalid file should fail")
	}
}
