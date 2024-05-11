package main

import (
	"encoding/json"
	"testing"
)

func TestToGELFTime(t *testing.T) {
	// Convert an owncloud time to GELF time format
	gelfTime, err := toGELFTime("2023-10-13 14:13:14.581000")
	if err != nil {
		t.Fatalf("Error converting time to GELF format: %v", err)
	}

	if gelfTime != 1697206394.5809999 {
		t.Errorf("Expected GELF time mismatch. Expected: %v, Got: %v", 1697206394.5809999, gelfTime)
	}

}

func TestConvertToGelf_AddsRequiredFields(t *testing.T) {
	// Convert a owncloud JSON log line to GELF format
	logLine := []byte(`{"time":"2023-10-13 14:13:14.581000","message":"hello world"}`)
	if err := convertToGELF(&logLine); err != nil {
		t.Fatalf("Error converting to GELF: %v", err)
	}

	assertFieldHasValue(t, logLine, "version", "1.1")
	assertFieldHasValue(t, logLine, "host", "localhost")
	assertFieldHasValue(t, logLine, "timestamp", 1697206394.5809999)
}

func TestConvertToGelf_PassesAllOtherFields(t *testing.T) {
	// Convert a owncloud JSON log line to GELF format
	logLine := []byte(`{"message":"hello world","time":"2023-10-13 14:13:14.581000","level":"info","user":"admin", "status":1337}`)
	if err := convertToGELF(&logLine); err != nil {
		t.Fatalf("Error converting to GELF: %v", err)
	}

	assertFieldHasValue(t, logLine, "message", "hello world")
	assertFieldHasValue(t, logLine, "level", "info")
	assertFieldHasValue(t, logLine, "user", "admin")
	assertFieldHasValue(t, logLine, "status", 1337)
}

func assertFieldHasValue(t *testing.T, jsonStr []byte, fieldName string, value interface{}) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonStr, &data); err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	rawValue, ok := data[fieldName]
	if !ok {
		t.Errorf("Expected field '%s' not found", fieldName)
	}

	switch value.(type) {
	case int:
		rawValue = int(rawValue.(float64))
		break
	case float64:
		rawValue = float64(rawValue.(float64))
		break
	case bool:
		rawValue = bool(rawValue.(bool))
	}

	if rawValue != value {
		t.Errorf("Field '%s' value mismatch. Expected: %v, Got: %v", fieldName, value, rawValue)
	}
}
