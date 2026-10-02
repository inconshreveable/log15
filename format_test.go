package log15

import (
	"strings"
	"testing"
)

func TestLogfmtNonStringKey(t *testing.T) {
	t.Parallel()

	l, buf := testFormatter(LogfmtFormat())
	l.Info("msg", 123, "val")

	got := buf.String()
	expected := `LOG15_ERROR="123 is not a string key"`
	if !strings.Contains(got, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, got)
	}
}

func TestLogfmtOddArgsNonStringKey(t *testing.T) {
	t.Parallel()

	l, buf := testFormatter(LogfmtFormat())
	l.Info("msg", 123)

	got := buf.String()
	expectedKey := `LOG15_ERROR="123 is not a string key"`
	expectedOdd := `LOG15_ERROR="Normalized odd number of arguments by adding nil"`
	if !strings.Contains(got, expectedKey) {
		t.Fatalf("expected output to contain %q, got %q", expectedKey, got)
	}
	if !strings.Contains(got, expectedOdd) {
		t.Fatalf("expected output to contain %q, got %q", expectedOdd, got)
	}
}
