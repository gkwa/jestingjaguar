package escaper

import (
	"os"

	"github.com/gkwa/jestingjaguar/internal/logger"
)

// FileProcessor defines an interface for processing files
type FileProcessor interface {
	ProcessFile(filePath string) (int, error)
}

// DefaultFileProcessor is the default implementation of FileProcessor
type DefaultFileProcessor struct {
	escaper Escaper
}

// ProcessFile reads a file, escapes the content, and writes it back
func (p *DefaultFileProcessor) ProcessFile(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	logger.Trace("Original content: %s", string(content))

	escaped, count := p.escaper.Escape(string(content))
	if count == 0 {
		logger.Debug("No escapes needed for file: %s", filePath)
		return 0, nil
	}

	logger.Trace("Escaped content: %s", escaped)
	logger.Debug("Performed %d escapes in file: %s", count, filePath)

	err = os.WriteFile(filePath, []byte(escaped), 0o644)
	if err != nil {
		return 0, err
	}

	return count, nil
}
