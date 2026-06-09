package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func TestAuthenticate_WithHeaderOverwrite(t *testing.T) {
	tests := []struct {
		name string
		headerValue string
		expectedStatus int
		expectedHeader string
	}{
		{
			name: "Valid clean username",
			headerValue: "moynul_islam",
			expectedStatus: http.StatusOK,
			expectedHeader: "moynul_islam",
		},
		{
			name: "Username with surrounding spaces",
			headerValue: "   moynul_islam   ",
			expectedStatus: http.StatusOK,
			expectedHeader: "moynul_islam",
		},
		{
			name: "Missing username header entirely",
			headerValue: "",
			expectedStatus: http.StatusUnauthorized,
			expectedHeader: "",
		},
		{
			name: "Username header containing only spaces",
			headerValue: "     ",
			expectedStatus: http.StatusUnauthorized,
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/wishlist", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tt.headerValue != "" {
				req.Header.Set("Username", tt.headerValue)
			}

			w := httptest.NewRecorder()
			ctx := context.NewContext()
			ctx.Reset(w, req)

			Authenticate(ctx)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
			if tt.expectedStatus == http.StatusOK {
				actualHeader := ctx.Input.Header("Username")
				if actualHeader != tt.expectedHeader {
					t.Errorf("Expected request header 'Username' to be %q, got %q", tt.expectedHeader, actualHeader)
				}
			}
		})
	}
}
