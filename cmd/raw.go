package cmd

import (
	"buty4649/omusoba/cmdutil"
	"buty4649/omusoba/sensor/usb"
	"strconv"

	"github.com/spf13/cobra"
)

var rawCmd = &cobra.Command{
	Use:   "raw <command> <address> [data]...",
	Short: "Send the raw data specified by the argument as is.",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRawCmd(args)
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)
}

func runRawCmd(args []string) error {
	cmd, err := strconv.ParseInt(args[0], 0, 8)
	if err != nil {
		return nil
	}

	addr, err := strconv.ParseInt(args[1], 0, 16)
	if err != nil {
		return nil
	}

	data := []byte{}
	for i := 2; i < len(args); i++ {
		d, err := strconv.ParseInt(args[2], 0, 8)
		if err != nil {
			return nil
		}
		data = append(data, byte(d))
	}

	s := usb.New(cfgUsb)
	response, err := s.RawSend(int(cmd), int(addr), data)
	if err != nil {
		return err
	}

	cmdutil.MarshalPrint(cfgFormat, struct {
		Response []byte `desc:"response"`
	}{response})
	return nil
}
