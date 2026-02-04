package auth

import (
	"os"
	"testing"
)

func TestGetCallbackBaseURL(t *testing.T) {
	tests := []struct {
		name               string
		callbackBaseURL    string
		railwayPublicDomain string
		backendURL         string
		expected           string
	}{
		{
			name:            "Uses CALLBACK_BASE_URL when set",
			callbackBaseURL: "https://pr-environment.railway.app",
			railwayPublicDomain: "ignored.railway.app",
			backendURL:      "https://production.example.com",
			expected:        "https://pr-environment.railway.app",
		},
		{
			name:                "Uses RAILWAY_PUBLIC_DOMAIN when CALLBACK_BASE_URL not set",
			callbackBaseURL:     "",
			railwayPublicDomain: "my-app-pr-30.up.railway.app",
			backendURL:          "https://production.example.com",
			expected:            "https://my-app-pr-30.up.railway.app",
		},
		{
			name:                "Falls back to BACKEND_URL when neither is set",
			callbackBaseURL:     "",
			railwayPublicDomain: "",
			backendURL:          "https://production.example.com",
			expected:            "https://production.example.com",
		},
		{
			name:                "Handles empty BACKEND_URL",
			callbackBaseURL:     "",
			railwayPublicDomain: "",
			backendURL:          "",
			expected:            "",
		},
		{
			name:               "CALLBACK_BASE_URL takes priority over RAILWAY_PUBLIC_DOMAIN",
			callbackBaseURL:    "https://explicit-override.com",
			railwayPublicDomain: "railway-domain.com",
			backendURL:         "https://backend.com",
			expected:           "https://explicit-override.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalCallbackURL := os.Getenv("CALLBACK_BASE_URL")
			originalRailwayDomain := os.Getenv("RAILWAY_PUBLIC_DOMAIN")
			originalBackendURL := BackendURL

			// Set test values
			if tt.callbackBaseURL != "" {
				os.Setenv("CALLBACK_BASE_URL", tt.callbackBaseURL)
			} else {
				os.Unsetenv("CALLBACK_BASE_URL")
			}

			if tt.railwayPublicDomain != "" {
				os.Setenv("RAILWAY_PUBLIC_DOMAIN", tt.railwayPublicDomain)
			} else {
				os.Unsetenv("RAILWAY_PUBLIC_DOMAIN")
			}

			BackendURL = tt.backendURL

			// Test
			result := getCallbackBaseURL()
			if result != tt.expected {
				t.Errorf("getCallbackBaseURL() = %v, want %v", result, tt.expected)
			}

			// Restore original values
			if originalCallbackURL != "" {
				os.Setenv("CALLBACK_BASE_URL", originalCallbackURL)
			} else {
				os.Unsetenv("CALLBACK_BASE_URL")
			}

			if originalRailwayDomain != "" {
				os.Setenv("RAILWAY_PUBLIC_DOMAIN", originalRailwayDomain)
			} else {
				os.Unsetenv("RAILWAY_PUBLIC_DOMAIN")
			}

			BackendURL = originalBackendURL
		})
	}
}

func TestOAuthConfigUsesCorrectCallbackURL(t *testing.T) {
	// This test verifies that the OAuth configs are using getCallbackBaseURL()
	// Note: Since conf and googleConf are package-level variables initialized at load time,
	// this test just verifies they exist and have the expected structure

	if conf == nil {
		t.Error("GitHub OAuth config (conf) should not be nil")
	}

	if googleConf == nil {
		t.Error("Google OAuth config (googleConf) should not be nil")
	}

	// Verify that the redirect URLs follow the expected pattern
	if conf != nil && conf.RedirectURL == "" {
		t.Error("GitHub OAuth RedirectURL should not be empty")
	}

	if googleConf != nil && googleConf.RedirectURL == "" {
		t.Error("Google OAuth RedirectURL should not be empty")
	}

	// Verify the URLs end with the expected callback paths
	if conf != nil {
		expectedSuffix := "/auth/github-callback"
		if len(conf.RedirectURL) < len(expectedSuffix) ||
			conf.RedirectURL[len(conf.RedirectURL)-len(expectedSuffix):] != expectedSuffix {
			t.Errorf("GitHub OAuth RedirectURL should end with %s, got: %s", expectedSuffix, conf.RedirectURL)
		}
	}

	if googleConf != nil {
		expectedSuffix := "/auth/google-callback"
		if len(googleConf.RedirectURL) < len(expectedSuffix) ||
			googleConf.RedirectURL[len(googleConf.RedirectURL)-len(expectedSuffix):] != expectedSuffix {
			t.Errorf("Google OAuth RedirectURL should end with %s, got: %s", expectedSuffix, googleConf.RedirectURL)
		}
	}
}
