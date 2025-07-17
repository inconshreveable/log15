package compat

import (
	"bytes"
	"testing"
	"time"

	legacy "github.com/inconshreveable/log15"
	"github.com/inconshreveable/log15/v3"
)

func TestCompatHandler(t *testing.T) {
	t.Parallel()

	// Create a v3 handler that captures records
	var capturedRecord *log15.Record
	v3Handler := log15.FuncHandler(func(r log15.Record) error {
		capturedRecord = &r
		return nil
	})

	// Wrap it for legacy use
	legacyHandler := CompatHandler(v3Handler)

	// Create a legacy record
	legacyRecord := &legacy.Record{
		Time: time.Now(),
		Lvl:  legacy.LvlInfo,
		Msg:  "test message",
		Ctx:  []interface{}{"key", "value"},
		KeyNames: legacy.RecordKeyNames{
			Time: "t",
			Lvl:  "lvl",
			Msg:  "msg",
		},
	}

	// Log through the legacy interface
	err := legacyHandler.Log(legacyRecord)
	if err != nil {
		t.Fatalf("CompatHandler.Log failed: %v", err)
	}

	// Verify the record was converted and passed to v3 handler
	if capturedRecord == nil {
		t.Fatal("No record captured by v3 handler")
	}

	if capturedRecord.Msg != "test message" {
		t.Errorf("Message not preserved: got %q, want %q", capturedRecord.Msg, "test message")
	}

	if capturedRecord.Lvl != log15.LvlInfo {
		t.Errorf("Level not converted: got %v, want %v", capturedRecord.Lvl, log15.LvlInfo)
	}

	if len(capturedRecord.Ctx) != 2 || capturedRecord.Ctx[0] != "key" || capturedRecord.Ctx[1] != "value" {
		t.Errorf("Context not preserved: got %v, want [key value]", capturedRecord.Ctx)
	}

	if capturedRecord.KeyNames == nil {
		t.Fatal("KeyNames not converted")
	}

	if capturedRecord.KeyNames.Time != "t" || capturedRecord.KeyNames.Lvl != "lvl" || capturedRecord.KeyNames.Msg != "msg" {
		t.Errorf("KeyNames not preserved: got %+v", capturedRecord.KeyNames)
	}
}

func TestCompatHandlerWithNilKeyNames(t *testing.T) {
	t.Parallel()

	var capturedRecord *log15.Record
	v3Handler := log15.FuncHandler(func(r log15.Record) error {
		capturedRecord = &r
		return nil
	})

	legacyHandler := CompatHandler(v3Handler)

	// Legacy record without KeyNames
	legacyRecord := &legacy.Record{
		Time: time.Now(),
		Lvl:  legacy.LvlError,
		Msg:  "error message",
		Ctx:  []interface{}{"error", "test"},
	}

	err := legacyHandler.Log(legacyRecord)
	if err != nil {
		t.Fatalf("CompatHandler.Log failed: %v", err)
	}

	if capturedRecord == nil {
		t.Fatal("No record captured")
	}

	// Should handle nil KeyNames gracefully
	if capturedRecord.KeyNames == nil {
		t.Error("KeyNames should be initialized even when legacy record has none")
	}
}

func TestCompatHandlerInterfaceCompliance(t *testing.T) {
	t.Parallel()

	v3Handler := log15.FuncHandler(func(r log15.Record) error { return nil })
	legacyHandler := CompatHandler(v3Handler)

	// Verify it implements legacy.Handler interface
	var _ legacy.Handler = legacyHandler
}

func TestCompatHandlerWithRealV3Handlers(t *testing.T) {
	t.Parallel()

	// Test with StreamHandler
	var buf bytes.Buffer
	v3StreamHandler := log15.StreamHandler(&buf, log15.JsonFormat())
	legacyHandler := CompatHandler(v3StreamHandler)

	legacyRecord := &legacy.Record{
		Time: time.Now(),
		Lvl:  legacy.LvlWarn,
		Msg:  "warning message",
		Ctx:  []interface{}{"component", "test"},
	}

	err := legacyHandler.Log(legacyRecord)
	if err != nil {
		t.Fatalf("CompatHandler with StreamHandler failed: %v", err)
	}

	// Verify output was written
	if buf.Len() == 0 {
		t.Error("No output written to buffer")
	}

	output := buf.String()
	if !containsAll(output, "warning message", "component", "test") {
		t.Errorf("Output missing expected content: %s", output)
	}
}

