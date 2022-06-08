package cmd

import (
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
	for i := 1; i < 5; i++ {
		name := fmt.Sprintf("hey %d", i)
		payload, err := client.Encode(User{Name: name})
		if err != nil {
			panic(err)
		}
		err = client.Publish(channel, payload)
		if err != nil {
			return
		}
	}

}

func init() {
	rootCmd.AddCommand(pubCmd)
	flagSet := pubCmd.Flags()
	flagSet.BoolVarP(&clientFlags.addPingPipeline, "add-ping", "", false, "add one ping job to the pipeline")
}
