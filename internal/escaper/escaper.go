package escaper

import (
	"regexp"
	"strings"

	"github.com/gkwa/jestingjaguar/internal/logger"
)

// Escaper defines an interface for escaping template delimiters
type Escaper interface {
	Escape(content string) (string, int)
}

// TemplateEscaper is the default implementation of Escaper
type TemplateEscaper struct{}

// Escape escapes golang template delimiters
func (e *TemplateEscaper) Escape(content string) (string, int) {
	// Check for already escaped delimiters to prevent double escaping
	alreadyEscapedPattern := regexp.MustCompile(`\{\{"\{\{"}}.*?\{\{"\}\}"}}`)
	alreadyEscaped := alreadyEscapedPattern.FindAllString(content, -1)

	// Temporarily replace already escaped patterns
	replacements := make(map[string]string)
	for i, match := range alreadyEscaped {
		placeholder := "___ESCAPED_PLACEHOLDER_" + string(rune(i)) + "___"
		replacements[placeholder] = match
		content = strings.Replace(content, match, placeholder, 1)
		logger.Trace("Temporarily replaced already escaped pattern: %s", match)
	}

	// Pattern to match unescaped template delimiters
	pattern := regexp.MustCompile(`\{\{([^}]*)\}\}`)

	count := 0
	// Replace unescaped delimiters
	result := pattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the content between {{ and }}
		inner := match[2 : len(match)-2]
		count++
		logger.Trace("Escaping: %s", match)
		// Format: {{"{{"}} inner {{"}}"}}
		return `{{"{{"}}` + inner + `{{"}}"}}`
	})

	// Put back the already escaped patterns
	for placeholder, original := range replacements {
		result = strings.Replace(result, placeholder, original, 1)
		logger.Trace("Restored placeholder: %s", placeholder)
	}

	return result, count
}
