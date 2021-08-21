package cmd

import (
	"buty4649/omusoba/cmdutil"
	"buty4649/omusoba/sensor/usb"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch latest sensor data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFetchCmd()
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func runFetchCmd() error {
	s := usb.New(cfgUsb)
	data, err := s.FetchLatestData()
	if err != nil {
		return err
	}

	cmdutil.MarshalPrint(cfgFormat, *data)
	return nil
}
