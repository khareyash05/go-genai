package genai

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TestAssertRequestBodyFormatting_ValidComparison_001 tests the behavior of `assertRequest` when comparing request body formatting and headers.
func TestAssertRequestBodyFormatting_ValidComparison_001(t *testing.T) {
	// Arrange
	rac := &replayAPIClient{t: t}
	sdkRequestBody := `{
         "key": "value",
         "nested": {
             "key2": "value2"
         }
     }`
	sdkRequest := &http.Request{
		Body: io.NopCloser(strings.NewReader(sdkRequestBody)),
		Header: http.Header{
			"Content-Type":      []string{"application/json"},
			"User-Agent":        []string{"test-agent"},
			"X-Goog-Api-Client": []string{"test-client"},
		},
		Method: "POST",
		URL:    &url.URL{Path: "/test-path"},
	}
	want := &replayRequest{
		Method: "POST",
		Url:    "/test-path",
		Headers: map[string]string{
			"content-type":      "application/json",
			"user-agent":        "test-agent",
			"x-goog-api-client": "test-client",
		},
		BodySegments: []map[string]any{
			{
				"key": "value",
				"nested": map[string]any{
					"key2": "value2",
				},
			},
		},
	}

	// Act
	rac.assertRequest(sdkRequest, want)

	// Assert
	// No explicit assertions here as the test will fail if the logic in assertRequest does not handle formatting correctly.
}
