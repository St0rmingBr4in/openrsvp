package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSanitizeStrictEntityEncodedTags covers the defect where SanitizeStrict
// unescaped the output of bluemonday. The unescape ran after the sanitizer, so
// an entity-encoded tag became a live tag in storage.
func TestSanitizeStrictEntityEncodedTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"named entities", "&lt;script&gt;alert(1)&lt;/script&gt;"},
		{"decimal entities", "&#60;script&#62;alert(1)&#60;/script&#62;"},
		{"hex entities", "&#x3C;script&#x3E;alert(1)&#x3C;/script&#x3E;"},
		{"uppercase hex entities", "&#X3c;script&#X3e;alert(1)&#X3c;/script&#X3e;"},
		{"double encoded", "&amp;lt;script&amp;gt;alert(1)&amp;lt;/script&amp;gt;"},
		{"triple encoded", "&amp;amp;lt;script&amp;amp;gt;alert(1)"},
		{"image onerror", "&lt;img src=x onerror=alert(1)&gt;"},
		{"split tag", "&lt;scr&lt;script&gt;ipt&gt;alert(1)"},
		{"svg onload", "&lt;svg/onload=alert(1)&gt;"},
		{"iframe src", "&lt;iframe src=javascript:alert(1)&gt;"},
		{"comment wrapper", "&lt;!--&gt;&lt;script&gt;alert(1)&lt;/script&gt;--&gt;"},
		{"mixed encoding", "&#60;scr&amp;#105;pt&#62;alert(1)"},
		{"legacy no semicolon", "&ltscript&gt alert(1)"},
		{"uppercase named", "&LT;script&GT;alert(1)"},
		{"entity assembly across passes", "&amp;l&amp;#116;;script&amp;gt;alert(1)"},
		{"exotic undecoded entities", "&nGt;&nLt;script&nGt;alert(1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeStrict(tt.input)

			// The stored value must not contain a tag that a browser can
			// parse. A sanitizer that returns its input re-escaped is
			// acceptable. A sanitizer that returns a live tag is not.
			assert.NotContains(t, got, "<script", "stored a live script tag: %q", got)
			assert.NotContains(t, got, "<img", "stored a live img tag: %q", got)
			assert.NotContains(t, got, "<svg", "stored a live svg tag: %q", got)
			assert.NotContains(t, got, "<iframe", "stored a live iframe tag: %q", got)

			// Stronger property: the result must be a fixed point of the
			// sanitizer. If it is not, a second pass would change it, which
			// means a tag survived in some encoded form.
			assert.Equal(t, got, SanitizeStrict(got),
				"result is not stable under a second pass: %q", got)
		})
	}
}

// TestSanitizeStrictKeepsPlainText makes sure the fix does not regress ordinary
// text. Users write comparison operators and punctuation in event descriptions.
func TestSanitizeStrictKeepsPlainText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"plain text", "plain text"},
		{"a < b", "a < b"},
		{"5 > 3 and 2 < 4", "5 > 3 and 2 < 4"},
		{"Beer & Wine < 5pm", "Beer & Wine < 5pm"},
		{"O'Brien", "O'Brien"},
		{"Tom & Jerry", "Tom & Jerry"},
		{`say "hi"`, `say "hi"`},
		{"café — 20€", "café — 20€"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, SanitizeStrict(tt.input))
		})
	}
}

// TestSanitizeStrictNonConvergent covers the iteration cap. Deeply nested
// encoding does not reach a fixed point within the cap, so the function returns
// the escaped form. The result must be inert.
func TestSanitizeStrictNonConvergent(t *testing.T) {
	input := "&amp;amp;amp;amp;lt;script&amp;amp;amp;amp;gt;alert(1)"
	got := SanitizeStrict(input)

	assert.NotContains(t, got, "<script", "cap path stored a live tag: %q", got)
	require.False(t, strings.Contains(got, "<"),
		"cap path must return fully escaped output, got %q", got)
}

// TestSanitizeStrictIdempotent makes sure a value that is edited and saved again
// does not drift. Storage passes through this function on every write.
func TestSanitizeStrictIdempotent(t *testing.T) {
	inputs := []string{
		"Hello World",
		"a < b",
		"&lt;script&gt;alert(1)&lt;/script&gt;",
		"<h1>Hello</h1> <script>alert(1)</script>World",
		"Beer & Wine < 5pm",
	}

	for _, in := range inputs {
		once := SanitizeStrict(in)
		twice := SanitizeStrict(once)
		assert.Equal(t, once, twice, "not idempotent for %q", in)
	}
}
