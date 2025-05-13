package cmd

import (
	"os"

	"github.com/gkwa/jestingjaguar/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose int
	rootCmd = &cobra.Command{
		Use:   "jestingjaguar",
		Short: "Escape golang template brackets",
		Long: `A tool to escape golang template brackets to prevent interpolation.

For example:
  artifacts/{{ workflow.name }}
becomes:
  artifacts/{{"{{"}} workflow.name {{"}}"}}`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetVerbosity(verbose)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jestingjaguar.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase verbosity (can be used multiple times)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jestingjaguar")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logger.Debug("Using config file: %s", viper.ConfigFileUsed())
	}
}
