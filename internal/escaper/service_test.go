package escaper

import (
	"os"
	"path/filepath"
	"testing"
)

// MockFileProcessor implements the FileProcessor interface for testing
type MockFileProcessor struct {
	processFileFunc func(string) (int, error)
}

func (m *MockFileProcessor) ProcessFile(filePath string) (int, error) {
	return m.processFileFunc(filePath)
}

func TestService_Process_File(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_file.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	mockProcessor := &MockFileProcessor{
		processFileFunc: func(filePath string) (int, error) {
			return 5, nil
		},
	}

	service := &Service{processor: mockProcessor}
	stats, err := service.Process(tempFile.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if stats.FilesProcessed != 1 {
		t.Errorf("Expected 1 file processed, got %d", stats.FilesProcessed)
	}

	if stats.EscapesPerformed != 5 {
		t.Errorf("Expected 5 escapes performed, got %d", stats.EscapesPerformed)
	}
}

func TestService_Process_Directory(t *testing.T) {
	// Create a temporary directory with files
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Create files
	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "file2.txt"),
		filepath.Join(subDir, "file3.txt"),
	}

	for _, file := range files {
		if err := os.WriteFile(file, []byte("test content"), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Mock the processor to return different escape counts
	mockCounts := map[string]int{
		files[0]: 3,
		files[1]: 0,
		files[2]: 2,
	}

	mockProcessor := &MockFileProcessor{
		processFileFunc: func(filePath string) (int, error) {
			return mockCounts[filePath], nil
		},
	}

	service := &Service{processor: mockProcessor}
	stats, err := service.Process(tempDir)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if stats.FilesProcessed != 3 {
		t.Errorf("Expected 3 files processed, got %d", stats.FilesProcessed)
	}

	// Expected escapes: 3 + 0 + 2 = 5
	if stats.EscapesPerformed != 5 {
		t.Errorf("Expected 5 escapes performed, got %d", stats.EscapesPerformed)
	}
}
