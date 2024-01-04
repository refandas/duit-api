package helper

import "net/http"

// SetupSecurityHeaders configures and set HTTP security headers on the provided
// http.ResponseWriter based on the configuration proposal outlined by OWASP.
// https://owasp.org/www-project-secure-headers/index.html#configuration-proposal
func SetupSecurityHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Cache-Control", "no-store, max-age=0")
	writer.Header().Set("Clear-Site-Data", "\"cache\",\"cookies\",\"storage\"")
	writer.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	writer.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.Header().Set("X-Frame-Options", "DENY")
	writer.Header().Set("X-Permitted-Cross-Domain-Policies", "none")

	// Provide additional security when responses are rendered as HTML
	writer.Header().Set("Content-Security-Policy", "default-src 'none'")
	writer.Header().Set("Permissions-Policy", "accelerometer=(),ambient-light-sensor=(),autoplay=(),battery=(),camera=(),display-capture=(),document-domain=(),encrypted-media=(),fullscreen=(),gamepad=(),geolocation=(),gyroscope=(),layout-animations=(self),legacy-image-formats=(self),magnetometer=(),microphone=(),midi=(),oversized-images=(self),payment=(),picture-in-picture=(),publickey-credentials-get=(),speaker-selection=(),sync-xhr=(self),unoptimized-images=(self),unsized-media=(self),usb=(),screen-wake-lock=(),web-share=(),xr-spatial-tracking=()")
	writer.Header().Set("Referrer-Policy", "no-referrer")
}
