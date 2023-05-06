package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/ikedam/pictmanager/pkg/config"
	"github.com/ikedam/pictmanager/pkg/log"
	"github.com/ikedam/pictmanager/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "server",
	Short:        "server starts api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		var cfg config.Config
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return err
		}
		err = cfg.Build()
		if err != nil {
			return err
		}
		project := viper.GetString("project")
		if project != "" {
			os.Setenv("CLOUDSDK_CORE_PROJECT", project)
		}
		s, err := server.New(ctx, &cfg)
		if err != nil {
			return err
		}
		return s.Start(ctx)
	},
}

func init() {
	cobra.OnInitialize(initLevel)
	cobra.OnInitialize(setDefault)

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().String("log-level", "debug", "Log level.")
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))
	rootCmd.Flags().String("project", "", "ID of Google Cloud Project.")
	viper.BindPFlag("project", rootCmd.Flags().Lookup("project"))
	rootCmd.Flags().Int("port", port, "Port to bind")
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	rootCmd.Flags().String("gcs", "", "GCS directory to store images.")
	viper.BindPFlag("gcs", rootCmd.Flags().Lookup("gcs"))
	rootCmd.Flags().String("gcs-public-base", "https://storage.googleapis.com", "GCS public base URL.")
	viper.BindPFlag("gcsPublicBase", rootCmd.Flags().Lookup("gcs-public-base"))
}

func initLevel() {
	level := viper.GetString("logLevel")
	err := log.SetLevelByName(level)
	if err != nil {
		log.Fatalf(context.Background(), "Invalid log level specified: %v: %+v", level, err)
	}
}

func setDefault() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("pictmanager")

	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			rootCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
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
