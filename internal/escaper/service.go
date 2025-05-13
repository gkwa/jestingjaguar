package escaper

import (
	"os"
	"path/filepath"

	"github.com/gkwa/jestingjaguar/internal/logger"
)

// Service provides methods to escape template delimiters
type Service struct {
	processor FileProcessor
}

// Stats holds statistics about the escaping operation
type Stats struct {
	FilesProcessed   int
	EscapesPerformed int
}

// NewService creates a new escaper service
func NewService() *Service {
	return &Service{
		processor: &DefaultFileProcessor{
			escaper: &TemplateEscaper{},
		},
	}
}

// Process handles escaping for a file or directory
func (s *Service) Process(path string) (Stats, error) {
	stats := Stats{}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return stats, err
	}

	if fileInfo.IsDir() {
		return s.processDirectory(path)
	}

	escapesPerformed, err := s.processor.ProcessFile(path)
	if err != nil {
		return stats, err
	}

	stats.FilesProcessed = 1
	stats.EscapesPerformed = escapesPerformed
	return stats, nil
}

// processDirectory recursively processes all files in a directory
func (s *Service) processDirectory(dirPath string) (Stats, error) {
	stats := Stats{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			logger.Debug("Processing file: %s", path)
			escapesPerformed, err := s.processor.ProcessFile(path)
			if err != nil {
				logger.Error("Error processing file %s: %v", path, err)
				return err
			}

			stats.FilesProcessed++
			stats.EscapesPerformed += escapesPerformed
		}
		return nil
	})

	return stats, err
}
