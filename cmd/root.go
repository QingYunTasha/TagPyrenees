package cmd

import (
	"fmt"
	"os"

	"tagpyrenees/usecase"

	"github.com/ory/viper"
	"github.com/spf13/cobra"
)

var (
	// used for flags
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "tagpyrenees",
		Short: "The command-line tool for tagPyrenees system",
		Long: ` The command-line tool for tagPyrenees system
		Complete documentation is available at https://github.com/QingYunTasha/TagPyrenees`,
	}

	queryCmd = &cobra.Command{
		Use:   "query [path] [tag]",
		Short: "query by the tag recursive the given path and subpath",
		Long: ` query by the tag
				Complete documentation is available at https://github.com/QingYunTasha/TagPyrenees`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			tag := args[1]
			err := usecase.QueryByTag(path, tag)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	listTagsCmd = &cobra.Command{
		Use:   "listtags [path]",
		Short: "list all tags recursive the given path and subpath",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			err := usecase.ListTags(path)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	//queryCmd.Flags().StringP("tag", "t", "", "the tag to query")
	queryCmd.PersistentFlags()

	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(listTagsCmd)
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
