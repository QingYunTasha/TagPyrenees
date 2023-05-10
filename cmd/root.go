package cmd

import (
	"fmt"
	"os"

	"github.com/ory/viper"
	"github.com/spf13/cobra"
)

var (
	// used for flags
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "tagpyrenees",
		Short: "TagPyrenees is a command for tag system",
		Long: ` A command for tag system in Go.
				Complete documentation is available at https://github.com/QingYunTasha/TagPyrenees`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
