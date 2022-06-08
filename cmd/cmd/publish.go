package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var pubCmd = &cobra.Command{
	Use:           "pub",
	Short:         "pub to a topic",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           pub,
}

var clientFlags struct {
	wipDB           bool
	addPingPipeline bool
}

func pub(cmd *cobra.Command, args []string) {
	client := initRedis()
	payload, err := json.Marshal(User{Name: "aidan"})
	if err != nil {
		panic(err)
	}
	err = client.Publish(channel, payload)
	fmt.Println(err)
	if err != nil {
		return
	}
}

func init() {
	rootCmd.AddCommand(pubCmd)
	flagSet := pubCmd.Flags()
	flagSet.BoolVarP(&clientFlags.addPingPipeline, "add-ping", "", false, "add one ping job to the pipeline")
}
