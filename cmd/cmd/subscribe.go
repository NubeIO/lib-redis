package cmd

import (
	"fmt"
	"github.com/NubeIO/lib-redis/libredis"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "subscribe",
	Short:         "subscribe to a topic",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           runRoot,
}

var rootFlags struct {
	server bool
	config string
	wipeDb bool
}

var channel = "test"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func initRedis() libredis.Client {
	client, err := libredis.New(&libredis.Config{})
	if err != nil {
		return nil
	}
	return client
}

func runRoot(cmd *cobra.Command, args []string) {
	client := initRedis()
	messages := make(chan string)
	go func() {
		for {
			msg := <-messages
			user := &User{}
			err := client.Decode(msg, user)
			fmt.Println(user.Name, err)
		}
	}()
	fmt.Println("subscribing to channel:", channel)
	client.Subscribe(channel, messages)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//color.Magenta(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pFlagSet := rootCmd.PersistentFlags()
	pFlagSet.StringVarP(&rootFlags.config, "config", "", "config.yaml", "set config path example ./config.yaml")
	pFlagSet.BoolVarP(&rootFlags.server, "server", "", false, "run server")
	pFlagSet.BoolVarP(&rootFlags.wipeDb, "wipe", "", false, "delete the db after server has started")
}
