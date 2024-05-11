package main

import (
	"encoding/json"
	"testing"
)

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

func assertFieldHasValue(t *testing.T, jsonStr []byte, fieldName string, value interface{}) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonStr, &data); err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	fieldValue, ok := data[fieldName]
	if !ok {
		t.Errorf("Expected field '%s' not found", fieldName)
	}

	if fieldValue != value {
		t.Errorf("Field '%s' value mismatch. Expected: %v, Got: %v", fieldName, value, fieldValue)
	}
}
