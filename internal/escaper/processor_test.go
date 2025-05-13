package escaper

import (
	"os"
	"path/filepath"
	"testing"
)

// MockEscaper implements the Escaper interface for testing
type MockEscaper struct {
	returnContent string
	returnCount   int
}

func (m *MockEscaper) Escape(content string) (string, int) {
	return m.returnContent, m.returnCount
}

func TestDefaultFileProcessor_ProcessFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "processor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "{{ test }}"
	err = os.WriteFile(testFilePath, []byte(originalContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test successful processing
	t.Run("Successful processing", func(t *testing.T) {
		mock := &MockEscaper{
			returnContent: `{{"{{"}} test {{"}}"}}`,
			returnCount:   1,
		}
		processor := &DefaultFileProcessor{escaper: mock}

		count, err := processor.ProcessFile(testFilePath)
		if err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		if count != 1 {
			t.Errorf("Expected count 1, got %d", count)
		}

		content, err := os.ReadFile(testFilePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(content) != mock.returnContent {
			t.Errorf("File content not updated correctly, got %s, want %s", string(content), mock.returnContent)
		}
	})

	// Test no changes needed
	t.Run("No changes needed", func(t *testing.T) {
		// Reset the file
		err = os.WriteFile(testFilePath, []byte(originalContent), 0o644)
		if err != nil {
			t.Fatalf("Failed to reset test file: %v", err)
		}

		mock := &MockEscaper{
			returnContent: originalContent,
			returnCount:   0,
		}
		processor := &DefaultFileProcessor{escaper: mock}

		count, err := processor.ProcessFile(testFilePath)
		if err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		if count != 0 {
			t.Errorf("Expected count 0, got %d", count)
		}
	})

	// Test error handling for non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		processor := &DefaultFileProcessor{escaper: &MockEscaper{}}
		_, err := processor.ProcessFile(filepath.Join(tempDir, "nonexistent.txt"))
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
}
