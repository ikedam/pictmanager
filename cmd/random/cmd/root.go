package cmd

import (
	"context"
	"os"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/log"
	"github.com/ikedam/pictmanager/pkg/random"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "random",
	Short:        "random generates random values for images. All images will be updated.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return err
		}
		project := viper.GetString("project")
		if project != "" {
			os.Setenv("CLOUDSDK_CORE_PROJECT", project)
		}
		ctx := context.Background()
		random, err := random.New(ctx, &cfg)
		if err != nil {
			return err
		}
		return random.Scan(ctx)
	},
}

func init() {
	cobra.OnInitialize(initLevel)

	rootCmd.PersistentFlags().String("log-level", "info", "Log level.")
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))
	rootCmd.Flags().String("project", "", "ID of Google Cloud Project.")
	viper.BindPFlag("project", rootCmd.Flags().Lookup("project"))
	rootCmd.Flags().String("gcs", "", "GCS directory to store images.")
	viper.BindPFlag("gcs", rootCmd.Flags().Lookup("gcs"))
}

func initLevel() {
	level := viper.GetString("logLevel")
	err := log.SetLevelByName(level)
	if err != nil {
		log.Fatalf(context.Background(), "Invalid log level specified: %v: %+v", level, err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(context.Background(), "command failed", zap.Error(err))
	}
}

// SetVersion sets version of the command
func SetVersion(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
