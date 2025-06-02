package genai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAssertRequest_StringNormalization_MethodAndBody_Cover_582 tests the string normalization
// logic (lines 247-253) in assertRequest for Method and BodySegment fields.
// It verifies that newlines are removed, multiple spaces are condensed, strings are trimmed,
// and comparisons are case-insensitive.
func TestAssertRequest_StringNormalization_MethodAndBody_Cover_582(t *testing.T) {
	t.Helper()

	rac := newReplayAPIClient(t) // This sets up rac.t and rac.server
	defer rac.server.Close()

	// 1. Define the SDK request (represents the "got" side after initial processing by assertRequest)
	// Method should be a clean, valid HTTP method.
	// Body strings should be in their "clean" or "normalized" form.
	sdkMethod := http.MethodPost // "POST"
	sdkBodyMap := map[string]any{
		"message":       "hello world",   // Target for normalization: mixed case, newlines, spaces
		"detail_field":  "value1 value2", // Target for normalization: multiple internal spaces
		"empty_ish":     "",              // Target for normalization: spaces, newline, tab to empty
		"already_clean": "clean",         // Target for EqualFold with different case
		"tab_test":      "a b",           // Target for tab normalization
	}
	sdkBodyJSON, err := json.Marshal(sdkBodyMap)
	if err != nil {
		t.Fatalf("Failed to marshal SDK body: %v", err)
	}

	sdkRequest := httptest.NewRequest(
		sdkMethod,
		rac.GetBaseURL()+"/test/normalization_path", // URL for the request
		bytes.NewReader(sdkBodyJSON),
	)
	// Set headers required by the custom header comparer in assertRequest to ensure it passes,
	// allowing focus on the general string comparer for Method and Body.
	sdkRequest.Header.Set("User-Agent", "test-agent/coverage/1.0")
	sdkRequest.Header.Set("X-Goog-Api-Client", "test-client/coverage/1.0")
	sdkRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

	// 2. Define the "want" replayRequest (represents the "want" side, e.g., from a replay file)
	// These strings are "dirty" and are expected to be normalized by the cmp.Comparer (lines 240-253)
	// to match the "clean" values derived from sdkRequest.
	wantRequest := &replayRequest{
		Method: "  pOSt\r\n  ",             // Tests trimming, newline removal, case-insensitivity
		Url:    "/test/normalization_path", // Must match for URL suffix comparison to pass
		Headers: map[string]string{ // Must match for the specific header comparer to pass
			"user-agent":        "test-agent/coverage/1.0",
			"x-goog-api-client": "test-client/coverage/1.0",
			"content-type":      "application/json; charset=utf-8",
		},
		BodySegments: []map[string]any{
			{
				"message":       "  HELLO\r\n   WORLD  ", // Tests newlines, leading/trailing/multiple spaces, case
				"detail_field":  "value1   value2",       // Tests multiple internal spaces
				"empty_ish":     "  \r\n\t  ",            // Tests normalization to empty (includes tab)
				"already_clean": "CLEAN",                 // Tests EqualFold with different case
				"tab_test":      "a\tb",                  // Tests tab to single space
			},
		},
	}

	// 3. Call assertRequest.
	// The method uses rac.t.Errorf internally. If no error is reported by t.Errorf,
	// it means all comparisons, including the string normalization, passed as expected.
	// This implicitly verifies that lines 247-253 are working correctly.
	rac.assertRequest(sdkRequest, wantRequest)
}
