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
		Use:   "query",
		Short: "query by the tag",
		Long: ` query by the tag
				Complete documentation is available at https://github.com/QingYunTasha/TagPyrenees`,
		Run: func(cmd *cobra.Command, args []string) {
			usecase.QueryByTag(cmd.Flag("tag").Value.String())
		},
	}

	listTagsCmd = &cobra.Command{
		Use:   "listtags",
		Short: "list all tags",
		Run: func(cmd *cobra.Command, args []string) {
			err := usecase.ListTags()
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

	queryCmd.Flags().StringP("tag", "t", "", "the tag to query")

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
