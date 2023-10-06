package cmd

import (
	"errors"
	"fmt"
	"github.com/geowa4/odx/pkg/scripts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"syscall"
)

const ConfigFileName = "odx.yaml"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "odx",
	Short: "Run one-off scripts and dispose of them",
	Args:  cobra.RangeArgs(1, 2),

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		sourceArgs := args
		if len(args) == 1 {
			aliases := viper.GetStringMap("aliases")
			if alias, ok := aliases[args[0]]; !ok {
				cobra.CheckErr(fmt.Errorf("solo argument %s provided but not found in aliases", args[0]))
			} else {
				aliasSlice := alias.([]interface{})
				sourceArgs = make([]string, len(aliasSlice))
				for i, elem := range aliasSlice {
					sourceArgs[i] = elem.(string)
				}
			}
		}

		sources := viper.GetStringMap("sources")
		if source, ok := sources[sourceArgs[0]]; !ok {
			cobra.CheckErr(fmt.Errorf("source %s does not exist", sourceArgs[0]))
		} else {
			sourceMap := source.(map[string]interface{})
			repo := sourceMap["github"].(string)
			branch := sourceMap["branch"].(string)
			path := sourceMap["path"].(string)
			localScriptFullPath := scripts.DownloadScriptFromGitHub(repo, branch, path, sourceArgs[1])
			//defer func(name string) {
			//	_ = os.Remove(name)
			//}(localScriptFullPath)
			cobra.CheckErr(syscall.Exec(
				localScriptFullPath,
				[]string{},
				os.Environ(),
			))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/odx.yaml)")
}

func getConfigDir() string {
	if xdgConfigHome, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		return xdgConfigHome
	} else {
		// TODO: do we care about this error?
		userConfigDir, _ := os.UserConfigDir()
		return userConfigDir
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		configDir := getConfigDir()
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(ConfigFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			cobra.CheckErr(fmt.Errorf("bad config file (%s): %q", viper.ConfigFileUsed(), err))
		}
	}
}