func TestCompatLogger(t *testing.T) {
	t.Parallel()

	// Create a v3 logger with a test handler
	var capturedRecord *log15.Record
	v3Logger := log15.New()
	v3Logger.SetHandler(log15.FuncHandler(func(r log15.Record) error {
		capturedRecord = &r
		return nil
	}))

	// Wrap it for legacy use
	legacyLogger := CompatLogger(v3Logger)

	// Test basic logging methods
	legacyLogger.Debug("debug message", "key", "value")
	if capturedRecord == nil || capturedRecord.Msg != "debug message" {
		t.Errorf("Debug logging failed: %+v", capturedRecord)
	}
	if capturedRecord.Lvl != log15.LvlDebug {
		t.Errorf("Debug level not preserved: got %v, want %v", capturedRecord.Lvl, log15.LvlDebug)
	}

	legacyLogger.Info("info message")
	if capturedRecord.Msg != "info message" || capturedRecord.Lvl != log15.LvlInfo {
		t.Errorf("Info logging failed: %+v", capturedRecord)
	}

	legacyLogger.Warn("warn message")
	if capturedRecord.Msg != "warn message" || capturedRecord.Lvl != log15.LvlWarn {
		t.Errorf("Warn logging failed: %+v", capturedRecord)
	}

	legacyLogger.Error("error message")
	if capturedRecord.Msg != "error message" || capturedRecord.Lvl != log15.LvlError {
		t.Errorf("Error logging failed: %+v", capturedRecord)
	}

	legacyLogger.Crit("crit message")
	if capturedRecord.Msg != "crit message" || capturedRecord.Lvl != log15.LvlCrit {
		t.Errorf("Crit logging failed: %+v", capturedRecord)
	}
}

func TestCompatLoggerNew(t *testing.T) {
	t.Parallel()

	var capturedRecord *log15.Record
	v3Logger := log15.New()
	v3Logger.SetHandler(log15.FuncHandler(func(r log15.Record) error {
		capturedRecord = &r
		return nil
	}))

	legacyLogger := CompatLogger(v3Logger)

	// Create child logger with context
	childLogger := legacyLogger.New("component", "test", "version", "1.0")

	// Child should implement legacy.Logger interface
	var _ legacy.Logger = childLogger

	// Log from child
	childLogger.Info("child message")

	if capturedRecord == nil {
		t.Fatal("No record captured from child logger")
	}

	if capturedRecord.Msg != "child message" {
		t.Errorf("Child message not preserved: got %q", capturedRecord.Msg)
	}

	// Check context was preserved
	expectedCtx := []interface{}{"component", "test", "version", "1.0"}
	if len(capturedRecord.Ctx) != 4 {
		t.Errorf("Context length wrong: got %d, want %d", len(capturedRecord.Ctx), 4)
	}
	for i, expected := range expectedCtx {
		if i >= len(capturedRecord.Ctx) || capturedRecord.Ctx[i] != expected {
			t.Errorf("Context[%d]: got %v, want %v", i, capturedRecord.Ctx[i], expected)
		}
	}
}

func TestCompatLoggerGetSetHandler(t *testing.T) {
	t.Parallel()

	v3Logger := log15.New()
	legacyLogger := CompatLogger(v3Logger)

	// Test GetHandler returns a legacy.Handler
	handler := legacyLogger.GetHandler()
	var _ legacy.Handler = handler

	// Test SetHandler with a legacy handler
	var capturedLegacyRecord *legacy.Record
	legacyHandler := legacy.FuncHandler(func(r *legacy.Record) error {
		capturedLegacyRecord = r
		return nil
	})

	legacyLogger.SetHandler(legacyHandler)

	// Log something
	legacyLogger.Info("test message", "key", "value")

	// Verify it went through the legacy handler
	if capturedLegacyRecord == nil {
		t.Fatal("No record captured by legacy handler")
	}

	if capturedLegacyRecord.Msg != "test message" {
		t.Errorf("Message not preserved: got %q", capturedLegacyRecord.Msg)
	}

	if capturedLegacyRecord.Lvl != legacy.LvlInfo {
		t.Errorf("Level not converted: got %v", capturedLegacyRecord.Lvl)
	}
}

func TestCompatLoggerInterfaceCompliance(t *testing.T) {
	t.Parallel()

	v3Logger := log15.New()
	legacyLogger := CompatLogger(v3Logger)

	// Verify it implements legacy.Logger interface
	var _ legacy.Logger = legacyLogger
}

func TestCompatLoggerWithContext(t *testing.T) {
	t.Parallel()

	var capturedRecord *log15.Record
	v3Logger := log15.New("global", "context")
	v3Logger.SetHandler(log15.FuncHandler(func(r log15.Record) error {
		capturedRecord = &r
		return nil
	}))

	legacyLogger := CompatLogger(v3Logger)
	legacyLogger.Info("test message", "local", "context")

	if capturedRecord == nil {
		t.Fatal("No record captured")
	}

	// Should have both global and local context
	expectedCtx := []interface{}{"global", "context", "local", "context"}
	if len(capturedRecord.Ctx) != 4 {
		t.Errorf("Context length wrong: got %d, want %d", len(capturedRecord.Ctx), 4)
	}
	for i, expected := range expectedCtx {
		if i >= len(capturedRecord.Ctx) || capturedRecord.Ctx[i] != expected {
			t.Errorf("Context[%d]: got %v, want %v", i, capturedRecord.Ctx[i], expected)
		}
	}
}

// Helper function to check if string contains all substrings
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !bytes.Contains([]byte(s), []byte(sub)) {
			return false
		}
	}
	return true
}
