package cmd

import (
	"context"
	"os"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/log"
	"github.com/ikedam/pictmanager/pkg/uploader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "uploader",
	Short:        "uploader scans a directory and upload all files to pictmanager",
	Args:         cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
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
		uploader, err := uploader.New(&cfg)
		if err != nil {
			return err
		}
		ctx := context.Background()
		return uploader.Scan(ctx, args[0])
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
