package cmd

import (
	"buty4649/omusoba/cmdutil"
	"buty4649/omusoba/sensor/usb"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show device information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInfoCmd()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfoCmd() error {
	s := usb.New(cfgUsb)
	info, err := s.Info()
	if err != nil {
		return err
	}

	cmdutil.MarshalPrint(cfgFormat, *info)
	return nil
}
