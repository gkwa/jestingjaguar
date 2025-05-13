package escaper

import (
	"testing"
)

func TestTemplateEscaper_Escape(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectedCount  int
	}{
		{
			name:           "Simple template",
			input:          "artifacts/{{ workflow.name }}",
			expectedOutput: `artifacts/{{"{{"}} workflow.name {{"}}"}}`,
			expectedCount:  1,
		},
		{
			name:           "Multiple templates",
			input:          "{{ first }} and {{ second }}",
			expectedOutput: `{{"{{"}} first {{"}}"}} and {{"{{"}} second {{"}}"}}`,
			expectedCount:  2,
		},
		{
			name:           "Already escaped template",
			input:          `artifacts/{{"{{"}} workflow.name {{"}}"}}`,
			expectedOutput: `artifacts/{{"{{"}} workflow.name {{"}}"}}`,
			expectedCount:  0,
		},
		{
			name:           "Mixed escaped and unescaped",
			input:          `{{ first }} and {{"{{"}} second {{"}}"}}`,
			expectedOutput: `{{"{{"}} first {{"}}"}} and {{"{{"}} second {{"}}"}}`,
			expectedCount:  1,
		},
		{
			name:           "No templates",
			input:          "No templates here",
			expectedOutput: "No templates here",
			expectedCount:  0,
		},
		{
			name:           "With single quotes",
			input:          "{{ 'workflow.name' }}",
			expectedOutput: `{{"{{"}} 'workflow.name' {{"}}"}}`,
			expectedCount:  1,
		},
		{
			name:           "With double quotes",
			input:          `{{ "workflow.name" }}`,
			expectedOutput: `{{"{{"}} "workflow.name" {{"}}"}}`,
			expectedCount:  1,
		},
		{
			name:           "Complex nested case",
			input:          `{{ if eq .Value "test" }}{{ .Name }}{{ else }}{{ .Default }}{{ end }}`,
			expectedOutput: `{{"{{"}} if eq .Value "test" {{"}}"}}{{"{{"}} .Name {{"}}"}}{{"{{"}} else {{"}}"}}{{"{{"}} .Default {{"}}"}}{{"{{"}} end {{"}}"}}`,
			expectedCount:  5,
		},
	}

	escaper := &TemplateEscaper{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, count := escaper.Escape(tt.input)
			if got != tt.expectedOutput {
				t.Errorf("TemplateEscaper.Escape() got = %v, want %v", got, tt.expectedOutput)
			}
			if count != tt.expectedCount {
				t.Errorf("TemplateEscaper.Escape() count = %v, want %v", count, tt.expectedCount)
			}
		})
	}
}

func TestTemplateEscaper_PreventDoubleEscaping(t *testing.T) {
	escaper := &TemplateEscaper{}

	// First escape
	input := "artifacts/{{ workflow.name }}"
	escaped, count := escaper.Escape(input)
	expected := `artifacts/{{"{{"}} workflow.name {{"}}"}}`

	if escaped != expected {
		t.Errorf("First escape failed, got = %v, want %v", escaped, expected)
	}
	if count != 1 {
		t.Errorf("First escape count incorrect, got = %v, want %v", count, 1)
	}

	// Second escape (should not change anything)
	doubleEscaped, count := escaper.Escape(escaped)
	if doubleEscaped != expected {
		t.Errorf("Double escaping occurred, got = %v, want %v", doubleEscaped, expected)
	}
	if count != 0 {
		t.Errorf("Second escape should have count 0, got = %v", count)
	}
}

func TestTemplateEscaper_PreventDoubleEscapingWithQuotes(t *testing.T) {
	escaper := &TemplateEscaper{}

	// Test with double quotes
	t.Run("Double quotes", func(t *testing.T) {
		input := `{{ "quoted value" }}`
		escaped, count := escaper.Escape(input)
		expected := `{{"{{"}} "quoted value" {{"}}"}}`

		if escaped != expected {
			t.Errorf("First escape with double quotes failed, got = %v, want %v", escaped, expected)
		}
		if count != 1 {
			t.Errorf("First escape count incorrect, got = %v, want %v", count, 1)
		}

		// Second escape (should not change anything)
		doubleEscaped, count := escaper.Escape(escaped)
		if doubleEscaped != expected {
			t.Errorf("Double escaping with double quotes occurred, got = %v, want %v", doubleEscaped, expected)
		}
		if count != 0 {
			t.Errorf("Second escape should have count 0, got = %v", count)
		}
	})

	// Test with single quotes
	t.Run("Single quotes", func(t *testing.T) {
		input := `{{ 'quoted value' }}`
		escaped, count := escaper.Escape(input)
		expected := `{{"{{"}} 'quoted value' {{"}}"}}`

		if escaped != expected {
			t.Errorf("First escape with single quotes failed, got = %v, want %v", escaped, expected)
		}
		if count != 1 {
			t.Errorf("First escape count incorrect, got = %v, want %v", count, 1)
		}

		// Second escape (should not change anything)
		doubleEscaped, count := escaper.Escape(escaped)
		if doubleEscaped != expected {
			t.Errorf("Double escaping with single quotes occurred, got = %v, want %v", doubleEscaped, expected)
		}
		if count != 0 {
			t.Errorf("Second escape should have count 0, got = %v", count)
		}
	})
}
