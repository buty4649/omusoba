package cmd

import (
	"buty4649/omusoba/cmdutil"
	"buty4649/omusoba/sensor/usb"

	"github.com/spf13/cobra"
)

var modeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Show device mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModeCmd()
	},
}

func init() {
	rootCmd.AddCommand(modeCmd)
}

func runModeCmd() error {
	s := usb.New(cfgUsb)
	mode, err := s.Mode()
	if err != nil {
		return err
	}

	cmdutil.MarshalPrint(cfgFormat, *mode)
	return nil
}
