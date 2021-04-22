package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	programName = "partition-watchdog"
)

var (
	rootCmd = &cobra.Command{
		Use:          programName,
		SilenceUsage: true,
	}
)

// Execute is the entrypoint of the cient-go application
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		if viper.GetBool("debug") {
			st := errors.WithStack(err)
			fmt.Printf("%+v", st)
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Duration("checkinterval", 10*time.Second, "time between connection checks")
	rootCmd.PersistentFlags().Int("tries", 6, "number of checks partition is considered (un)available")
	rootCmd.PersistentFlags().String("target", "", "target to check (e.g. 212.1.5.7:80")
	rootCmd.PersistentFlags().Duration("timeout", 2*time.Second, "connection timeout for checks")
	rootCmd.PersistentFlags().String("deployment", "kube-controller-manager", "name of deployment to scale")

	rootCmd.AddCommand(checkPartition)

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatalf("error setup root cmd:%v", err)
	}
}

func initConfig() {
	viper.SetEnvPrefix(strings.ToUpper(programName))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
