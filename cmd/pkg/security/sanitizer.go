package security

import (
	"html"
	"strings"
	"unicode"
)

const (
	MaxNameLength = 50 // Hard limit to prevent memory exhaustion or log bloat
)

// SanitizeName takes an untrusted input and returns a cryptographically safe string.
func SanitizeName(raw string) string {
	// 1. Initial Trim
	raw = strings.TrimSpace(raw)

	// 2. Length Truncation (O(1) prevention of payload bombs)
	if len(raw) > MaxNameLength {
		// Cut safely by runes, not bytes, to avoid breaking multi-byte characters
		runes := []rune(raw)
		if len(runes) > MaxNameLength {
			raw = string(runes[:MaxNameLength])
		}
	}

	// 3. Zero-Trust Whitelist Evaluation (The Core Security Layer)
	var safeBuilder strings.Builder
	// Pre-allocate memory to prevent reallocation overhead during the loop
	safeBuilder.Grow(len(raw)) 

	for _, r := range raw {
		// Allow standard letters (includes international characters like é, ñ, etc.)
		if unicode.IsLetter(r) {
			safeBuilder.WriteRune(r)
			continue
		}

		// Allow safe separators (spaces, hyphens, apostrophes for names like O'Connor or Anne-Marie)
		if r == ' ' || r == '-' || r == '\'' {
			safeBuilder.WriteRune(r)
			continue
		}

		// If a character makes it here, it is considered hostile (numbers, symbols, <script> tags, emojis, control chars).
		// It is silently dropped.
	}

	safeString := safeBuilder.String()

	// 4. Double Cleanse (Whitespace collapse)
	// If a malicious user entered "A     B", collapse it to "A B"
	safeString = strings.Join(strings.Fields(safeString), " ")

	// 5. Final Output Encoding
	// Even though we stripped symbols, passing it through standard HTML escaping 
	// ensures absolute safety before it hits the Next.js frontend.
	return html.EscapeString(safeString)
}

// SanitizeString is a slightly more permissive generic sanitizer for things like text bodies,
// allowing basic punctuation but stripping control characters and HTML.
func SanitizeString(raw string) string {
	raw = strings.TrimSpace(raw)
	
	var safeBuilder strings.Builder
	for _, r := range raw {
		// Drop invisible control characters (e.g., \n, \r, null bytes)
		if unicode.IsControl(r) {
			continue
		}
		safeBuilder.WriteRune(r)
	}

	return html.EscapeString(safeBuilder.String())
}