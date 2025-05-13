package cmd

import (
	"github.com/gkwa/jestingjaguar/internal/escaper"
	"github.com/gkwa/jestingjaguar/internal/logger"
	"github.com/spf13/cobra"
)

var escapeCmd = &cobra.Command{
	Use:   "escape [file or directory]",
	Short: "Escape golang template delimiters",
	Long: `Escape golang template delimiters to prevent interpolation.
If a directory is provided, it will recursively escape all files within it.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		service := escaper.NewService()
		stats, err := service.Process(path)
		if err != nil {
			return err
		}

		logger.Info("Processed %d files, performed %d escapes", stats.FilesProcessed, stats.EscapesPerformed)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(escapeCmd)
}
