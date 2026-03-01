package service

import "testing"

func TestNormalizeDeleteMode_KeepExplicitMode(t *testing.T) {
	if got := normalizeDeleteMode("reassign"); got != "reassign" {
		t.Fatalf("expected reassign, got %q", got)
	}
}
