/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/auth"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/files"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/infrastructure"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/projects"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/psql"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/server"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/users"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/utils"
)

var port int
var connStr string
var verbose bool
var jwtSecret string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackmap-postgis",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		if verbose {
			logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
		}

		logger.Info("Starting")
		slog.SetDefault(logger)

		if err := serve(logger, connStr, port); err != nil {
			logger.Error("failed starting server", "error", err)
			panic(err)
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stackmap-postgis.yaml)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8000, "Port to answer on")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&connStr, "connect", "c", "", "PSQL connection string: postgresql://[user[:password]@][netloc][:port]")

	defjwtsecret, err := utils.GetMacAddr()
	if err != nil {
		defjwtsecret = utils.RandStringBytes(10)
	}

	rootCmd.Flags().StringVarP(&jwtSecret, "jwtsecret", "j", defjwtsecret, "Set jwt secret to use when issuing jwt tokens")

	if rootCmd.Flags().Lookup("jwtsecret").Changed == false {
		slog.Info("Setting jwtsecret to default", "jwtsecret", defjwtsecret)
	}
	//viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag(utils.JWT_SECRET_KEY, rootCmd.Flags().Lookup("jwtsecret"))
	// Initialize Viper
	//cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initConfig)

}

func initConfig() {
	// Optionally load configurations (e.g., from file or environment variables)
	viper.AutomaticEnv() // Bind environment variables automatically
}

func serve(logger *slog.Logger, connStr string, port int) error {
	var db *sql.DB
	var err error
	if connStr != "" {
		if db, err = psql.SetupFromConnStr(connStr); err != nil {
			return fmt.Errorf("unable to connect to database using connection string")
		}
	} else {
		db, err = psql.SetupFromEnv()
		if err != nil {
			return err
		}
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}
	server := server.Server{
		//ProjectStore:   projects.NewSQLProjectStore(db),
		ProjectStore:      projects.NewSQLProjectStore(db),
		FileStore:         files.NewPSQLFilestStore(db),
		Infrastructure:    infrastructure.NewPSQLInfrastructure(db),
		UserStore:         users.NewPSQLUserStore(db),
		RefreshTokenStore: auth.NewPSQLTokenStore(db),
		Port:              port,
		Logger:            logger,
	}
	return server.Serve()

}
